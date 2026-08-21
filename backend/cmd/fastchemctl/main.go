// Command fastchemctl manages a FastChem database from the shell.
//
// It talks to SQLite directly rather than to the admin API, which is the whole
// point: it works before the first admin exists, while the server is down, and
// over `docker exec`. The web console is the better tool for everything else —
// this covers what a browser cannot reach.
//
// Two things are safe to change while the server is running, by design:
//
//   - is_admin is read on every admin request, so grant/revoke applies at once
//   - anticheat_rules is reloaded every 30 seconds, and `rules set` reloads the
//     local cache too, so a retune lands without a restart
//
// Writes that touch a player mid-match (user delete) are not — see the warning
// on that command.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/takumi/fastchem/internal/anticheat"
	"github.com/takumi/fastchem/internal/database"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

const usage = `fastchemctl — manage a FastChem database

Usage:
  fastchemctl [-db path] <command> [flags]

Commands:
  status                    Counts, table sizes and anti-cheat posture
  user list                 List accounts        [-search s] [-limit n] [-sort points|rating|games|findings|created|username]
  user show <username>      One account in full
  user grant <username>     Give admin console access
  user revoke <username>    Take it away
  user passwd <username>    Set a new password (prompted, not echoed)
  user delete <username>    Delete the account and everything it owns
  rules list                Anti-cheat rules, their thresholds and hit counts
  rules set <name>          Retune a rule        [-enabled=false] [-action observe|reject] [-param k=v]
  findings                  Recent detections    [-rule r] [-mode m] [-user u] [-limit n]
  backup <dest>             Consistent copy of the database, safe while running

The database defaults to $DB_PATH, then fastchem.db — the same order the
server uses, so both see the same file when run from the same directory.
`

func main() {
	// The server logs at info; a CLI that printed the same startup chatter
	// would bury its own output.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})))

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	command, args := os.Args[1], os.Args[2:]
	var err error

	switch command {
	case "status":
		err = cmdStatus(args)
	case "user":
		err = cmdUser(args)
	case "rules":
		err = cmdRules(args)
	case "findings":
		err = cmdFindings(args)
	case "backup":
		err = cmdBackup(args)
	case "help", "-h", "--help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", command, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// open resolves the database path the same way the server does and runs the
// migrations, so the CLI works against a database file that predates a column
// it needs. database.Init exits on failure rather than returning.
func open(fs *flag.FlagSet, args []string) *flag.FlagSet {
	dbPath := fs.String("db", "", "path to the SQLite database (default $DB_PATH, then fastchem.db)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: fastchemctl %s [flags]\n\nFlags:\n", fs.Name())
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	path := *dbPath
	if path == "" {
		path = os.Getenv("DB_PATH")
	}
	if path == "" {
		path = "fastchem.db"
	}
	database.Init(path)
	return fs
}

// openSubject is open() for the commands that name one thing — a username, a
// rule, a destination file.
//
// The flag package stops parsing at the first non-flag token, so
// `rules set impossible_speed -action reject` would otherwise hand -action
// straight through as a positional and silently ignore it. Peeling a leading
// subject off before parsing makes both orders work, and the documented order
// is the one that reads naturally.
func openSubject(fs *flag.FlagSet, args []string, usage string) (string, error) {
	var subject string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		subject, args = args[0], args[1:]
	}
	open(fs, args)

	if subject == "" {
		if fs.NArg() != 1 {
			return "", errors.New("usage: fastchemctl " + usage)
		}
		return fs.Arg(0), nil
	}
	if fs.NArg() != 0 {
		return "", errors.New("usage: fastchemctl " + usage)
	}
	return subject, nil
}

// ─── status ────────────────────────────────────────────────────

func cmdStatus(args []string) error {
	fs := open(flag.NewFlagSet("status", flag.ExitOnError), args)
	_ = fs
	defer database.Close()
	ctx := context.Background()

	var users, admins, scores, matches, ranked, findings, findings24h int
	var points sql.NullInt64
	row := func(query string, dest ...any) {
		if err := database.DB.QueryRowContext(ctx, query).Scan(dest...); err != nil {
			fmt.Fprintln(os.Stderr, "warning:", err)
		}
	}
	row("SELECT COUNT(*), COALESCE(SUM(is_admin), 0), COALESCE(SUM(total_points), 0) FROM users", &users, &admins, &points)
	row("SELECT COUNT(*) FROM scores", &scores)
	row("SELECT COUNT(*) FROM matches", &matches)
	row("SELECT COUNT(*) FROM ranked_matches", &ranked)
	row("SELECT COUNT(*), COALESCE(SUM(at >= datetime('now', '-1 day')), 0) FROM anticheat_findings", &findings, &findings24h)

	w := table()
	fmt.Fprintf(w, "accounts\t%d (%d admin)\n", users, admins)
	fmt.Fprintf(w, "total points\t%d\n", points.Int64)
	fmt.Fprintf(w, "casual games\t%d\n", scores)
	fmt.Fprintf(w, "solo matches\t%d\n", matches)
	fmt.Fprintf(w, "ranked matches\t%d\n", ranked)
	fmt.Fprintf(w, "findings\t%d (%d in 24h)\n", findings, findings24h)
	w.Flush()

	if admins == 0 {
		fmt.Println("\nNo admin accounts: the console is unreachable.")
		fmt.Println("Fix with:  fastchemctl user grant <username>")
	}

	store, err := rulesStore(ctx)
	if err != nil {
		return err
	}
	var enforcing []string
	for _, r := range store.Snapshot() {
		if r.Enabled && r.Action == anticheat.ActionReject {
			enforcing = append(enforcing, r.Name)
		}
	}
	if len(enforcing) == 0 {
		fmt.Println("\nAnti-cheat: observing only, no rule rejects answers.")
	} else {
		fmt.Printf("\nAnti-cheat: REJECTING answers on %s.\n", strings.Join(enforcing, ", "))
	}
	return nil
}

// ─── user ──────────────────────────────────────────────────────

func cmdUser(args []string) error {
	if len(args) == 0 {
		return errors.New("user needs a subcommand: list, show, grant, revoke, passwd, delete")
	}
	sub, rest := args[0], args[1:]

	switch sub {
	case "list":
		return userList(rest)
	case "show":
		return userShow(rest)
	case "grant":
		return userSetAdmin(rest, true)
	case "revoke":
		return userSetAdmin(rest, false)
	case "passwd":
		return userPasswd(rest)
	case "delete":
		return userDelete(rest)
	default:
		return fmt.Errorf("unknown user subcommand %q", sub)
	}
}

// userSortColumns whitelists -sort. The value is spliced into SQL because a
// placeholder cannot carry a column name, so it may only come from this map.
var userSortColumns = map[string]string{
	"points":   "u.total_points DESC",
	"rating":   "u.rating DESC",
	"games":    "games DESC",
	"findings": "findings DESC",
	"created":  "u.created_at DESC",
	"username": "u.username ASC",
}

func userList(args []string) error {
	fs := flag.NewFlagSet("user list", flag.ExitOnError)
	search := fs.String("search", "", "match usernames containing this text")
	limit := fs.Int("limit", 50, "maximum rows")
	sortKey := fs.String("sort", "points", "points|rating|games|findings|created|username")
	open(fs, args)
	defer database.Close()

	order, ok := userSortColumns[strings.ToLower(*sortKey)]
	if !ok {
		return fmt.Errorf("unknown sort %q", *sortKey)
	}

	needle := strings.ToLower(*search)
	rows, err := database.DB.QueryContext(context.Background(), `
		SELECT u.id, u.username, u.is_admin, u.total_points, u.rating,
		       u.ranked_wins, u.ranked_losses,
		       COUNT(s.id) AS games,
		       (SELECT COUNT(*) FROM anticheat_findings f WHERE f.user_id = u.id) AS findings,
		       u.created_at
		FROM users u LEFT JOIN scores s ON s.user_id = u.id
		WHERE ? = '' OR LOWER(u.username) LIKE ?
		GROUP BY u.id
		ORDER BY `+order+`
		LIMIT ?`, needle, "%"+needle+"%", *limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	w := table()
	fmt.Fprintln(w, "ID\tUSERNAME\tADMIN\tPOINTS\tRATING\tW/L\tGAMES\tFLAGS\tCREATED")
	count := 0
	for rows.Next() {
		var (
			id, points, rating, wins, losses, games, findings int64
			username                                          string
			isAdmin                                           bool
			created                                           time.Time
		)
		if err := rows.Scan(&id, &username, &isAdmin, &points, &rating, &wins, &losses, &games, &findings, &created); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%d/%d\t%d\t%d\t%s\n",
			id, username, yesNo(isAdmin), points, rating, wins, losses, games, findings,
			created.Format(time.DateOnly))
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	w.Flush()

	if count == 0 {
		fmt.Println("no accounts matched")
	}
	return nil
}

func userShow(args []string) error {
	fs := flag.NewFlagSet("user show", flag.ExitOnError)
	name, err := openSubject(fs, args, "user show <username>")
	defer database.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()

	var (
		id, points, rating, wins, losses, highest int64
		username                                  string
		isAdmin                                   bool
		created                                   time.Time
	)
	err = database.DB.QueryRowContext(ctx, `
		SELECT id, username, is_admin, total_points, rating, ranked_wins, ranked_losses, highest_rating, created_at
		FROM users WHERE LOWER(username) = LOWER(?)`, name,
	).Scan(&id, &username, &isAdmin, &points, &rating, &wins, &losses, &highest, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no such account %q", name)
	}
	if err != nil {
		return err
	}

	var games, answered, correct, findings int64
	database.DB.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(total_answered), 0), COALESCE(SUM(correct_answers), 0)
		FROM scores WHERE user_id = ?`, id).Scan(&games, &answered, &correct)
	database.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM anticheat_findings WHERE user_id = ?", id).Scan(&findings)

	w := table()
	fmt.Fprintf(w, "id\t%d\n", id)
	fmt.Fprintf(w, "username\t%s\n", username)
	fmt.Fprintf(w, "admin\t%s\n", yesNo(isAdmin))
	fmt.Fprintf(w, "created\t%s\n", created.Format(time.DateTime))
	fmt.Fprintf(w, "points\t%d\n", points)
	fmt.Fprintf(w, "rating\t%d (peak %d)\n", rating, highest)
	fmt.Fprintf(w, "ranked\t%d W / %d L\n", wins, losses)
	fmt.Fprintf(w, "games\t%d\n", games)
	if answered > 0 {
		fmt.Fprintf(w, "accuracy\t%.1f%% (%d/%d)\n", float64(correct)/float64(answered)*100, correct, answered)
	}
	fmt.Fprintf(w, "findings\t%d\n", findings)
	return w.Flush()
}

func userSetAdmin(args []string, grant bool) error {
	name := "user revoke"
	if grant {
		name = "user grant"
	}
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	target, err := openSubject(fs, args, name+" <username>")
	defer database.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()

	// Refuse to remove the last admin: nothing in the running system can put it
	// back, and the console stays unreachable until someone runs this tool
	// again — which works, but the operator should mean it.
	if !grant {
		var admins int
		if err := database.DB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM users WHERE is_admin = 1").Scan(&admins); err != nil {
			return err
		}
		var targetIsAdmin bool
		database.DB.QueryRowContext(ctx,
			"SELECT is_admin FROM users WHERE LOWER(username) = LOWER(?)", target).Scan(&targetIsAdmin)
		if admins <= 1 && targetIsAdmin {
			return errors.New("that is the only admin account; grant another one first")
		}
	}

	res, err := database.DB.ExecContext(ctx,
		"UPDATE users SET is_admin = ? WHERE LOWER(username) = LOWER(?)", boolToInt(grant), target)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("no such account %q", target)
	}

	if grant {
		fmt.Printf("%s is now an admin. The change is live — is_admin is read on every request.\n", target)
	} else {
		fmt.Printf("%s is no longer an admin, effective immediately.\n", target)
	}
	return nil
}

func userPasswd(args []string) error {
	fs := flag.NewFlagSet("user passwd", flag.ExitOnError)
	name, err := openSubject(fs, args, "user passwd <username>")
	defer database.Close()
	if err != nil {
		return err
	}

	var id int64
	err = database.DB.QueryRowContext(context.Background(),
		"SELECT id FROM users WHERE LOWER(username) = LOWER(?)", name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no such account %q", name)
	}
	if err != nil {
		return err
	}

	password, err := readPassword("New password: ")
	if err != nil {
		return err
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	again, err := readPassword("Repeat: ")
	if err != nil {
		return err
	}
	if password != again {
		return errors.New("passwords do not match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if _, err := database.DB.ExecContext(context.Background(),
		"UPDATE users SET password_hash = ? WHERE id = ?", string(hash), id); err != nil {
		return err
	}

	fmt.Printf("password updated for %s. Existing tokens stay valid for up to 7 days.\n", name)
	return nil
}

func userDelete(args []string) error {
	fs := flag.NewFlagSet("user delete", flag.ExitOnError)
	yes := fs.Bool("yes", false, "skip the confirmation prompt")
	name, err := openSubject(fs, args, "user delete <username> [-yes]")
	defer database.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()

	var id int64
	var username string
	err = database.DB.QueryRowContext(ctx,
		"SELECT id, username FROM users WHERE LOWER(username) = LOWER(?)", name).Scan(&id, &username)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no such account %q", name)
	}
	if err != nil {
		return err
	}

	if !*yes {
		fmt.Printf("Delete %s (id %d) and every score, match and ranked result it owns?\n", username, id)
		fmt.Println("Its ranked matches disappear from the opponents' history too. This cannot be undone.")
		typed, _ := prompt("Type the username to confirm: ")
		if typed != username {
			return errors.New("aborted")
		}
	}

	if err := database.DeleteUser(ctx, id); err != nil {
		return err
	}

	fmt.Printf("deleted %s.\n", username)
	fmt.Println("If they were mid-match, the server still holds that match in memory until it ends or is swept.")
	return nil
}

// ─── rules ─────────────────────────────────────────────────────

// rulesStore seeds and loads the anti-cheat settings through the same Store
// the server uses, so the CLI reports and writes with identical semantics —
// including params merging over the compiled defaults.
func rulesStore(ctx context.Context) (*anticheat.Store, error) {
	store := anticheat.NewStore(database.DB)
	if err := store.Seed(ctx); err != nil {
		return nil, err
	}
	if err := store.Reload(ctx); err != nil {
		return nil, err
	}
	return store, nil
}

func cmdRules(args []string) error {
	if len(args) == 0 {
		return errors.New("rules needs a subcommand: list, set")
	}
	switch args[0] {
	case "list":
		return rulesList(args[1:])
	case "set":
		return rulesSet(args[1:])
	default:
		return fmt.Errorf("unknown rules subcommand %q", args[0])
	}
}

func rulesList(args []string) error {
	fs := flag.NewFlagSet("rules list", flag.ExitOnError)
	open(fs, args)
	defer database.Close()
	ctx := context.Background()

	store, err := rulesStore(ctx)
	if err != nil {
		return err
	}

	hits := map[string]int64{}
	rows, err := database.DB.QueryContext(ctx,
		"SELECT rule, COUNT(*) FROM anticheat_findings GROUP BY rule")
	if err != nil {
		return err
	}
	for rows.Next() {
		var rule string
		var n int64
		if err := rows.Scan(&rule, &n); err == nil {
			hits[rule] = n
		}
	}
	rows.Close()

	w := table()
	fmt.Fprintln(w, "RULE\tENABLED\tACTION\tFINDINGS\tPARAMS")
	for _, r := range store.Snapshot() {
		params := make([]string, 0, len(r.Params))
		for _, k := range sortedKeys(r.Params) {
			params = append(params, fmt.Sprintf("%s=%g", k, r.Params[k]))
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			r.Name, yesNo(r.Enabled), r.Action, hits[r.Name], strings.Join(params, " "))
	}
	return w.Flush()
}

// paramFlag collects repeated -param key=value pairs.
type paramFlag map[string]float64

func (p paramFlag) String() string { return "" }

func (p paramFlag) Set(value string) error {
	key, raw, ok := strings.Cut(value, "=")
	if !ok {
		return fmt.Errorf("expected key=value, got %q", value)
	}
	number, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fmt.Errorf("param %s: %w", key, err)
	}
	p[strings.TrimSpace(key)] = number
	return nil
}

func rulesSet(args []string) error {
	fs := flag.NewFlagSet("rules set", flag.ExitOnError)
	enabled := fs.Bool("enabled", true, "whether the rule runs")
	action := fs.String("action", "", "observe or reject (default: unchanged)")
	params := paramFlag{}
	fs.Var(params, "param", "threshold as key=value; repeatable")
	name, err := openSubject(fs, args, "rules set <name> [-enabled=false] [-action reject] [-param key=value]")
	defer database.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()

	store, err := rulesStore(ctx)
	if err != nil {
		return err
	}

	current := store.Rule(name)
	known := false
	for _, r := range store.Snapshot() {
		if r.Name == name {
			known = true
		}
	}
	if !known {
		return fmt.Errorf("unknown rule %q — see `fastchemctl rules list`", name)
	}

	// Unset flags leave the stored value alone, so retuning one threshold does
	// not silently reset the action or the other thresholds.
	next := current.Params.Clone()
	if next == nil {
		next = anticheat.Params{}
	}
	for k, v := range params {
		next[k] = v
	}

	wantAction := current.Action
	if *action != "" {
		parsed, ok := anticheat.ParseAction(*action)
		if !ok {
			return fmt.Errorf("unknown action %q (want observe or reject)", *action)
		}
		wantAction = parsed
	}
	wantEnabled := current.Enabled
	if isFlagSet(fs, "enabled") {
		wantEnabled = *enabled
	}

	if err := store.Save(ctx, name, wantEnabled, wantAction, next); err != nil {
		return err
	}

	fmt.Printf("%s: enabled=%s action=%s params=%v\n", name, yesNo(wantEnabled), wantAction, next)
	if wantAction == anticheat.ActionReject && wantEnabled {
		fmt.Println("This rule now marks flagged answers wrong. Players are not told which rule caught them.")
	}
	fmt.Println("A running server picks this up within 30 seconds.")
	return nil
}

// ─── findings ──────────────────────────────────────────────────

func cmdFindings(args []string) error {
	fs := flag.NewFlagSet("findings", flag.ExitOnError)
	rule := fs.String("rule", "", "only this rule")
	mode := fs.String("mode", "", "casual, match, ranked or room")
	user := fs.String("user", "", "only this username")
	limit := fs.Int("limit", 25, "maximum rows")
	open(fs, args)
	defer database.Close()

	rows, err := database.DB.QueryContext(context.Background(), `
		SELECT f.at, COALESCE(u.username, ''), f.subject, f.mode, f.rule, f.action, f.time_spent, f.detail
		FROM anticheat_findings f LEFT JOIN users u ON u.id = f.user_id
		WHERE (? = '' OR f.rule = ?)
		  AND (? = '' OR f.mode = ?)
		  AND (? = '' OR LOWER(u.username) = LOWER(?))
		ORDER BY f.at DESC, f.id DESC
		LIMIT ?`, *rule, *rule, *mode, *mode, *user, *user, *limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	w := table()
	fmt.Fprintln(w, "WHEN\tWHO\tMODE\tRULE\tACTION\tSECONDS\tDETAIL")
	count := 0
	for rows.Next() {
		var at time.Time
		var username, subject, mode, rule, action, detail string
		var seconds float64
		if err := rows.Scan(&at, &username, &subject, &mode, &rule, &action, &seconds, &detail); err != nil {
			return err
		}
		who := username
		if who == "" {
			who = subject // casual play has no account
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%.2f\t%s\n",
			at.Format(time.DateTime), who, mode, rule, action, seconds, detail)
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	w.Flush()

	if count == 0 {
		fmt.Println("no findings matched")
	}
	return nil
}

// ─── backup ────────────────────────────────────────────────────

func cmdBackup(args []string) error {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	target, err := openSubject(fs, args, "backup <destination.db>")
	defer database.Close()
	if err != nil {
		return err
	}
	dest, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("%s already exists", dest)
	}

	// VACUUM INTO takes a consistent snapshot of a live database. Copying the
	// file with cp while the server is writing would capture a torn WAL.
	if _, err := database.DB.ExecContext(context.Background(), "VACUUM INTO ?", dest); err != nil {
		return err
	}

	info, err := os.Stat(dest)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %s (%d bytes)\n", dest, info.Size())
	return nil
}

// ─── helpers ───────────────────────────────────────────────────

func table() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func sortedKeys(p anticheat.Params) []string {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	// Sorted so repeated runs print the same order; map order is random.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// isFlagSet reports whether a flag was given on the command line, which is how
// -enabled distinguishes "leave it alone" from "set it to false".
func isFlagSet(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// stdin is read through one shared buffered reader. A second bufio.Reader over
// os.Stdin would find nothing left: the first one buffers ahead, so piping two
// lines into a command that prompts twice gave EOF on the second prompt.
var stdin = bufio.NewReader(os.Stdin)

func prompt(text string) (string, error) {
	fmt.Print(text)
	line, err := stdin.ReadString('\n')
	return strings.TrimSpace(line), err
}

// readPassword prompts without echoing. When stdin is not a terminal — a pipe
// in a script, or `docker exec` without -t — it falls back to a plain read so
// the command still works unattended.
func readPassword(text string) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return prompt(text)
	}
	fmt.Print(text)
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(raw), err
}

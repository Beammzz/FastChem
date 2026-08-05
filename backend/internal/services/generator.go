package services

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/takumi/fastchem/internal/models"
)

// QuestionGen is the interface that all difficulty generators implement
type QuestionGen interface {
	Generate() models.Question
}

// QuestionGenerator is the main entry point – delegates to the correct difficulty generator
type QuestionGenerator struct {
	easy   QuestionGen
	medium QuestionGen
	hard   QuestionGen
}

// NewQuestionGenerator creates a new generator instance
func NewQuestionGenerator() *QuestionGenerator {
	return &QuestionGenerator{
		easy:   &EasyGenerator{},
		medium: &MediumGenerator{},
		hard:   &HardGenerator{},
	}
}

// GenerateQuestion produces a question for the given difficulty
func (g *QuestionGenerator) GenerateQuestion(difficulty string) models.Question {
	switch difficulty {
	case "medium":
		return g.medium.Generate()
	case "hard":
		return g.hard.Generate()
	default:
		return g.easy.Generate()
	}
}

// GenerateQuestionByCategories picks a random category from the list and generates a matching question
func (g *QuestionGenerator) GenerateQuestionByCategories(categories []string) models.Question {
	if len(categories) == 0 {
		return g.easy.Generate()
	}
	cat := categories[rand.Intn(len(categories))]
	switch cat {
	case "atomic_structure":
		return (&EasyGenerator{}).generateAtomicStructureQuestion()
	case "oxidation_number":
		return (&EasyGenerator{}).generateOxidationQuestion()
	case "state_of_matter":
		return (&EasyGenerator{}).generateStateOfMatterQuestion()
	case "mole_concept":
		return (&MediumGenerator{}).Generate()
	case "dilution":
		return (&HardGenerator{}).dilution()
	case "preparing_solution":
		return (&HardGenerator{}).preparingFromSolid()
	case "freezing_point":
		return (&HardGenerator{}).freezingPointDepression()
	default:
		return g.easy.Generate()
	}
}

// ValidateAnswer checks if the selected answer is correct (legacy)
func (g *QuestionGenerator) ValidateAnswer(req models.ValidateRequest) models.ValidateResponse {
	correct := req.SelectedIndex == req.CorrectIndex
	msg := "ผิด!"
	if correct {
		msg = "ถูกต้อง!"
	}
	return models.ValidateResponse{
		Correct:      correct,
		CorrectIndex: req.CorrectIndex,
		Message:      msg,
	}
}

// ─── Randomness helpers ─────────────────────────────────────────

// recentWindow is how many previous indices to remember per category
const recentWindow = 8

// recentMu guards the recent* slices below. Gin serves requests concurrently,
// so pickFreshIndex is called from many goroutines at once.
var (
	recentMu        sync.Mutex
	recentOxidation []int
	recentState     []int
)

// pickFreshIndex returns a random index in [0, n) that was not recently used.
// It mutates the recent slice to record the picked index.
func pickFreshIndex(n int, recent *[]int) int {
	if n <= 1 {
		return 0
	}
	recentMu.Lock()
	defer recentMu.Unlock()
	for {
		idx := rand.Intn(n)
		dup := false
		for _, r := range *recent {
			if r == idx {
				dup = true
				break
			}
		}
		if !dup {
			*recent = append(*recent, idx)
			if len(*recent) > recentWindow {
				*recent = (*recent)[1:]
			}
			return idx
		}
		// If the pool is smaller than the window, clear history to avoid deadloop
		if n <= len(*recent) {
			*recent = nil
		}
	}
}

// ─── Easy Generator ──────────────────────────────────────────────

// EasyGenerator produces atomic-structure, oxidation-number, and state-of-matter questions
type EasyGenerator struct{}

// Generate creates one easy question, cycling randomly among three topics
func (g *EasyGenerator) Generate() models.Question {
	switch rand.Intn(3) {
	case 0:
		return g.generateAtomicStructureQuestion()
	case 1:
		return g.generateOxidationQuestion()
	default:
		return g.generateStateOfMatterQuestion()
	}
}

func (g *EasyGenerator) generateAtomicStructureQuestion() models.Question {
	elem := elements[rand.Intn(len(elements))]
	atomicNumber := elem.AtomicNumber
	massNumber := atomicNumber + rand.Intn(11)

	questionType := rand.Intn(4)

	var questionText string
	var correctAnswer int

	switch questionType {
	case 0:
		// Ask for atomic number — no hint, must look up periodic table
		questionText = fmt.Sprintf(
			"ธาตุ %s (%s) มีเลขอะตอมเท่าใด?",
			elem.NameTH, elem.Symbol,
		)
		correctAnswer = atomicNumber
	case 1:
		// Ask for proton count — no hint
		questionText = fmt.Sprintf(
			"ธาตุ %s (%s) มีจำนวนโปรตอนเท่าใด?",
			elem.NameTH, elem.Symbol,
		)
		correctAnswer = atomicNumber
	case 2:
		// Ask for electron count in neutral atom — no hint
		questionText = fmt.Sprintf(
			"อะตอมที่เป็นกลางของ %s (%s) มีจำนวนอิเล็กตรอนเท่าใด?",
			elem.NameTH, elem.Symbol,
		)
		correctAnswer = atomicNumber
	case 3:
		// Ask for neutron count — give mass number only, must know atomic number
		neutrons := massNumber - atomicNumber
		questionText = fmt.Sprintf(
			"ไอโซโทปของ %s (%s) ที่มีเลขมวล %d มีจำนวนนิวตรอนเท่าใด?",
			elem.NameTH, elem.Symbol, massNumber,
		)
		correctAnswer = neutrons
	}

	choices, correctIndex := generateDistractors(correctAnswer, 0, 40)

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("easy").TimeLimit,
		Category:     "atomic_structure",
		Difficulty:   "easy",
	}
}

func (g *EasyGenerator) generateOxidationQuestion() models.Question {
	idx := pickFreshIndex(len(compounds), &recentOxidation)
	compound := compounds[idx]
	target := compound.Elements[rand.Intn(len(compound.Elements))]

	questionText := fmt.Sprintf(
		"เลขออกซิเดชันของ %s ใน %s คือเท่าใด?",
		target.Symbol, compound.Formula,
	)

	choices, correctIndex := generateOxidationDistractors(target.OxidationNumber)

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("easy").TimeLimit,
		Category:     "oxidation_number",
		Difficulty:   "easy",
	}
}

// stateLabels maps internal state codes to Thai display labels
var stateLabels = map[string]string{
	"s":  "ของแข็ง (s)",
	"l":  "ของเหลว (l)",
	"g":  "แก๊ส (g)",
	"aq": "สารละลายในน้ำ (aq)",
}

func (g *EasyGenerator) generateStateOfMatterQuestion() models.Question {
	idx := pickFreshIndex(len(stateCompounds), &recentState)
	compound := stateCompounds[idx]

	correctLabel := stateLabels[compound.State]

	// Build distractors: the 3 states that are not the correct one
	allStates := []string{"s", "l", "g", "aq"}
	distractors := make([]string, 0, 3)
	for _, s := range allStates {
		if s != compound.State {
			distractors = append(distractors, stateLabels[s])
		}
	}
	rand.Shuffle(len(distractors), func(i, j int) {
		distractors[i], distractors[j] = distractors[j], distractors[i]
	})

	choices := append(distractors[:3], correctLabel)
	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	correctIndex := 0
	for i, c := range choices {
		if c == correctLabel {
			correctIndex = i
			break
		}
	}

	questionText := fmt.Sprintf(
		"%s (%s) มีสถานะที่ภาวะมาตรฐาน (25°C, 1 atm) เป็นอย่างไร?",
		compound.NameTH, compound.Formula,
	)

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("easy").TimeLimit,
		Category:     "state_of_matter",
		Difficulty:   "easy",
	}
}

// ─── Shared helpers ──────────────────────────────────────────────

// generateDistractors creates 4 unique plausible choices for integer answers
func generateDistractors(correct, minVal, maxVal int) ([]string, int) {
	choices := make(map[int]bool)
	choices[correct] = true

	for len(choices) < 4 {
		offset := rand.Intn(9) - 4
		if offset == 0 {
			offset = rand.Intn(2)*2 - 1
		}
		distractor := correct + offset
		if distractor < minVal {
			distractor = minVal + rand.Intn(5)
		}
		if distractor > maxVal {
			distractor = maxVal - rand.Intn(5)
		}
		if distractor != correct {
			choices[distractor] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, strconv.Itoa(c))
	}

	rand.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctStr := strconv.Itoa(correct)
	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}

	return choiceSlice, correctIndex
}

// generateOxidationDistractors creates 4 unique plausible oxidation number choices
func generateOxidationDistractors(correct int) ([]string, int) {
	choices := make(map[int]bool)
	choices[correct] = true

	commonOx := []int{-3, -2, -1, 0, +1, +2, +3, +4, +5, +6, +7}

	for len(choices) < 4 {
		candidate := commonOx[rand.Intn(len(commonOx))]
		if candidate != correct {
			choices[candidate] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		if c > 0 {
			choiceSlice = append(choiceSlice, fmt.Sprintf("+%d", c))
		} else {
			choiceSlice = append(choiceSlice, strconv.Itoa(c))
		}
	}

	rand.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	var correctStr string
	if correct > 0 {
		correctStr = fmt.Sprintf("+%d", correct)
	} else {
		correctStr = strconv.Itoa(correct)
	}

	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}

	return choiceSlice, correctIndex
}

// distractorAttempts bounds the random phase of distractor generation.
// Past this, the deterministic ladder guarantees termination.
const distractorAttempts = 200

// formatPrecisionStep returns one unit in the last decimal place of a
// printf-style float format, e.g. "%.2f mol" → 0.01.
func formatPrecisionStep(format string) float64 {
	i := strings.Index(format, "%.")
	if i < 0 {
		return 0.01
	}
	prec := 0
	for j := i + 2; j < len(format) && format[j] >= '0' && format[j] <= '9'; j++ {
		prec = prec*10 + int(format[j]-'0')
	}
	return math.Pow(10, -float64(prec))
}

// fillFloatLadder tops choices up to 4 using values one format-step apart,
// walking away from zero so the sign of the answer is preserved. Each step is
// exactly one unit in the last printed decimal, so every k yields a distinct
// string — this always terminates, whatever the magnitude of correct.
func fillFloatLadder(choices map[string]bool, correct float64, format string) {
	step := formatPrecisionStep(format)
	correctStr := fmt.Sprintf(format, correct)
	dir := 1.0
	if correct < 0 {
		dir = -1.0
	}
	for k := 1; len(choices) < 4; k++ {
		v := math.Round((correct+dir*float64(k)*step)/step) * step
		if s := fmt.Sprintf(format, v); s != correctStr {
			choices[s] = true
		}
	}
}

// generateFloatDistractors creates 4 unique plausible float choices
func generateFloatDistractors(correct float64, format string) ([]string, int) {
	choices := make(map[string]bool)
	correctStr := fmt.Sprintf(format, correct)
	choices[correctStr] = true

	for attempt := 0; len(choices) < 4 && attempt < distractorAttempts; attempt++ {
		factor := 0.5 + rand.Float64()*1.0
		if factor > 0.9 && factor < 1.1 {
			factor = 1.3
		}
		distractor := correct * factor
		dStr := fmt.Sprintf(format, distractor)
		if dStr != correctStr {
			choices[dStr] = true
		}
	}
	// Small answers have too few distinct values within ±50% to fill 4 slots.
	fillFloatLadder(choices, correct, format)

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, c)
	}

	rand.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}

	return choiceSlice, correctIndex
}

"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Navbar from "@/components/Navbar";
import { useAuth } from "@/components/AuthProvider";
import AnalyticsPanel from "@/components/admin/AnalyticsPanel";
import AnticheatPanel from "@/components/admin/AnticheatPanel";
import LeaderboardPanel from "@/components/admin/LeaderboardPanel";
import LivePanel from "@/components/admin/LivePanel";
import QuestionsPanel from "@/components/admin/QuestionsPanel";
import UsersPanel from "@/components/admin/UsersPanel";
import { Spinner } from "@/components/admin/ui";
import {
  deleteAdminUser,
  fetchAdminAnalytics,
  fetchAdminFindings,
  fetchAdminLeaderboards,
  fetchAdminLive,
  fetchAdminMatches,
  fetchAdminOverview,
  fetchAdminRules,
  fetchAdminTopics,
  fetchAdminUsers,
  previewAdminQuestion,
  updateAdminRule,
  updateAdminUser,
} from "@/lib/api";
import {
  AdminAnalytics,
  AdminFindingsResponse,
  AdminLeaderboards,
  AdminLive,
  AdminMatchesResponse,
  AdminOverview,
  AdminRule,
  AdminTopic,
  AdminUserUpdate,
  AdminUsersResponse,
} from "@/types";

type Tab = "analytics" | "users" | "anticheat" | "live" | "leaderboard" | "questions";

const TABS: { id: Tab; label: string; icon: string }[] = [
  { id: "analytics", label: "ภาพรวม", icon: "📊" },
  { id: "users", label: "ผู้เล่น", icon: "👥" },
  { id: "anticheat", label: "Anti-cheat", icon: "🛡️" },
  { id: "live", label: "สด", icon: "🟢" },
  { id: "leaderboard", label: "กระดานผู้นำ", icon: "🏆" },
  { id: "questions", label: "คลังคำถาม", icon: "🧪" },
];

// The console is one page of tabs. Each tab owns its own fetch, so opening the
// admin page does not pull every table in the database at once.
export default function AdminPage() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const [tab, setTab] = useState<Tab>("analytics");
  const [error, setError] = useState("");

  // analytics
  const [overview, setOverview] = useState<AdminOverview | null>(null);
  const [analytics, setAnalytics] = useState<AdminAnalytics | null>(null);
  const [days, setDays] = useState(30);

  // users
  const [users, setUsers] = useState<AdminUsersResponse | null>(null);
  const [usersLoading, setUsersLoading] = useState(false);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState("points");
  const [userPage, setUserPage] = useState(1);

  // anticheat
  const [rules, setRules] = useState<AdminRule[]>([]);
  const [findings, setFindings] = useState<AdminFindingsResponse | null>(null);
  const [findingsLoading, setFindingsLoading] = useState(false);
  const [filters, setFilters] = useState({ rule: "", mode: "", action: "" });
  const [findingPage, setFindingPage] = useState(1);

  // live
  const [live, setLive] = useState<AdminLive | null>(null);
  const [matches, setMatches] = useState<AdminMatchesResponse | null>(null);
  const [liveLoading, setLiveLoading] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState(true);

  // leaderboard + questions
  const [boards, setBoards] = useState<AdminLeaderboards | null>(null);
  const [boardsLoading, setBoardsLoading] = useState(false);
  const [topics, setTopics] = useState<AdminTopic[]>([]);
  const [topicsLoading, setTopicsLoading] = useState(false);

  // The flag on the token holder is a hint for what to render; the server
  // rejects a non-admin regardless of what the client believes.
  useEffect(() => {
    if (!authLoading && !user) router.replace("/login/");
  }, [authLoading, user, router]);

  const report = (e: unknown) => setError(e instanceof Error ? e.message : "โหลดข้อมูลไม่สำเร็จ");

  useEffect(() => {
    if (!user?.isAdmin || tab !== "analytics") return;
    Promise.all([fetchAdminOverview(), fetchAdminAnalytics(days)])
      .then(([o, a]) => {
        setOverview(o);
        setAnalytics(a);
        setError("");
      })
      .catch(report);
  }, [user?.isAdmin, tab, days]);

  const loadUsers = useCallback(() => {
    setUsersLoading(true);
    fetchAdminUsers({ page: userPage, search, sort })
      .then(setUsers)
      .catch(report)
      .finally(() => setUsersLoading(false));
  }, [userPage, search, sort]);

  useEffect(() => {
    if (!user?.isAdmin || tab !== "users") return;
    // Debounced so typing in the search box does not fire a query per keystroke.
    const timer = setTimeout(loadUsers, 250);
    return () => clearTimeout(timer);
  }, [user?.isAdmin, tab, loadUsers]);

  const loadFindings = useCallback(() => {
    setFindingsLoading(true);
    fetchAdminFindings({ page: findingPage, ...filters })
      .then(setFindings)
      .catch(report)
      .finally(() => setFindingsLoading(false));
  }, [findingPage, filters]);

  useEffect(() => {
    if (!user?.isAdmin || tab !== "anticheat") return;
    fetchAdminRules().then(setRules).catch(report);
    loadFindings();
  }, [user?.isAdmin, tab, loadFindings]);

  const loadLive = useCallback(() => {
    setLiveLoading(true);
    Promise.all([fetchAdminLive(), fetchAdminMatches()])
      .then(([l, m]) => {
        setLive(l);
        setMatches(m);
      })
      .catch(report)
      .finally(() => setLiveLoading(false));
  }, []);

  useEffect(() => {
    if (!user?.isAdmin || tab !== "live") return;
    loadLive();
    if (!autoRefresh) return;
    const timer = setInterval(loadLive, 5000);
    return () => clearInterval(timer);
  }, [user?.isAdmin, tab, autoRefresh, loadLive]);

  useEffect(() => {
    if (!user?.isAdmin || tab !== "leaderboard") return;
    setBoardsLoading(true);
    fetchAdminLeaderboards()
      .then(setBoards)
      .catch(report)
      .finally(() => setBoardsLoading(false));
  }, [user?.isAdmin, tab]);

  useEffect(() => {
    if (!user?.isAdmin || tab !== "questions" || topics.length > 0) return;
    setTopicsLoading(true);
    fetchAdminTopics()
      .then(setTopics)
      .catch(report)
      .finally(() => setTopicsLoading(false));
  }, [user?.isAdmin, tab, topics.length]);

  const onUpdateUser = async (update: AdminUserUpdate) => {
    await updateAdminUser(update);
    loadUsers();
  };

  const onDeleteUser = async (userId: number) => {
    await deleteAdminUser(userId);
    loadUsers();
  };

  const onSaveRule = async (rule: {
    name: string;
    enabled: boolean;
    action: string;
    params: Record<string, number>;
  }) => {
    await updateAdminRule(rule);
    setRules(await fetchAdminRules());
  };

  if (authLoading) {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
        <Navbar />
        <Spinner />
      </div>
    );
  }

  if (!user?.isAdmin) {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
        <Navbar />
        <div className="max-w-md mx-auto px-4 py-24 text-center">
          <div className="text-4xl mb-4">🔒</div>
          <h1 className="text-xl font-bold mb-2">ไม่มีสิทธิ์เข้าถึง</h1>
          <p className="text-sm text-gray-500">
            หน้านี้สำหรับผู้ดูแลระบบเท่านั้น ตั้งค่า ADMIN_USERNAMES
            แล้วรีสตาร์ตเซิร์ฟเวอร์เพื่อให้สิทธิ์บัญชีแรก
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 pt-6 sm:pt-10 pb-16 pb-safe">
        <header className="mb-6">
          <h1 className="text-2xl sm:text-3xl font-extrabold">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
              ระบบผู้ดูแล
            </span>
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            เข้าสู่ระบบเป็น {user.username} · ทุกคำสั่งถูกตรวจสิทธิ์ที่เซิร์ฟเวอร์ทุกครั้ง
          </p>
        </header>

        <nav className="flex gap-1.5 overflow-x-auto pb-2 mb-4">
          {TABS.map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`px-3 sm:px-4 py-2 rounded-lg text-xs sm:text-sm font-medium whitespace-nowrap transition-all cursor-pointer ${
                tab === t.id
                  ? "bg-violet-500/20 text-violet-300 border border-violet-500/30"
                  : "text-gray-500 hover:text-gray-300 border border-transparent"
              }`}
            >
              <span className="mr-1.5">{t.icon}</span>
              {t.label}
            </button>
          ))}
        </nav>

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 rounded-xl px-4 py-3 text-sm text-red-300 mb-4">
            {error}
          </div>
        )}

        {tab === "analytics" && (
          <AnalyticsPanel
            overview={overview}
            analytics={analytics}
            days={days}
            onDaysChange={setDays}
          />
        )}

        {tab === "users" && (
          <UsersPanel
            data={users}
            loading={usersLoading}
            search={search}
            sort={sort}
            page={userPage}
            currentUserId={user.id}
            onSearch={(value) => {
              setSearch(value);
              setUserPage(1);
            }}
            onSort={(key) => {
              setSort(key);
              setUserPage(1);
            }}
            onPage={setUserPage}
            onUpdate={onUpdateUser}
            onDelete={onDeleteUser}
          />
        )}

        {tab === "anticheat" && (
          <AnticheatPanel
            rules={rules}
            findings={findings}
            loading={findingsLoading}
            filters={filters}
            page={findingPage}
            onFilter={(next) => {
              setFilters(next);
              setFindingPage(1);
            }}
            onPage={setFindingPage}
            onSaveRule={onSaveRule}
          />
        )}

        {tab === "live" && (
          <LivePanel
            live={live}
            matches={matches}
            loading={liveLoading}
            autoRefresh={autoRefresh}
            onAutoRefresh={setAutoRefresh}
          />
        )}

        {tab === "leaderboard" && <LeaderboardPanel boards={boards} loading={boardsLoading} />}

        {tab === "questions" && (
          <QuestionsPanel
            topics={topics}
            loading={topicsLoading}
            onPreview={previewAdminQuestion}
          />
        )}
      </div>
    </div>
  );
}

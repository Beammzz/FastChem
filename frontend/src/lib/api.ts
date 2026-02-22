const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

import {
  Question,
  ValidateRequest,
  ValidateResponse,
  AnswerRequest,
  AnswerResponse,
  AuthResponse,
  MeResponse,
  LeaderboardResponse,
  ProfileResponse,
  StartMatchResponse,
  MatchAnswerRequest,
  MatchAnswerResponse,
  EndMatchResponse,
  RankedStats,
  RankedMatchHistoryEntry,
  RankedLeaderboardEntry,
} from "@/types";

// Helper to get auth headers
function authHeaders(): Record<string, string> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  if (token) {
    return { Authorization: `Bearer ${token}` };
  }
  return {};
}

export async function fetchQuestion(
  difficulty: string = "easy",
  categories?: string[]
): Promise<Question> {
  const params = new URLSearchParams();
  params.set("difficulty", difficulty);
  if (categories && categories.length > 0) {
    params.set("categories", categories.join(","));
  }
  const res = await fetch(`${API_BASE}/api/question?${params.toString()}`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("Failed to fetch question");
  }
  return res.json();
}

export async function validateAnswer(
  req: ValidateRequest
): Promise<ValidateResponse> {
  const res = await fetch(`${API_BASE}/api/validate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    throw new Error("Failed to validate answer");
  }
  return res.json();
}

// New per-question answer submission (server-side scoring + anti-cheat)
export async function submitAnswer(
  req: AnswerRequest
): Promise<AnswerResponse> {
  const res = await fetch(`${API_BASE}/api/answer`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    throw new Error("Failed to submit answer");
  }
  return res.json();
}

// Auth API
export async function register(
  username: string,
  password: string
): Promise<AuthResponse> {
  const res = await fetch(`${API_BASE}/api/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Registration failed");
  }
  return res.json();
}

export async function login(
  username: string,
  password: string
): Promise<AuthResponse> {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Login failed");
  }
  return res.json();
}

export async function fetchMe(): Promise<MeResponse> {
  const res = await fetch(`${API_BASE}/api/auth/me`, {
    headers: { ...authHeaders() },
  });
  if (!res.ok) {
    throw new Error("Not authenticated");
  }
  return res.json();
}

// Leaderboard API
export async function fetchLeaderboard(): Promise<LeaderboardResponse> {
  const res = await fetch(`${API_BASE}/api/leaderboard`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("Failed to fetch leaderboard");
  }
  return res.json();
}

// Score API
export async function submitScore(data: {
  score: number;
  totalAnswered: number;
  correctAnswers: number;
  difficulty: string;
  timeLimit: number;
  timeSpent: number;
}): Promise<{ id: number; message: string; totalPoints: number }> {
  const res = await fetch(`${API_BASE}/api/scores`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    throw new Error("Failed to submit score");
  }
  return res.json();
}

// Profile API
export async function fetchProfile(
  username: string,
  page: number = 1
): Promise<ProfileResponse> {
  const res = await fetch(
    `${API_BASE}/api/profile/${encodeURIComponent(username)}?page=${page}`,
    { cache: "no-store" }
  );
  if (!res.ok) {
    throw new Error("Failed to fetch profile");
  }
  return res.json();
}

// ─── Match API ─────────────────────────────────────────────────

export async function startMatch(
  difficulty: string
): Promise<StartMatchResponse> {
  const res = await fetch(`${API_BASE}/api/match/start`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify({ difficulty }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Failed to start match");
  }
  return res.json();
}

export async function submitMatchAnswer(
  req: MatchAnswerRequest
): Promise<MatchAnswerResponse> {
  const res = await fetch(`${API_BASE}/api/match/answer`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Failed to submit answer");
  }
  return res.json();
}

export async function endMatch(matchId: number): Promise<EndMatchResponse> {
  const res = await fetch(`${API_BASE}/api/match/end`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
    },
    body: JSON.stringify({ matchId }),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.error || "Failed to end match");
  }
  return res.json();
}

// ─── Ranked API ────────────────────────────────────────────────

export function getRankedWebSocketURL(): string {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;
  // Derive WebSocket URL from the current page origin (works for both dev and prod)
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const host = window.location.host;
  return `${protocol}//${host}/api/ranked/ws?token=${token || ""}`;
}

export async function fetchRankedStats(): Promise<RankedStats> {
  const res = await fetch(`${API_BASE}/api/ranked/stats`, {
    headers: { ...authHeaders() },
  });
  if (!res.ok) {
    throw new Error("Failed to fetch ranked stats");
  }
  return res.json();
}

export async function fetchRankedHistory(): Promise<RankedMatchHistoryEntry[]> {
  const res = await fetch(`${API_BASE}/api/ranked/history`, {
    headers: { ...authHeaders() },
  });
  if (!res.ok) {
    throw new Error("Failed to fetch ranked history");
  }
  const data = await res.json();
  return data.history;
}

export async function fetchRankedLeaderboard(): Promise<
  RankedLeaderboardEntry[]
> {
  const res = await fetch(`${API_BASE}/api/ranked/leaderboard`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("Failed to fetch ranked leaderboard");
  }
  const data = await res.json();
  return data.entries;
}

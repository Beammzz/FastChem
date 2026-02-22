export interface Question {
  id: string;
  question: string;
  choices: string[];
  correctIndex: number;
  timeLimit: number;
  category: string;
  difficulty: string;
}

export interface ValidateRequest {
  questionId: string;
  selectedIndex: number;
  correctIndex: number;
}

export interface ValidateResponse {
  correct: boolean;
  correctIndex: number;
  message: string;
}

// New per-question answer submission
export interface AnswerRequest {
  questionId: string;
  selectedIndex: number;
}

export interface AnswerResponse {
  correct: boolean;
  correctIndex: number;
  scoreEarned: number;
  speedBonus: number;
  timeSpent: number;
  difficulty: string;
  timedOut: boolean;
}

// ─── Match types ─────────────────────────────────────────────────

export interface StartMatchRequest {
  difficulty: string;
}

export interface StartMatchResponse {
  matchId: number;
  difficulty: string;
  questions: number;
  timeLimit: number;
}

export interface MatchAnswerRequest {
  matchId: number;
  questionId: string;
  selectedIndex: number;
}

export interface MatchAnswerResponse {
  correct: boolean;
  correctIndex: number;
  scoreEarned: number;
  speedBonus: number;
  comboMultiplier: number;
  combo: number;
  timeSpent: number;
  timedOut: boolean;
  totalScore: number;
  questionNumber: number;
}

export interface EndMatchResponse {
  matchId: number;
  totalScore: number;
  totalAnswered: number;
  correctAnswers: number;
  bestCombo: number;
  difficulty: string;
  attempts: QuestionAttemptResult[];
  totalPoints: number;
}

export interface QuestionAttemptResult {
  correct: boolean;
  timedOut: boolean;
  timeSpentMs: number;
  scoreAwarded: number;
  comboAtAnswer: number;
  comboMultiplier: number;
}

// ─── Game types ──────────────────────────────────────────────────

export type GameState = "idle" | "playing" | "finished";

// Difficulty scoring configuration (mirrors backend)
export interface DifficultyConfig {
  timeLimit: number;
  baseScore: number;
  multiplier: number;
}

export const DIFFICULTY_CONFIGS: Record<string, DifficultyConfig> = {
  easy: { timeLimit: 30, baseScore: 5, multiplier: 1 },
  medium: { timeLimit: 60, baseScore: 10, multiplier: 2 },
  hard: { timeLimit: 120, baseScore: 20, multiplier: 3 },
};

// Combo multiplier tiers (mirrors backend)
export function getComboMultiplier(combo: number): number {
  if (combo >= 10) return 2.0;
  if (combo >= 6) return 1.5;
  if (combo >= 3) return 1.2;
  return 1.0;
}

// Auth types
export interface User {
  id: number;
  username: string;
  totalPoints?: number;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface UserStats {
  totalGames: number;
  totalPoints: number;
  highScore: number;
  totalCorrect: number;
  totalAnswered: number;
  averageScore: number;
  averageAccuracy: number;
  fastestTime: number;
  averageTime: number;
}

export interface MeResponse {
  user: User;
  stats: UserStats;
}

// Leaderboard types
export interface LeaderboardEntry {
  rank: number;
  username: string;
  userId: number;
  totalPoints: number;
  totalGames: number;
}

export interface LeaderboardResponse {
  entries: LeaderboardEntry[];
}

export interface ScoreEntry {
  id: number;
  userId: number;
  score: number;
  totalAnswered: number;
  correctAnswers: number;
  difficulty: string;
  timeLimit: number;
  timeSpent: number;
  playedAt: string;
}

export interface ProfileResponse {
  user: User;
  stats: UserStats;
  history: ScoreEntry[];
  page: number;
  total: number;
}

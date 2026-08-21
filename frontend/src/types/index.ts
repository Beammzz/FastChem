// No correctIndex: the server withholds the answer until the question is
// answered, then returns it on AnswerResponse / MatchAnswerResponse.
export interface Question {
  id: string;
  question: string;
  choices: string[];
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
  // Set only by /api/auth/me. Decides whether the console is offered; it
  // grants nothing — every admin route re-checks the flag server-side.
  isAdmin?: boolean;
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

// ─── Ranked types ────────────────────────────────────────────────

export type RankedGamePhase =
  | "idle"
  | "queuing"
  | "matched"
  | "playing"
  | "finished"
  | "error";

// WebSocket event types (mirror backend)
export type WSEventType =
  | "MATCH_FOUND"
  | "MATCH_START"
  | "QUESTION_START"
  | "SUBMIT_ANSWER"
  | "ANSWER_RESULT"
  | "MATCH_END"
  | "OPPONENT_PROGRESS"
  | "ERROR"
  | "PING"
  | "PONG"
  | "QUEUE_JOINED"
  | "ROOM_CREATED"
  | "ROOM_JOINED";

export interface WSMessage {
  type: WSEventType;
  payload?: unknown;
}

export interface MatchFoundPayload {
  matchId: number;
  opponentId: number;
  opponent: string;
}

export interface MatchStartPayload {
  matchId: number;
  totalQuestions: number;
}

export interface QuestionStartPayload {
  index: number;
  question: string;
  choices: string[];
  timeLimit: number;
  difficulty: string;
  category: string;
}

export interface SubmitAnswerPayload {
  index: number;
  selectedIndex: number;
}

export interface AnswerResultPayload {
  index: number;
  correct: boolean;
  correctIndex: number;
  scoreEarned: number;
  speedBonus: number;
  comboMultiplier: number;
  combo: number;
  totalScore: number;
  timedOut: boolean;
  timeSpent: number;
  opponentAnswered: number;
  opponentScore: number;
}

export interface MatchEndPayload {
  matchId: number;
  winner: string;
  winnerId: number;
  player1: MatchEndPlayerSummary;
  player2: MatchEndPlayerSummary;
  ratingChange: number;
  newRating: number;
}

export interface MatchEndPlayerSummary {
  userId: number;
  username: string;
  totalScore: number;
  correctAnswers: number;
  hardScore: number;
  totalTime: number;
  bestCombo: number;
}

export interface OpponentProgressPayload {
  answered: number;
  totalScore: number;
}

export interface RankedStats {
  rating: number;
  rankedWins: number;
  rankedLosses: number;
  highestRating: number;
}

export interface RankedMatchHistoryEntry {
  matchId: number;
  player1Id: number;
  player2Id: number;
  player1Score: number;
  player2Score: number;
  winnerId: number | null;
  status: string;
  createdAt: string;
  finishedAt: string | null;
  player1Name: string;
  player2Name: string;
  won: boolean;
}

export interface RankedLeaderboardEntry {
  rank: number;
  username: string;
  userId: number;
  rating: number;
  wins: number;
  losses: number;
  highestRating: number;
}

// ─── Admin types ─────────────────────────────────────────────────
// Mirrors backend/internal/models/admin.go.

export interface AdminBucket {
  label: string;
  count: number;
  value: number;
}

export interface AdminOverview {
  users: number;
  admins: number;
  newUsers7d: number;
  activePlayers7d: number;
  games: number;
  games24h: number;
  soloMatches: number;
  rankedMatches: number;
  questionsAnswered: number;
  totalPoints: number;
  averageAccuracy: number;
  findings: number;
  findings24h: number;
  findingsDropped: number;
  rejectingRules: number;
  queueSize: number;
  liveMatches: number;
  openRooms: number;
  liveSoloMatches: number;
  storedQuestions: number;
  topics: number;
  uptimeSeconds: number;
  databaseBytes: number;
  goroutines: number;
}

export interface AdminDailyPoint {
  date: string;
  games: number;
  players: number;
  points: number;
  rankedMatches: number;
  findings: number;
}

export interface AdminTopicStat {
  topic: string;
  difficulty: string;
  attempts: number;
  correct: number;
  accuracy: number;
  averageTime: number;
}

export interface AdminAnalytics {
  daily: AdminDailyPoint[];
  difficulty: AdminBucket[];
  topics: AdminTopicStat[];
  hours: AdminBucket[];
  ratings: AdminBucket[];
  rules: AdminBucket[];
}

export interface AdminUserRow {
  id: number;
  username: string;
  isAdmin: boolean;
  totalPoints: number;
  rating: number;
  rankedWins: number;
  rankedLosses: number;
  highestRating: number;
  games: number;
  highScore: number;
  totalAnswered: number;
  totalCorrect: number;
  accuracy: number;
  findings: number;
  createdAt: string;
  lastPlayed: string | null;
  online: boolean;
}

export interface AdminUsersResponse {
  users: AdminUserRow[];
  total: number;
  page: number;
  pageSize: number;
}

export interface AdminUserUpdate {
  userId: number;
  isAdmin?: boolean;
  rating?: number;
  totalPoints?: number;
  password?: string;
}

export interface AdminFinding {
  id: number;
  at: string;
  subject: string;
  userId: number;
  username: string;
  matchId: number;
  mode: string;
  questionId: string;
  difficulty: string;
  timeSpent: number;
  rule: string;
  detail: string;
  action: string;
}

export interface AdminFindingsResponse {
  findings: AdminFinding[];
  total: number;
  page: number;
  pageSize: number;
  byRule: AdminBucket[];
  byMode: AdminBucket[];
  byAction: AdminBucket[];
  dropped: number;
}

export interface AdminRule {
  name: string;
  enabled: boolean;
  action: string;
  params: Record<string, number>;
  findings: number;
  findings24h: number;
}

export interface AdminQueueEntry {
  userId: number;
  username: string;
  rating: number;
  waitingSeconds: number;
}

export interface AdminLivePlayer {
  userId: number;
  username: string;
  rating: number;
  totalScore: number;
  answered: number;
  correct: number;
  combo: number;
  connected: boolean;
}

export interface AdminLiveMatch {
  matchId: number;
  status: string;
  question: number;
  createdAt: string;
  player1: AdminLivePlayer;
  player2: AdminLivePlayer;
}

export interface AdminRoom {
  code: string;
  hostId: number;
  hostName: string;
  guestId: number;
  guestName: string;
  matchId: number;
  status: string;
  createdAt: string;
}

export interface AdminLive {
  queue: AdminQueueEntry[];
  matches: AdminLiveMatch[];
  rooms: AdminRoom[];
  soloMatches: number;
  storedQuestions: number;
}

export interface AdminRankedMatchRow {
  matchId: number;
  player1: string;
  player2: string;
  player1Score: number;
  player2Score: number;
  winner: string;
  status: string;
  createdAt: string;
  finishedAt: string | null;
}

export interface AdminSoloMatchRow {
  matchId: number;
  username: string;
  difficulty: string;
  totalScore: number;
  bestCombo: number;
  attempts: number;
  startedAt: string;
  endedAt: string | null;
}

export interface AdminScoreRow {
  id: number;
  userId: number;
  username: string;
  score: number;
  totalAnswered: number;
  correctAnswers: number;
  difficulty: string;
  timeSpent: number;
  playedAt: string;
}

export interface AdminMatchesResponse {
  ranked: AdminRankedMatchRow[];
  solo: AdminSoloMatchRow[];
  scores: AdminScoreRow[];
}

export interface AdminLeaderboards {
  casual: LeaderboardEntry[];
  ranked: RankedLeaderboardEntry[];
}

export interface AdminTopic {
  category: string;
  difficulty: string;
}

export interface AdminQuestionPreview {
  question: string;
  choices: string[];
  correctIndex: number;
  category: string;
  difficulty: string;
  timeLimit: number;
}

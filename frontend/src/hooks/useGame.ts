"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import {
  Question,
  GameState,
  MatchAnswerResponse,
  DIFFICULTY_CONFIGS,
  EndMatchResponse,
} from "@/types";
import {
  fetchQuestion,
  submitAnswer,
  startMatch,
  submitMatchAnswer,
  endMatch,
} from "@/lib/api";
import { triggerHaptic } from "@/lib/haptic";

interface UseGameOptions {
  totalQuestions: number;
  difficulty: string;
  categories?: string[];
  /** Pass true when the user is logged in to use the match API */
  useMatchApi?: boolean;
}

interface GameHook {
  gameState: GameState;
  question: Question | null;
  score: number;
  timeLeft: number;
  timeLimit: number;
  totalAnswered: number;
  correctAnswers: number;
  totalQuestions: number;
  selectedIndex: number | null;
  isCorrect: boolean | null;
  /** Which choice was right, once the server has told us. Null until then. */
  revealedIndex: number | null;
  loading: boolean;
  combo: number;
  bestCombo: number;
  lastScoreEarned: number;
  lastSpeedBonus: number;
  lastComboMultiplier: number;
  difficulty: string;
  matchResult: EndMatchResponse | null;
  startGame: () => void;
  answerQuestion: (index: number) => void;
}

export function useGame({
  totalQuestions,
  difficulty,
  categories,
  useMatchApi = false,
}: UseGameOptions): GameHook {
  const cfg = DIFFICULTY_CONFIGS[difficulty] || DIFFICULTY_CONFIGS.easy;
  const questionTimeLimit = cfg.timeLimit;

  const [gameState, setGameState] = useState<GameState>("idle");
  const [question, setQuestion] = useState<Question | null>(null);
  const [score, setScore] = useState(0);
  const [timeLeft, setTimeLeft] = useState(questionTimeLimit);
  const [totalAnswered, setTotalAnswered] = useState(0);
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [isCorrect, setIsCorrect] = useState<boolean | null>(null);
  const [revealedIndex, setRevealedIndex] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [combo, setCombo] = useState(0);
  const [bestCombo, setBestCombo] = useState(0);
  const [lastScoreEarned, setLastScoreEarned] = useState(0);
  const [lastSpeedBonus, setLastSpeedBonus] = useState(0);
  const [lastComboMultiplier, setLastComboMultiplier] = useState(1);
  const [matchResult, setMatchResult] = useState<EndMatchResponse | null>(null);

  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const questionNumberRef = useRef(0);
  const answeringRef = useRef(false);
  const matchIdRef = useRef<number | null>(null);

  const getTimeLimit = useCallback(() => {
    if (question?.timeLimit) return question.timeLimit;
    return questionTimeLimit;
  }, [question, questionTimeLimit]);

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const loadQuestion = useCallback(async () => {
    setLoading(true);
    setSelectedIndex(null);
    setIsCorrect(null);
    setRevealedIndex(null);
    setLastScoreEarned(0);
    setLastSpeedBonus(0);
    setLastComboMultiplier(1);
    answeringRef.current = false;

    try {
      const q = await fetchQuestion(difficulty, categories);
      setQuestion(q);
      const limit = q.timeLimit || questionTimeLimit;
      setTimeLeft(limit);
    } catch (err) {
      console.error("Failed to load question:", err);
    } finally {
      setLoading(false);
    }
  }, [difficulty, categories, questionTimeLimit]);

  const startGame = useCallback(async () => {
    clearTimer();
    setGameState("playing");
    setScore(0);
    setTotalAnswered(0);
    setCorrectAnswers(0);
    setSelectedIndex(null);
    setIsCorrect(null);
    setRevealedIndex(null);
    setCombo(0);
    setBestCombo(0);
    setLastScoreEarned(0);
    setLastSpeedBonus(0);
    setLastComboMultiplier(1);
    setMatchResult(null);
    questionNumberRef.current = 0;
    answeringRef.current = false;
    matchIdRef.current = null;

    // If authenticated, start a match via API
    if (useMatchApi) {
      try {
        const matchRes = await startMatch(difficulty);
        matchIdRef.current = matchRes.matchId;
      } catch (err) {
        console.error("Failed to start match, fallback to legacy:", err);
        // fallback – continue without match API
      }
    }

    loadQuestion();
  }, [loadQuestion, clearTimer, useMatchApi, difficulty]);

  // End game and finalize match
  const finishGame = useCallback(async () => {
    if (matchIdRef.current) {
      try {
        const result = await endMatch(matchIdRef.current);
        setMatchResult(result);
        // Update local state with server-validated values
        setScore(result.totalScore);
        setCorrectAnswers(result.correctAnswers);
        setBestCombo(result.bestCombo);
      } catch {
        // Result can't be fetched, keep local values
      }
      matchIdRef.current = null;
    }
    setGameState("finished");
  }, []);

  // Handle timeout — called when timer reaches 0
  const handleTimeout = useCallback(async () => {
    if (answeringRef.current || !question) return;
    answeringRef.current = true;
    clearTimer();

    setSelectedIndex(-1);
    setIsCorrect(false);
    triggerHaptic("timeout");
    setCombo(0);
    setTotalAnswered((prev) => prev + 1);
    questionNumberRef.current += 1;

    // Submit timeout to backend. Its reply is also how we learn which choice
    // was right — the question itself never carried the answer.
    if (matchIdRef.current) {
      try {
        const result: MatchAnswerResponse = await submitMatchAnswer({
          matchId: matchIdRef.current,
          questionId: question.id,
          selectedIndex: -1,
        });
        setScore(result.totalScore);
        setCombo(result.combo);
        setRevealedIndex(result.correctIndex);
      } catch {
        // best-effort
      }
    } else {
      try {
        const result = await submitAnswer({
          questionId: question.id,
          selectedIndex: -1,
        });
        setRevealedIndex(result.correctIndex);
      } catch {
        // best-effort
      }
    }

    if (questionNumberRef.current >= totalQuestions) {
      setTimeout(() => finishGame(), 800);
    } else {
      setTimeout(() => loadQuestion(), 800);
    }
  }, [question, totalQuestions, clearTimer, loadQuestion, finishGame]);

  const answerQuestion = useCallback(
    async (index: number) => {
      if (!question || selectedIndex !== null || answeringRef.current) return;
      answeringRef.current = true;
      clearTimer();

      setSelectedIndex(index);
      setTotalAnswered((prev) => prev + 1);
      questionNumberRef.current += 1;

      // ─── Match API path ────────────────────────────────────
      if (matchIdRef.current) {
        try {
          const result: MatchAnswerResponse = await submitMatchAnswer({
            matchId: matchIdRef.current,
            questionId: question.id,
            selectedIndex: index,
          });

          setIsCorrect(result.correct);
          setRevealedIndex(result.correctIndex);
          triggerHaptic(result.correct ? "correct" : "wrong");
          setScore(result.totalScore);
          setCombo(result.combo);
          setBestCombo((prev) => Math.max(prev, result.combo));
          setLastScoreEarned(result.scoreEarned);
          setLastSpeedBonus(result.speedBonus);
          setLastComboMultiplier(result.comboMultiplier);

          if (result.correct) {
            setCorrectAnswers((prev) => prev + 1);
          }
        } catch {
          // Only the server knows the answer, so a failed submit leaves this
          // question unresolved rather than guessed at. endMatch reconciles
          // the score from what the server actually recorded.
          setIsCorrect(null);
        }
      } else {
        // ─── Legacy (unauthenticated) path ─────────────────
        try {
          const result = await submitAnswer({
            questionId: question.id,
            selectedIndex: index,
          });

          setIsCorrect(result.correct);
          setRevealedIndex(result.correctIndex);
          triggerHaptic(result.correct ? "correct" : "wrong");

          if (result.correct) {
            setScore((prev) => prev + result.scoreEarned);
            setCorrectAnswers((prev) => prev + 1);
            setLastScoreEarned(result.scoreEarned);
            setLastSpeedBonus(result.speedBonus);
            setCombo((prev) => {
              const newCombo = prev + 1;
              setBestCombo((best) => Math.max(best, newCombo));
              return newCombo;
            });
          } else {
            setLastScoreEarned(0);
            setLastSpeedBonus(0);
            setCombo(0);
          }
        } catch {
          // Same as above — without a reply there is nothing to mark against,
          // and the question passes unscored.
          setIsCorrect(null);
        }
      }

      if (questionNumberRef.current >= totalQuestions) {
        setTimeout(() => finishGame(), 800);
      } else {
        setTimeout(() => loadQuestion(), 800);
      }
    },
    [question, selectedIndex, totalQuestions, clearTimer, loadQuestion, finishGame]
  );

  // Per-question countdown timer
  useEffect(() => {
    if (gameState !== "playing" || loading || !question) return;

    clearTimer();
    timerRef.current = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          handleTimeout();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearTimer();
  }, [gameState, loading, question, clearTimer, handleTimeout]);

  // Cleanup on unmount
  useEffect(() => {
    return () => clearTimer();
  }, [clearTimer]);

  return {
    gameState,
    question,
    score,
    timeLeft,
    timeLimit: getTimeLimit(),
    totalAnswered,
    correctAnswers,
    totalQuestions,
    selectedIndex,
    isCorrect,
    revealedIndex,
    loading,
    combo,
    bestCombo,
    lastScoreEarned,
    lastSpeedBonus,
    lastComboMultiplier,
    difficulty,
    matchResult,
    startGame,
    answerQuestion,
  };
}

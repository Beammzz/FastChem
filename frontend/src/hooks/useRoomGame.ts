"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import {
  WSMessage,
  MatchFoundPayload,
  MatchStartPayload,
  QuestionStartPayload,
  AnswerResultPayload,
  MatchEndPayload,
  OpponentProgressPayload,
} from "@/types";
import { getRoomWebSocketURL } from "@/lib/api";

export type RoomPhase =
  | "idle"
  | "creating"
  | "waiting"    // host waiting for guest
  | "joining"
  | "matched"
  | "playing"
  | "finished"
  | "error";

interface UseRoomGameReturn {
  phase: RoomPhase;
  roomCode: string;
  // Match info
  matchId: number | null;
  opponentName: string;
  totalQuestions: number;
  // Current question
  questionIndex: number;
  questionText: string;
  choices: string[];
  timeLimit: number;
  timeLeft: number;
  difficulty: string;
  category: string;
  // Player state
  score: number;
  combo: number;
  correctAnswers: number;
  answeredCount: number;
  // Last answer result
  selectedIndex: number | null;
  isCorrect: boolean | null;
  lastScoreEarned: number;
  lastSpeedBonus: number;
  lastComboMultiplier: number;
  showingResult: boolean;
  // Opponent state
  opponentScore: number;
  opponentAnswered: number;
  // Match end
  matchResult: MatchEndPayload | null;
  // Error
  errorMessage: string;
  // Actions
  createRoom: () => void;
  joinRoom: (code: string) => void;
  submitAnswer: (selectedIndex: number) => void;
  leaveGame: () => void;
}

export function useRoomGame(): UseRoomGameReturn {
  const [phase, setPhase] = useState<RoomPhase>("idle");
  const [roomCode, setRoomCode] = useState("");
  const [matchId, setMatchId] = useState<number | null>(null);
  const [opponentName, setOpponentName] = useState("");
  const [totalQuestions, setTotalQuestions] = useState(10);

  // Question state
  const [questionIndex, setQuestionIndex] = useState(0);
  const [questionText, setQuestionText] = useState("");
  const [choices, setChoices] = useState<string[]>([]);
  const [timeLimit, setTimeLimit] = useState(30);
  const [timeLeft, setTimeLeft] = useState(30);
  const [difficulty, setDifficulty] = useState("");
  const [category, setCategory] = useState("");

  // Player state
  const [score, setScore] = useState(0);
  const [combo, setCombo] = useState(0);
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [answeredCount, setAnsweredCount] = useState(0);

  // Last answer
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [isCorrect, setIsCorrect] = useState<boolean | null>(null);
  const [lastScoreEarned, setLastScoreEarned] = useState(0);
  const [lastSpeedBonus, setLastSpeedBonus] = useState(0);
  const [lastComboMultiplier, setLastComboMultiplier] = useState(1);
  const [showingResult, setShowingResult] = useState(false);

  // Opponent
  const [opponentScore, setOpponentScore] = useState(0);
  const [opponentAnswered, setOpponentAnswered] = useState(0);

  // Match end
  const [matchResult, setMatchResult] = useState<MatchEndPayload | null>(null);

  // Error
  const [errorMessage, setErrorMessage] = useState("");

  const wsRef = useRef<WebSocket | null>(null);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const pingRef = useRef<NodeJS.Timeout | null>(null);
  const answeringRef = useRef(false);

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const clearPing = useCallback(() => {
    if (pingRef.current) {
      clearInterval(pingRef.current);
      pingRef.current = null;
    }
  }, []);

  const cleanup = useCallback(() => {
    clearTimer();
    clearPing();
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  }, [clearTimer, clearPing]);

  const startTimer = useCallback(
    (limit: number) => {
      clearTimer();
      setTimeLeft(limit);
      timerRef.current = setInterval(() => {
        setTimeLeft((prev) => {
          if (prev <= 1) {
            clearTimer();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    },
    [clearTimer]
  );

  const handleMessage = useCallback(
    (msg: WSMessage) => {
      switch (msg.type) {
        case "ROOM_CREATED": {
          const payload = msg.payload as { code: string; message: string };
          setRoomCode(payload.code);
          setPhase("waiting");
          break;
        }

        case "ROOM_JOINED": {
          const payload = msg.payload as { code: string; message: string };
          setRoomCode(payload.code);
          // Stay in joining phase until match found
          break;
        }

        case "MATCH_FOUND": {
          const payload = msg.payload as MatchFoundPayload;
          setMatchId(payload.matchId);
          setOpponentName(payload.opponent);
          setPhase("matched");
          break;
        }

        case "MATCH_START": {
          const payload = msg.payload as MatchStartPayload;
          setTotalQuestions(payload.totalQuestions);
          setPhase("playing");
          setScore(0);
          setCombo(0);
          setCorrectAnswers(0);
          setAnsweredCount(0);
          setOpponentScore(0);
          setOpponentAnswered(0);
          break;
        }

        case "QUESTION_START": {
          const payload = msg.payload as QuestionStartPayload;
          setQuestionIndex(payload.index);
          setQuestionText(payload.question);
          setChoices(payload.choices);
          setTimeLimit(payload.timeLimit);
          setDifficulty(payload.difficulty);
          setCategory(payload.category);
          setSelectedIndex(null);
          setIsCorrect(null);
          setShowingResult(false);
          setLastScoreEarned(0);
          setLastSpeedBonus(0);
          setLastComboMultiplier(1);
          answeringRef.current = false;
          startTimer(payload.timeLimit);
          break;
        }

        case "ANSWER_RESULT": {
          const payload = msg.payload as AnswerResultPayload;
          clearTimer();
          setIsCorrect(payload.correct);
          setLastScoreEarned(payload.scoreEarned);
          setLastSpeedBonus(payload.speedBonus);
          setLastComboMultiplier(payload.comboMultiplier);
          setScore(payload.totalScore);
          setCombo(payload.combo);
          setOpponentScore(payload.opponentScore);
          setOpponentAnswered(payload.opponentAnswered);
          setShowingResult(true);
          if (payload.correct) {
            setCorrectAnswers((prev) => prev + 1);
          }
          setAnsweredCount((prev) => prev + 1);
          if (payload.timedOut) {
            setSelectedIndex(-1);
          }
          break;
        }

        case "OPPONENT_PROGRESS": {
          const payload = msg.payload as OpponentProgressPayload;
          setOpponentScore(payload.totalScore);
          setOpponentAnswered(payload.answered);
          break;
        }

        case "MATCH_END": {
          const payload = msg.payload as MatchEndPayload;
          clearTimer();
          setMatchResult(payload);
          setPhase("finished");
          break;
        }

        case "ERROR": {
          const payload = msg.payload as { message: string };
          setErrorMessage(payload.message);
          setPhase("error");
          break;
        }

        case "PONG":
          break;

        default:
          console.log("Unknown WS event:", msg.type);
      }
    },
    [startTimer, clearTimer]
  );

  const connectWS = useCallback(
    (url: string) => {
      cleanup();
      setMatchResult(null);
      setErrorMessage("");
      setScore(0);
      setCombo(0);
      setCorrectAnswers(0);
      setAnsweredCount(0);
      setOpponentScore(0);
      setOpponentAnswered(0);

      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        pingRef.current = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: "PING" }));
          }
        }, 10000);
      };

      ws.onmessage = (event) => {
        try {
          const msg: WSMessage = JSON.parse(event.data);
          handleMessage(msg);
        } catch (e) {
          console.error("Failed to parse WS message:", e);
        }
      };

      ws.onerror = () => {
        setErrorMessage("การเชื่อมต่อล้มเหลว");
        setPhase("error");
      };

      ws.onclose = () => {
        clearPing();
        setPhase((prev) => {
          if (prev === "finished" || prev === "idle" || prev === "error") {
            return prev;
          }
          setErrorMessage("การเชื่อมต่อขาดหาย");
          return "error";
        });
      };
    },
    [cleanup, handleMessage, clearPing]
  );

  const createRoom = useCallback(() => {
    setPhase("creating");
    const url = getRoomWebSocketURL("create");
    connectWS(url);
  }, [connectWS]);

  const joinRoom = useCallback(
    (code: string) => {
      setPhase("joining");
      const url = getRoomWebSocketURL("join", code);
      connectWS(url);
    },
    [connectWS]
  );

  const submitAnswer = useCallback(
    (selected: number) => {
      if (answeringRef.current) return;
      answeringRef.current = true;
      setSelectedIndex(selected);

      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(
          JSON.stringify({
            type: "SUBMIT_ANSWER",
            payload: {
              index: questionIndex,
              selectedIndex: selected,
            },
          })
        );
      }
    },
    [questionIndex]
  );

  const leaveGame = useCallback(() => {
    cleanup();
    setPhase("idle");
    setRoomCode("");
    setMatchResult(null);
    setErrorMessage("");
  }, [cleanup]);

  useEffect(() => {
    return () => {
      cleanup();
    };
  }, [cleanup]);

  return {
    phase,
    roomCode,
    matchId,
    opponentName,
    totalQuestions,
    questionIndex,
    questionText,
    choices,
    timeLimit,
    timeLeft,
    difficulty,
    category,
    score,
    combo,
    correctAnswers,
    answeredCount,
    selectedIndex,
    isCorrect,
    lastScoreEarned,
    lastSpeedBonus,
    lastComboMultiplier,
    showingResult,
    opponentScore,
    opponentAnswered,
    matchResult,
    errorMessage,
    createRoom,
    joinRoom,
    submitAnswer,
    leaveGame,
  };
}

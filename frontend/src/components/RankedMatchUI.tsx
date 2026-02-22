"use client";

import { useState, useEffect } from "react";

interface RankedMatchUIProps {
  // Question
  questionIndex: number;
  totalQuestions: number;
  questionText: string;
  choices: string[];
  difficulty: string;
  category: string;
  timeLeft: number;
  timeLimit: number;
  // Player
  score: number;
  combo: number;
  correctAnswers: number;
  answeredCount: number;
  // Opponent
  opponentName: string;
  opponentScore: number;
  opponentAnswered: number;
  // Answer state
  selectedIndex: number | null;
  isCorrect: boolean | null;
  showingResult: boolean;
  lastScoreEarned: number;
  lastSpeedBonus: number;
  lastComboMultiplier: number;
  // Actions
  onAnswer: (index: number) => void;
}

export default function RankedMatchUI({
  questionIndex,
  totalQuestions,
  questionText,
  choices,
  difficulty,
  category,
  timeLeft,
  timeLimit,
  score,
  combo,
  correctAnswers,
  answeredCount,
  opponentName,
  opponentScore,
  opponentAnswered,
  selectedIndex,
  isCorrect,
  showingResult,
  lastScoreEarned,
  lastComboMultiplier,
  onAnswer,
}: RankedMatchUIProps) {
  const [preSelected, setPreSelected] = useState<number | null>(null);

  // Reset pre-selection when a new question arrives (selectedIndex resets to null)
  useEffect(() => {
    if (selectedIndex === null) {
      setPreSelected(null);
    }
  }, [selectedIndex, questionIndex]);
  const getDiffLabel = (d: string) => {
    switch (d) {
      case "easy":
        return "ง่าย";
      case "medium":
        return "ปานกลาง";
      case "hard":
        return "ยาก";
      default:
        return d;
    }
  };

  const getDiffColor = (d: string) => {
    switch (d) {
      case "easy":
        return "text-green-400 bg-green-500/10 border-green-500/20";
      case "medium":
        return "text-amber-400 bg-amber-500/10 border-amber-500/20";
      case "hard":
        return "text-red-400 bg-red-500/10 border-red-500/20";
      default:
        return "text-gray-400 bg-gray-500/10 border-gray-500/20";
    }
  };

  const getCategoryLabel = (cat: string) => {
    switch (cat) {
      case "atomic_structure":
        return "โครงสร้างอะตอม";
      case "oxidation_number":
        return "เลขออกซิเดชัน";
      case "mole_concept":
        return "โมลคอนเซ็ปต์";
      case "dilution":
        return "การเจือจาง";
      case "preparing_solution":
        return "เตรียมสารละลาย";
      case "freezing_point":
        return "จุดเยือกแข็ง";
      default:
        return cat;
    }
  };

  const timerPerc = timeLimit > 0 ? (timeLeft / timeLimit) * 100 : 0;
  const timerColor =
    timerPerc > 50
      ? "bg-blue-500"
      : timerPerc > 25
      ? "bg-amber-500"
      : "bg-red-500";
  const timerTextColor =
    timerPerc > 50
      ? "text-blue-400"
      : timerPerc > 25
      ? "text-amber-400"
      : "text-red-400";

  const getChoiceStyle = (index: number) => {
    // Already submitted — show results
    if (selectedIndex !== null) {
      if (showingResult) {
        if (isCorrect && index === selectedIndex) {
          return "bg-green-500/10 border-green-500/50 text-green-300";
        }
        if (!isCorrect && index === selectedIndex) {
          return "bg-red-500/10 border-red-500/50 text-red-300";
        }
      }
      if (index === selectedIndex) {
        return "bg-blue-500/10 border-blue-500/50 text-blue-300";
      }
      return "bg-[#12122a] border-white/5 text-gray-600";
    }
    // Pre-selected but not yet submitted
    if (preSelected === index) {
      return "bg-blue-500/10 border-blue-500/50 text-blue-300 cursor-pointer";
    }
    // Not yet selected
    return "bg-[#12122a] border-white/5 hover:border-blue-500/40 hover:bg-blue-500/5 text-white cursor-pointer";
  };

  const scoreDiff = score - opponentScore;

  return (
    <div className="w-full max-w-3xl mx-auto animate-fade-in">
      {/* Score comparison bar */}
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-3 sm:p-4 mb-3 sm:mb-4">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-1.5 sm:gap-2">
            <span className="text-blue-400 font-bold text-xs sm:text-base">คุณ</span>
            <span className="text-white font-bold text-base sm:text-lg tabular-nums">
              {score}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <span
              className={`text-[10px] sm:text-xs font-bold px-1.5 sm:px-2 py-0.5 rounded-full ${
                scoreDiff > 0
                  ? "bg-green-500/10 text-green-400"
                  : scoreDiff < 0
                  ? "bg-red-500/10 text-red-400"
                  : "bg-gray-500/10 text-gray-400"
              }`}
            >
              {scoreDiff > 0 ? "+" : ""}
              {scoreDiff}
            </span>
          </div>
          <div className="flex items-center gap-1.5 sm:gap-2">
            <span className="text-white font-bold text-base sm:text-lg tabular-nums">
              {opponentScore}
            </span>
            <span className="text-red-400 font-bold text-xs sm:text-base truncate max-w-[80px] sm:max-w-none">{opponentName}</span>
          </div>
        </div>

        {/* Progress bars */}
        <div className="flex gap-2 items-center text-[10px] sm:text-xs text-gray-500">
          <div className="flex-1">
            <div className="flex justify-between mb-1">
              <span>ข้อ {answeredCount}/{totalQuestions}</span>
              <span>ถูก {correctAnswers}</span>
            </div>
            <div className="h-1 bg-white/5 rounded-full overflow-hidden">
              <div
                className="h-full bg-blue-500 rounded-full transition-all duration-300"
                style={{
                  width: `${(answeredCount / totalQuestions) * 100}%`,
                }}
              />
            </div>
          </div>
          <div className="w-px h-6 bg-white/10" />
          <div className="flex-1">
            <div className="flex justify-between mb-1">
              <span>ข้อ {opponentAnswered}/{totalQuestions}</span>
              <span className="truncate max-w-[60px] sm:max-w-none">{opponentName}</span>
            </div>
            <div className="h-1 bg-white/5 rounded-full overflow-hidden">
              <div
                className="h-full bg-red-500 rounded-full transition-all duration-300"
                style={{
                  width: `${(opponentAnswered / totalQuestions) * 100}%`,
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Timer & Meta */}
      <div className="flex flex-wrap items-center justify-between gap-2 mb-3 sm:mb-4">
        <div className="flex items-center gap-2 sm:gap-3">
          <span
            className={`text-[10px] sm:text-xs font-semibold px-2 sm:px-3 py-1 rounded-full border ${getDiffColor(
              difficulty
            )}`}
          >
            {getDiffLabel(difficulty)}
          </span>
          <span className="text-[10px] sm:text-xs text-gray-500 hidden sm:inline">
            {getCategoryLabel(category)}
          </span>
        </div>

        <div className="flex items-center gap-2">
          <div className="w-16 sm:w-24 bg-white/5 rounded-full h-1.5 overflow-hidden">
            <div
              className={`h-full rounded-full transition-all duration-1000 ease-linear ${timerColor}`}
              style={{ width: `${timerPerc}%` }}
            />
          </div>
          <span
            className={`text-xs sm:text-sm font-mono font-bold tabular-nums ${timerTextColor} ${
              timerPerc <= 25 ? "animate-pulse" : ""
            }`}
          >
            {timeLeft}s
          </span>
        </div>

        <div className="text-xs sm:text-sm text-gray-500">
          ข้อ {questionIndex + 1}/{totalQuestions}
        </div>
      </div>

      {/* Combo indicator */}
      {combo >= 2 && (
        <div className="flex justify-center mb-3">
          <span className="text-amber-400 font-bold animate-pulse text-sm px-3 py-1 bg-amber-500/10 rounded-full">
            🔥 Combo x{combo}
            {lastComboMultiplier > 1 && (
              <span className="text-amber-300 ml-1">(×{lastComboMultiplier})</span>
            )}
          </span>
        </div>
      )}

      {/* Question */}
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-8 mb-4 sm:mb-6">
        <p className="text-base sm:text-xl md:text-2xl font-semibold text-white text-center leading-relaxed">
          {questionText}
        </p>
      </div>

      {/* Choices */}
      <div className="grid gap-2 sm:gap-3">
        {choices.map((choice, index) => (
          <button
            key={index}
            onClick={() => {
              if (selectedIndex === null && timeLeft > 0) {
                setPreSelected(index);
              }
            }}
            disabled={selectedIndex !== null || timeLeft === 0}
            className={`w-full text-left px-4 py-3 sm:px-6 sm:py-4 rounded-xl border transition-all duration-200 text-sm sm:text-base font-medium min-h-[44px] ${getChoiceStyle(
              index
            )}`}
          >
            <span className="text-gray-500 mr-2 sm:mr-3 font-mono text-xs sm:text-sm">
              {String.fromCharCode(65 + index)}.
            </span>
            {choice}
          </button>
        ))}
      </div>

      {/* Submit button — also hide when time has expired to prevent race with server timeout */}
      {selectedIndex === null && preSelected !== null && timeLeft > 0 && (
        <div className="flex justify-center mt-4 sm:mt-5">
          <button
            onClick={() => {
              onAnswer(preSelected);
            }}
            className="bg-gradient-to-r from-blue-500 to-violet-500 text-white font-bold py-3 px-10 rounded-xl text-base sm:text-lg hover:opacity-90 active:scale-[0.97] transition-all shadow-lg shadow-blue-500/20 cursor-pointer min-h-[44px]"
          >
            ส่งคำตอบ
          </button>
        </div>
      )}

      {/* Score earned popup */}
      {showingResult && lastScoreEarned > 0 && (
        <div className="flex justify-center mt-4 animate-bounce-in">
          <span className="text-green-400 font-bold text-lg">
            +{lastScoreEarned} คะแนน
          </span>
        </div>
      )}
      {showingResult && isCorrect === false && (
        <div className="flex justify-center mt-4 animate-bounce-in">
          <span className="text-red-400 font-bold text-lg">
            {selectedIndex === -1 ? "⏱️ หมดเวลา!" : "✗ ผิด!"}
          </span>
        </div>
      )}

      {/* Waiting for opponent */}
      {showingResult && (
        <div className="flex justify-center mt-4">
          <span className="text-gray-400 text-sm animate-pulse flex items-center gap-2">
            <div className="w-3 h-3 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
            รอคู่ต่อสู้...
          </span>
        </div>
      )}
    </div>
  );
}

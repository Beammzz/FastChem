"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { Suspense } from "react";
import Link from "next/link";
import { useAuth } from "@/components/AuthProvider";
import { useGame } from "@/hooks/useGame";
import QuestionCard from "@/components/QuestionCard";
import TimerBar from "@/components/TimerBar";
import ScoreBoard from "@/components/ScoreBoard";
import GameOver from "@/components/GameOver";
import PeriodicTable from "@/components/PeriodicTable";
import Calculator from "@/components/Calculator";

function GameContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { user } = useAuth();
  const totalQuestions = Number(searchParams.get("questions")) || 10;
  const difficulty = searchParams.get("difficulty") || "easy";
  const categoriesParam = searchParams.get("categories");
  const categories = categoriesParam
    ? categoriesParam.split(",").filter(Boolean)
    : undefined;

  const difficultyLabel =
    difficulty === "custom"
      ? "กำหนดเอง"
      : difficulty === "hard"
      ? "ยาก"
      : difficulty === "medium"
      ? "ปานกลาง"
      : "ง่าย";

  const {
    gameState,
    question,
    score,
    timeLeft,
    timeLimit,
    totalAnswered,
    correctAnswers,
    totalQuestions: numQuestions,
    selectedIndex,
    isCorrect,
    loading,
    combo,
    bestCombo,
    lastScoreEarned,
    lastSpeedBonus,
    lastComboMultiplier,
    difficulty: gameDifficulty,
    matchResult,
    startGame,
    answerQuestion,
  } = useGame({ totalQuestions, difficulty, categories, useMatchApi: !!user });

  // Idle state — start screen
  if (gameState === "idle") {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white flex flex-col">
        {/* Navbar */}
        <nav className="border-b border-white/5 bg-[#0a0a1a]/80 backdrop-blur-md">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
            <Link href="/" className="flex items-center gap-2 text-lg font-bold">
              <span className="text-xl">⚗️</span>
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
                FastChem
              </span>
            </Link>
            <Link href="/" className="text-sm text-gray-400 hover:text-white transition-colors">
              ← กลับ
            </Link>
          </div>
        </nav>

        <div className="flex-1 flex items-center justify-center p-4 sm:p-6">
          <div className="text-center space-y-8 w-full max-w-sm">
            <div>
              <div className="text-6xl mb-4">🧪</div>
              <h1 className="text-3xl font-bold text-white mb-2">พร้อมเล่นหรือยัง?</h1>
              <p className="text-gray-400">
                โหมด{difficultyLabel} · {totalQuestions} ข้อ
              </p>
            </div>

            <div className="bg-[#12122a] border border-white/5 rounded-xl p-6 text-left max-w-sm mx-auto">
              <div className="space-y-3 text-sm text-gray-400">
                <div className="flex items-center gap-3">
                  <span className="w-6 h-6 rounded bg-violet-500/10 flex items-center justify-center text-xs">⚡</span>
                  <span>ตอบเร็ว = คะแนนสูง (Speed Bonus)</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="w-6 h-6 rounded bg-violet-500/10 flex items-center justify-center text-xs">⏱</span>
                  <span>จับเวลาแต่ละข้อ หมดเวลา = 0 คะแนน</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="w-6 h-6 rounded bg-violet-500/10 flex items-center justify-center text-xs">🔥</span>
                  <span>ตอบถูกต่อเนื่องเพื่อ Combo!</span>
                </div>
              </div>
            </div>

            <button
              onClick={startGame}
              className="w-full sm:w-auto bg-violet-500 hover:bg-violet-600 text-white font-bold text-lg px-12 py-4 rounded-xl transition-all hover:scale-[1.02] active:scale-[0.98] shadow-lg shadow-violet-500/20 min-h-[44px]"
            >
              เริ่มเกม
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Finished state
  if (gameState === "finished") {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white flex flex-col">
        <nav className="border-b border-white/5 bg-[#0a0a1a]/80 backdrop-blur-md">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
            <Link href="/" className="flex items-center gap-2 text-lg font-bold">
              <span className="text-xl">⚗️</span>
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
                FastChem
              </span>
            </Link>
          </div>
        </nav>
        <div className="flex-1 flex items-center justify-center p-6">
          <GameOver
            score={score}
            totalAnswered={totalAnswered}
            correctAnswers={correctAnswers}
            totalQuestions={numQuestions}
            difficulty={gameDifficulty}
            bestCombo={bestCombo}
            matchResult={matchResult}
            onPlayAgain={startGame}
            onHome={() => router.push("/")}
          />
        </div>
      </div>
    );
  }

  // Playing state
  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white flex flex-col">
      {/* Compact game header */}
      <div className="border-b border-white/5 bg-[#0a0a1a]/80 backdrop-blur-md">
        <div className="max-w-3xl mx-auto px-3 sm:px-6 py-2 sm:py-0 sm:h-14 flex flex-col sm:flex-row items-center justify-between gap-1 sm:gap-0">
          <div className="hidden sm:flex items-center gap-2 text-lg font-bold">
            <span className="text-xl">⚗️</span>
            <span className="text-violet-400 text-sm font-medium">FastChem</span>
          </div>
          <div className="flex items-center gap-3 sm:gap-4 w-full sm:w-auto justify-between sm:justify-end">
            <TimerBar timeLeft={timeLeft} totalTime={timeLimit} />
            <ScoreBoard score={score} totalAnswered={totalAnswered} correctAnswers={correctAnswers} totalQuestions={numQuestions} combo={combo} comboMultiplier={lastComboMultiplier} />
          </div>
        </div>
      </div>

      {/* Question area */}
      <div className="flex-1 flex items-center justify-center p-3 sm:p-6 pb-20 sm:pb-6">
        <div className="w-full max-w-2xl">
          {loading ? (
            <div className="flex justify-center py-16">
              <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
            </div>
          ) : question ? (
            <QuestionCard
              question={question}
              selectedIndex={selectedIndex}
              isCorrect={isCorrect}
              scoreEarned={lastScoreEarned}
              speedBonus={lastSpeedBonus}
              comboMultiplier={lastComboMultiplier}
              onAnswer={answerQuestion}
            />
          ) : (
            <div className="text-center text-gray-500 py-16 bg-[#12122a] border border-white/5 rounded-2xl">
              <div className="text-4xl mb-4">⚠️</div>
              <p>โหลดคำถามไม่สำเร็จ</p>
              <p className="text-sm mt-1">ตรวจสอบว่า backend ทำงานอยู่บนพอร์ต 8080</p>
            </div>
          )}
        </div>
      </div>

      {/* Mini Periodic Table & Calculator */}
      <PeriodicTable />
      <Calculator />
    </div>
  );
}

export default function GamePage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-[#0a0a1a] flex items-center justify-center">
          <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
        </div>
      }
    >
      <GameContent />
    </Suspense>
  );
}

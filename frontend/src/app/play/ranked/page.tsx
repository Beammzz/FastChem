"use client";

import { useRouter } from "next/navigation";
import Navbar from "@/components/Navbar";
import RankedQueue from "@/components/RankedQueue";
import RankedMatchUI from "@/components/RankedMatchUI";
import RankedGameOver from "@/components/RankedGameOver";
import PeriodicTable from "@/components/PeriodicTable";
import Calculator from "@/components/Calculator";
import { useRankedGame } from "@/hooks/useRankedGame";
import { useAuth } from "@/components/AuthProvider";
import { useEffect, useState } from "react";
import { fetchRankedStats } from "@/lib/api";
import { RankedStats } from "@/types";

export default function RankedPage() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const ranked = useRankedGame();
  const [stats, setStats] = useState<RankedStats | null>(null);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [authLoading, user, router]);

  // Fetch ranked stats on mount
  useEffect(() => {
    if (user) {
      fetchRankedStats()
        .then(setStats)
        .catch(() => setStats(null));
    }
  }, [user]);

  if (authLoading) {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
        <Navbar />
        <div className="flex justify-center py-32">
          <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
        </div>
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 pt-6 sm:pt-8 pb-16 pb-safe">
        {/* Idle - Show ranked lobby */}
        {ranked.phase === "idle" && (
          <div className="animate-fade-in">
            <div className="text-center mb-10">
              <h1 className="text-4xl font-extrabold mb-2">
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-violet-400">
                  แข่งขัน 1vs1
                </span>
              </h1>
              <p className="text-gray-400">
                ท้าทายผู้เล่นคนอื่นแบบเรียลไทม์ ไต่อันดับ Rating
              </p>
            </div>

            {/* Ranked stats card */}
            <div className="max-w-md mx-auto mb-8">
              <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
                <h3 className="text-sm text-gray-500 font-medium mb-4">
                  สถิติ Ranked ของคุณ
                </h3>
                {stats ? (
                  <div className="grid grid-cols-2 gap-4">
                    <div className="text-center">
                      <div className="text-3xl font-extrabold text-blue-400">
                        {stats.rating}
                      </div>
                      <div className="text-xs text-gray-500 mt-1">Rating</div>
                    </div>
                    <div className="text-center">
                      <div className="text-3xl font-extrabold text-violet-400">
                        {stats.highestRating}
                      </div>
                      <div className="text-xs text-gray-500 mt-1">
                        สูงสุด
                      </div>
                    </div>
                    <div className="text-center">
                      <div className="text-xl font-bold text-green-400">
                        {stats.rankedWins}
                      </div>
                      <div className="text-xs text-gray-500 mt-1">ชนะ</div>
                    </div>
                    <div className="text-center">
                      <div className="text-xl font-bold text-red-400">
                        {stats.rankedLosses}
                      </div>
                      <div className="text-xs text-gray-500 mt-1">แพ้</div>
                    </div>
                  </div>
                ) : (
                  <div className="text-center text-gray-500 py-4">
                    <div className="text-3xl font-extrabold text-blue-400 mb-1">
                      1200
                    </div>
                    <div className="text-xs">Rating เริ่มต้น</div>
                  </div>
                )}
              </div>
            </div>

            {/* How it works */}
            <div className="max-w-md mx-auto mb-8">
              <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
                <h3 className="text-sm text-gray-500 font-medium mb-4">
                  วิธีการเล่น
                </h3>
                <div className="space-y-3 text-sm text-gray-400">
                  <div className="flex items-start gap-3">
                    <span className="text-blue-400 font-bold mt-0.5">1.</span>
                    <span>กดปุ่มค้นหาคู่ต่อสู้ ระบบจะจับคู่ผู้เล่นที่มี Rating ใกล้เคียง</span>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-blue-400 font-bold mt-0.5">2.</span>
                    <span>ตอบคำถามเคมี 10 ข้อ (4 ง่าย, 3 ปานกลาง, 3 ยาก) ทั้งสองฝ่ายได้คำถามเดียวกัน</span>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-blue-400 font-bold mt-0.5">3.</span>
                    <span>คะแนนคำนวณจากความถูกต้อง ความเร็ว และ Combo ผู้เล่นที่ได้คะแนนมากกว่าชนะ</span>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-blue-400 font-bold mt-0.5">4.</span>
                    <span>Rating จะเปลี่ยนตามระบบ ELO — ชนะคนที่แข็งกว่าได้ Rating มากกว่า</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Find match button */}
            <div className="max-w-md mx-auto text-center">
              <button
                onClick={ranked.joinQueue}
                className="w-full bg-gradient-to-r from-blue-500 to-violet-500 text-white font-bold py-4 px-8 rounded-xl text-lg hover:opacity-90 transition-opacity cursor-pointer min-h-[48px]"
              >
                ⚔️ ค้นหาคู่ต่อสู้
              </button>
              <button
                onClick={() => router.push("/play")}
                className="mt-4 text-gray-500 hover:text-gray-300 transition-colors text-sm cursor-pointer"
              >
                ← กลับไปเลือกโหมด
              </button>
            </div>
          </div>
        )}

        {/* Queuing */}
        {ranked.phase === "queuing" && (
          <RankedQueue
            rating={ranked.queueRating}
            onCancel={ranked.leaveGame}
          />
        )}

        {/* Match found transition */}
        {ranked.phase === "matched" && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="text-center">
              <div className="text-5xl mb-4">⚔️</div>
              <h2 className="text-3xl font-extrabold text-white mb-2">
                พบคู่ต่อสู้!
              </h2>
              <p className="text-xl text-blue-400 font-bold">
                vs {ranked.opponentName}
              </p>
              <p className="text-gray-500 text-sm mt-4 animate-pulse">
                กำลังเริ่มแมตช์...
              </p>
            </div>
          </div>
        )}

        {/* Playing */}
        {ranked.phase === "playing" && (
          <>
          <RankedMatchUI
            questionIndex={ranked.questionIndex}
            totalQuestions={ranked.totalQuestions}
            questionText={ranked.questionText}
            choices={ranked.choices}
            difficulty={ranked.difficulty}
            category={ranked.category}
            timeLeft={ranked.timeLeft}
            timeLimit={ranked.timeLimit}
            score={ranked.score}
            combo={ranked.combo}
            correctAnswers={ranked.correctAnswers}
            answeredCount={ranked.answeredCount}
            opponentName={ranked.opponentName}
            opponentScore={ranked.opponentScore}
            opponentAnswered={ranked.opponentAnswered}
            selectedIndex={ranked.selectedIndex}
            isCorrect={ranked.isCorrect}
            showingResult={ranked.showingResult}
            lastScoreEarned={ranked.lastScoreEarned}
            lastSpeedBonus={ranked.lastSpeedBonus}
            lastComboMultiplier={ranked.lastComboMultiplier}
            onAnswer={ranked.submitAnswer}
          />
          <PeriodicTable />
          <Calculator />
          </>
        )}

        {/* Finished */}
        {ranked.phase === "finished" && ranked.matchResult && (
          <RankedGameOver
            result={ranked.matchResult}
            onPlayAgain={ranked.joinQueue}
            onHome={() => router.push("/play")}
          />
        )}

        {/* Error */}
        {ranked.phase === "error" && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="bg-[#12122a] border border-red-500/20 rounded-2xl p-8 max-w-md w-full text-center">
              <div className="text-5xl mb-4">⚠️</div>
              <h2 className="text-xl font-bold text-white mb-2">
                เกิดข้อผิดพลาด
              </h2>
              <p className="text-gray-400 text-sm mb-6">
                {ranked.errorMessage || "การเชื่อมต่อล้มเหลว"}
              </p>
              <div className="flex gap-3">
                <button
                  onClick={ranked.joinQueue}
                  className="flex-1 bg-blue-500 text-white font-bold py-3 px-6 rounded-xl hover:bg-blue-600 transition-colors cursor-pointer"
                >
                  ลองอีกครั้ง
                </button>
                <button
                  onClick={() => router.push("/play")}
                  className="flex-1 bg-white/5 border border-white/10 text-gray-300 font-bold py-3 px-6 rounded-xl hover:bg-white/10 transition-colors cursor-pointer"
                >
                  กลับ
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

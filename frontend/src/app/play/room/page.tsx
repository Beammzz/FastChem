"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useEffect } from "react";
import Navbar from "@/components/Navbar";
import RankedMatchUI from "@/components/RankedMatchUI";
import RoomGameOver from "@/components/RoomGameOver";
import PeriodicTable from "@/components/PeriodicTable";
import Calculator from "@/components/Calculator";
import { useRoomGame } from "@/hooks/useRoomGame";
import { useAuth } from "@/components/AuthProvider";

export default function RoomPage() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const room = useRoomGame();
  const [joinCode, setJoinCode] = useState("");

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [authLoading, user, router]);

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
        {/* Idle - Show room lobby */}
        {room.phase === "idle" && (
          <div className="animate-fade-in">
            <div className="text-center mb-10">
              <h1 className="text-4xl font-extrabold mb-2">
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 to-teal-400">
                  สร้างห้อง
                </span>
              </h1>
              <p className="text-gray-400">
                เล่นกับเพื่อน — ไม่มีผลต่อ Rating
              </p>
            </div>

            <div className="max-w-md mx-auto space-y-6">
              {/* Create room */}
              <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
                <h3 className="text-sm text-gray-500 font-medium mb-4">
                  สร้างห้องใหม่
                </h3>
                <p className="text-gray-400 text-sm mb-4">
                  สร้างห้องเเล้วแชร์รหัส 6 หลักให้เพื่อนเข้าร่วม
                </p>
                <button
                  onClick={room.createRoom}
                  className="w-full bg-gradient-to-r from-emerald-500 to-teal-500 text-white font-bold py-3 px-6 rounded-xl text-lg hover:opacity-90 transition-opacity cursor-pointer min-h-[48px]"
                >
                  🏠 สร้างห้อง
                </button>
              </div>

              {/* Join room */}
              <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
                <h3 className="text-sm text-gray-500 font-medium mb-4">
                  เข้าร่วมห้อง
                </h3>
                <p className="text-gray-400 text-sm mb-4">
                  ใส่รหัส 6 หลักจากเพื่อนเพื่อเข้าร่วมเกม
                </p>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input
                    type="text"
                    inputMode="numeric"
                    maxLength={6}
                    value={joinCode}
                    onChange={(e) =>
                      setJoinCode(e.target.value.replace(/\D/g, "").slice(0, 6))
                    }
                    placeholder="000000"
                    className="flex-1 bg-[#0a0a1a] border border-white/10 rounded-xl px-4 py-3 text-center text-xl sm:text-2xl font-mono tracking-[0.3em] sm:tracking-[0.5em] text-white placeholder-gray-600 focus:outline-none focus:border-emerald-500/50 min-h-[44px]"
                  />
                  <button
                    onClick={() => {
                      if (joinCode.length === 6) {
                        room.joinRoom(joinCode);
                      }
                    }}
                    disabled={joinCode.length !== 6}
                    className="bg-emerald-500 text-white font-bold py-3 px-6 rounded-xl hover:bg-emerald-600 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer min-h-[44px]"
                  >
                    เข้าร่วม
                  </button>
                </div>
              </div>

              {/* How it works */}
              <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
                <h3 className="text-sm text-gray-500 font-medium mb-4">
                  วิธีการเล่น
                </h3>
                <div className="space-y-3 text-sm text-gray-400">
                  <div className="flex items-start gap-3">
                    <span className="text-emerald-400 font-bold mt-0.5">1.</span>
                    <span>สร้างห้องเพื่อรับรหัส 6 หลัก หรือใส่รหัสเพื่อเข้าร่วมห้อง</span>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-emerald-400 font-bold mt-0.5">2.</span>
                    <span>เมื่อผู้เล่นทั้งสองพร้อม เกมจะเริ่มอัตโนมัติ</span>
                  </div>
                  <div className="flex items-start gap-3">
                    <span className="text-emerald-400 font-bold mt-0.5">3.</span>
                    <span>ตอบคำถามเคมี 10 ข้อเหมือน Ranked แต่ไม่มีผลต่อ Rating</span>
                  </div>
                </div>
              </div>

              <div className="text-center">
                <button
                  onClick={() => router.push("/play")}
                  className="text-gray-500 hover:text-gray-300 transition-colors text-sm cursor-pointer"
                >
                  ← กลับไปเลือกโหมด
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Creating / Waiting for guest */}
        {(room.phase === "creating" || room.phase === "waiting") && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-10 max-w-md w-full text-center">
              {/* Animated icon */}
              <div className="relative mx-auto mb-6 w-20 h-20">
                <div className="absolute inset-0 rounded-full border-4 border-emerald-500/20 animate-ping" />
                <div className="absolute inset-0 rounded-full border-4 border-emerald-500/40 animate-pulse" />
                <div className="relative w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center text-3xl">
                  🏠
                </div>
              </div>

              <h2 className="text-2xl font-bold text-white mb-2">
                รอผู้เล่นเข้าร่วม...
              </h2>

              {room.roomCode && (
                <div className="my-6">
                  <div className="text-xs text-gray-500 mb-2">รหัสห้อง</div>
                  <div className="text-5xl font-mono font-extrabold tracking-[0.3em] text-emerald-400">
                    {room.roomCode}
                  </div>
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(room.roomCode);
                    }}
                    className="mt-3 text-xs text-gray-500 hover:text-emerald-400 transition-colors cursor-pointer"
                  >
                    📋 คัดลอกรหัส
                  </button>
                </div>
              )}

              <p className="text-gray-400 text-sm mb-6">
                แชร์รหัสนี้ให้เพื่อนเพื่อเข้าร่วมเกม
              </p>

              {/* Loading bar */}
              <div className="w-full h-1 bg-white/5 rounded-full overflow-hidden mb-6">
                <div className="h-full bg-gradient-to-r from-emerald-500 to-teal-500 rounded-full animate-loading-bar" />
              </div>

              <button
                onClick={room.leaveGame}
                className="text-gray-500 hover:text-red-400 transition-colors text-sm font-medium cursor-pointer"
              >
                ยกเลิก
              </button>
            </div>

            <style jsx>{`
              @keyframes loading-bar {
                0% { width: 0%; margin-left: 0%; }
                50% { width: 60%; margin-left: 20%; }
                100% { width: 0%; margin-left: 100%; }
              }
              .animate-loading-bar {
                animation: loading-bar 2s ease-in-out infinite;
              }
            `}</style>
          </div>
        )}

        {/* Joining */}
        {room.phase === "joining" && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-10 max-w-md w-full text-center">
              <div className="relative mx-auto mb-6 w-20 h-20">
                <div className="absolute inset-0 rounded-full border-4 border-emerald-500/20 animate-pulse" />
                <div className="relative w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center text-3xl">
                  🔗
                </div>
              </div>
              <h2 className="text-2xl font-bold text-white mb-2">
                กำลังเข้าร่วมห้อง...
              </h2>
              <p className="text-gray-400 text-sm animate-pulse">
                รอเริ่มเกม...
              </p>
            </div>
          </div>
        )}

        {/* Match found transition */}
        {room.phase === "matched" && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="text-center">
              <div className="text-5xl mb-4">🎮</div>
              <h2 className="text-3xl font-extrabold text-white mb-2">
                พร้อมแล้ว!
              </h2>
              <p className="text-xl text-emerald-400 font-bold">
                vs {room.opponentName}
              </p>
              <p className="text-gray-500 text-sm mt-4 animate-pulse">
                กำลังเริ่มเกม...
              </p>
            </div>
          </div>
        )}

        {/* Playing */}
        {room.phase === "playing" && (
          <>
            <RankedMatchUI
              questionIndex={room.questionIndex}
              totalQuestions={room.totalQuestions}
              questionText={room.questionText}
              choices={room.choices}
              difficulty={room.difficulty}
              category={room.category}
              timeLeft={room.timeLeft}
              timeLimit={room.timeLimit}
              score={room.score}
              combo={room.combo}
              correctAnswers={room.correctAnswers}
              answeredCount={room.answeredCount}
              opponentName={room.opponentName}
              opponentScore={room.opponentScore}
              opponentAnswered={room.opponentAnswered}
              selectedIndex={room.selectedIndex}
              isCorrect={room.isCorrect}
              showingResult={room.showingResult}
              lastScoreEarned={room.lastScoreEarned}
              lastSpeedBonus={room.lastSpeedBonus}
              lastComboMultiplier={room.lastComboMultiplier}
              onAnswer={room.submitAnswer}
            />
            <PeriodicTable />
            <Calculator />
          </>
        )}

        {/* Finished */}
        {room.phase === "finished" && room.matchResult && (
          <RoomGameOver
            result={room.matchResult}
            onPlayAgain={room.createRoom}
            onHome={() => router.push("/play")}
          />
        )}

        {/* Error */}
        {room.phase === "error" && (
          <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
            <div className="bg-[#12122a] border border-red-500/20 rounded-2xl p-8 max-w-md w-full text-center">
              <div className="text-5xl mb-4">⚠️</div>
              <h2 className="text-xl font-bold text-white mb-2">
                เกิดข้อผิดพลาด
              </h2>
              <p className="text-gray-400 text-sm mb-6">
                {room.errorMessage || "การเชื่อมต่อล้มเหลว"}
              </p>
              <div className="flex gap-3">
                <button
                  onClick={room.leaveGame}
                  className="flex-1 bg-emerald-500 text-white font-bold py-3 px-6 rounded-xl hover:bg-emerald-600 transition-colors cursor-pointer"
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

"use client";

import { MatchEndPayload } from "@/types";
import { useAuth } from "@/components/AuthProvider";

interface RoomGameOverProps {
  result: MatchEndPayload;
  onPlayAgain: () => void;
  onHome: () => void;
}

export default function RoomGameOver({
  result,
  onPlayAgain,
  onHome,
}: RoomGameOverProps) {
  const { user } = useAuth();
  const isWinner = user && result.winnerId === user.id;

  const me =
    user && result.player1.userId === user.id ? result.player1 : result.player2;
  const opponent =
    user && result.player1.userId === user.id ? result.player2 : result.player1;

  return (
    <div className="w-full max-w-lg mx-auto animate-fade-in">
      {/* Result header */}
      <div className="text-center mb-8">
        <div className="text-6xl mb-4">{isWinner ? "🏆" : "😔"}</div>
        <h1 className="text-3xl font-extrabold mb-2">
          <span
            className={`bg-clip-text text-transparent bg-gradient-to-r ${
              isWinner
                ? "from-yellow-400 to-amber-400"
                : "from-gray-400 to-gray-500"
            }`}
          >
            {isWinner ? "ชนะ!" : "แพ้"}
          </span>
        </h1>

        {/* No rating change for custom rooms */}
        <div className="flex items-center justify-center gap-2 mt-2">
          <span className="text-gray-500 text-sm">Friendly Match — ไม่มีผลต่อ Rating</span>
        </div>
      </div>

      {/* Score comparison */}
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6 mb-4">
        <div className="flex items-center justify-between">
          {/* Me */}
          <div className="text-center flex-1">
            <div className="w-14 h-14 rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center text-xl font-bold text-white mx-auto mb-2">
              {me.username.charAt(0).toUpperCase()}
            </div>
            <div className="text-sm text-emerald-400 font-medium mb-1">
              {me.username}
            </div>
            <div className="text-3xl font-extrabold text-white">
              {me.totalScore}
            </div>
          </div>

          {/* VS */}
          <div className="text-gray-600 font-bold text-lg px-4">VS</div>

          {/* Opponent */}
          <div className="text-center flex-1">
            <div className="w-14 h-14 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center text-xl font-bold text-white mx-auto mb-2">
              {opponent.username.charAt(0).toUpperCase()}
            </div>
            <div className="text-sm text-orange-400 font-medium mb-1">
              {opponent.username}
            </div>
            <div className="text-3xl font-extrabold text-white">
              {opponent.totalScore}
            </div>
          </div>
        </div>
      </div>

      {/* Detailed stats */}
      <div className="grid grid-cols-2 gap-3 mb-6">
        {/* My stats */}
        <div className="bg-[#12122a] border border-white/5 rounded-xl p-4">
          <h3 className="text-xs text-gray-500 font-medium mb-3 uppercase tracking-wider">
            สถิติของคุณ
          </h3>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-400">ตอบถูก</span>
              <span className="text-white font-bold">
                {me.correctAnswers}/10
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Combo สูงสุด</span>
              <span className="text-amber-400 font-bold">{me.bestCombo}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">คะแนนข้อยาก</span>
              <span className="text-red-400 font-bold">{me.hardScore}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">เวลารวม</span>
              <span className="text-gray-300 font-bold">
                {me.totalTime.toFixed(1)}s
              </span>
            </div>
          </div>
        </div>

        {/* Opponent stats */}
        <div className="bg-[#12122a] border border-white/5 rounded-xl p-4">
          <h3 className="text-xs text-gray-500 font-medium mb-3 uppercase tracking-wider">
            สถิติคู่ต่อสู้
          </h3>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-400">ตอบถูก</span>
              <span className="text-white font-bold">
                {opponent.correctAnswers}/10
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Combo สูงสุด</span>
              <span className="text-amber-400 font-bold">
                {opponent.bestCombo}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">คะแนนข้อยาก</span>
              <span className="text-red-400 font-bold">
                {opponent.hardScore}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">เวลารวม</span>
              <span className="text-gray-300 font-bold">
                {opponent.totalTime.toFixed(1)}s
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex gap-3">
        <button
          onClick={onPlayAgain}
          className="flex-1 bg-gradient-to-r from-emerald-500 to-teal-500 text-white font-bold py-3 px-6 rounded-xl hover:opacity-90 transition-opacity cursor-pointer"
        >
          สร้างห้องใหม่
        </button>
        <button
          onClick={onHome}
          className="flex-1 bg-white/5 border border-white/10 text-gray-300 font-bold py-3 px-6 rounded-xl hover:bg-white/10 transition-colors cursor-pointer"
        >
          กลับหน้าหลัก
        </button>
      </div>
    </div>
  );
}

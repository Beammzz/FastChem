"use client";

import { MatchEndPayload } from "@/types";
import { useAuth } from "@/components/AuthProvider";

interface RankedGameOverProps {
  result: MatchEndPayload;
  onPlayAgain: () => void;
  onHome: () => void;
}

export default function RankedGameOver({
  result,
  onPlayAgain,
  onHome,
}: RankedGameOverProps) {
  const { user } = useAuth();
  const isWinner = user && result.winnerId === user.id;

  // Determine which player is "me" and which is "opponent"
  const me =
    user && result.player1.userId === user.id ? result.player1 : result.player2;
  const opponent =
    user && result.player1.userId === user.id ? result.player2 : result.player1;

  const ratingPrefix = result.ratingChange >= 0 ? "+" : "";
  const ratingColor = result.ratingChange >= 0 ? "text-green-400" : "text-red-400";

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

        {/* Rating change */}
        <div className="flex items-center justify-center gap-2 mt-2">
          <span className="text-gray-400 text-sm">Rating:</span>
          <span className="text-white font-bold">{result.newRating}</span>
          <span className={`text-sm font-bold ${ratingColor}`}>
            ({ratingPrefix}{result.ratingChange})
          </span>
        </div>
      </div>

      {/* Score comparison */}
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6 mb-4">
        <div className="flex items-center justify-between">
          {/* Me */}
          <div className="text-center flex-1">
            <div className="w-14 h-14 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-xl font-bold text-white mx-auto mb-2">
              {me.username.charAt(0).toUpperCase()}
            </div>
            <div className="text-sm text-blue-400 font-medium mb-1">
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
            <div className="w-14 h-14 rounded-full bg-gradient-to-br from-red-400 to-red-600 flex items-center justify-center text-xl font-bold text-white mx-auto mb-2">
              {opponent.username.charAt(0).toUpperCase()}
            </div>
            <div className="text-sm text-red-400 font-medium mb-1">
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
      <div className="flex flex-col sm:flex-row gap-3">
        <button
          onClick={onPlayAgain}
          className="flex-1 bg-gradient-to-r from-blue-500 to-violet-500 text-white font-bold py-3 px-6 rounded-xl hover:opacity-90 transition-opacity cursor-pointer min-h-[44px]"
        >
          เล่นอีกครั้ง
        </button>
        <button
          onClick={onHome}
          className="flex-1 bg-white/5 border border-white/10 text-gray-300 font-bold py-3 px-6 rounded-xl hover:bg-white/10 transition-colors cursor-pointer min-h-[44px]"
        >
          กลับหน้าหลัก
        </button>
      </div>
    </div>
  );
}

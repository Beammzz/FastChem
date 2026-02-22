"use client";

interface ScoreBoardProps {
  score: number;
  totalAnswered: number;
  correctAnswers: number;
  totalQuestions: number;
  combo: number;
  comboMultiplier?: number;
}

export default function ScoreBoard({
  score,
  totalAnswered,
  correctAnswers,
  totalQuestions,
  combo,
  comboMultiplier = 1,
}: ScoreBoardProps) {
  return (
    <div className="flex items-center gap-2 sm:gap-4 text-xs sm:text-sm">
      <div className="flex items-center gap-1 sm:gap-1.5">
        <span className="text-gray-500 hidden sm:inline">คะแนน</span>
        <span className="font-bold text-white tabular-nums">{score}</span>
      </div>
      <div className="w-px h-4 bg-white/10" />
      <div className="flex items-center gap-1 sm:gap-1.5">
        <span className="text-gray-500 hidden sm:inline">ถูก</span>
        <span className="font-bold text-violet-400 tabular-nums">
          {correctAnswers}/{totalAnswered}
        </span>
      </div>
      <div className="w-px h-4 bg-white/10" />
      <div className="flex items-center gap-1 sm:gap-1.5">
        <span className="text-gray-500 hidden sm:inline">ข้อ</span>
        <span className="font-bold text-gray-300 tabular-nums">
          {totalAnswered}/{totalQuestions}
        </span>
      </div>
      {combo >= 2 && (
        <>
          <div className="w-px h-4 bg-white/10" />
          <div className="flex items-center gap-1.5">
            <span className="text-amber-400 font-bold tabular-nums animate-pulse">
              🔥 x{combo}
            </span>
            {comboMultiplier > 1 && (
              <span className="text-amber-300 text-xs font-semibold">
                (×{comboMultiplier})
              </span>
            )}
          </div>
        </>
      )}
    </div>
  );
}

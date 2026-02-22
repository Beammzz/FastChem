"use client";

interface TimerBarProps {
  timeLeft: number;
  totalTime: number;
}

export default function TimerBar({ timeLeft, totalTime }: TimerBarProps) {
  const percentage = totalTime > 0 ? (timeLeft / totalTime) * 100 : 0;

  const getColor = () => {
    if (percentage > 50) return "text-violet-400";
    if (percentage > 25) return "text-amber-400";
    return "text-red-400";
  };

  const getBarColor = () => {
    if (percentage > 50) return "bg-violet-500";
    if (percentage > 25) return "bg-amber-500";
    return "bg-red-500";
  };

  return (
    <div className="flex items-center gap-2">
      <div className="w-24 bg-white/5 rounded-full h-1.5 overflow-hidden">
        <div
          className={`h-full rounded-full transition-all duration-1000 ease-linear ${getBarColor()}`}
          style={{ width: `${percentage}%` }}
        />
      </div>
      <span
        className={`text-sm font-mono font-bold tabular-nums ${getColor()} ${
          percentage <= 25 ? "animate-pulse" : ""
        }`}
      >
        {timeLeft}s
      </span>
    </div>
  );
}

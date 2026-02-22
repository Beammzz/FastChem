"use client";

import { useState, useEffect } from "react";

interface RankedQueueProps {
  rating: number | null;
  onCancel: () => void;
}

export default function RankedQueue({ rating, onCancel }: RankedQueueProps) {
  const [dots, setDots] = useState("");
  const [elapsed, setElapsed] = useState(0);

  useEffect(() => {
    const dotInterval = setInterval(() => {
      setDots((prev) => (prev.length >= 3 ? "" : prev + "."));
    }, 500);
    const timeInterval = setInterval(() => {
      setElapsed((prev) => prev + 1);
    }, 1000);
    return () => {
      clearInterval(dotInterval);
      clearInterval(timeInterval);
    };
  }, []);

  const formatTime = (s: number) => {
    const m = Math.floor(s / 60);
    const sec = s % 60;
    return `${m}:${sec.toString().padStart(2, "0")}`;
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] animate-fade-in">
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-10 max-w-md w-full text-center">
        {/* Animated icon */}
        <div className="relative mx-auto mb-6 w-20 h-20">
          <div className="absolute inset-0 rounded-full border-4 border-blue-500/20 animate-ping" />
          <div className="absolute inset-0 rounded-full border-4 border-blue-500/40 animate-pulse" />
          <div className="relative w-20 h-20 rounded-full bg-blue-500/10 flex items-center justify-center text-3xl">
            ⚔️
          </div>
        </div>

        <h2 className="text-2xl font-bold text-white mb-2">
          กำลังค้นหาคู่ต่อสู้{dots}
        </h2>
        <p className="text-gray-400 text-sm mb-6">
          ระบบกำลังจับคู่ผู้เล่นที่มี Rating ใกล้เคียง
        </p>

        {/* Stats */}
        <div className="flex justify-center gap-6 mb-8">
          {rating !== null && (
            <div className="text-center">
              <div className="text-xs text-gray-500 mb-1">Rating ของคุณ</div>
              <div className="text-xl font-bold text-blue-400">{rating}</div>
            </div>
          )}
          <div className="text-center">
            <div className="text-xs text-gray-500 mb-1">เวลารอ</div>
            <div className="text-xl font-bold text-gray-300 tabular-nums">
              {formatTime(elapsed)}
            </div>
          </div>
        </div>

        {/* Loading bar */}
        <div className="w-full h-1 bg-white/5 rounded-full overflow-hidden mb-6">
          <div className="h-full bg-gradient-to-r from-blue-500 to-violet-500 rounded-full animate-loading-bar" />
        </div>

        <button
          onClick={onCancel}
          className="text-gray-500 hover:text-red-400 transition-colors text-sm font-medium cursor-pointer"
        >
          ยกเลิกการค้นหา
        </button>
      </div>

      <style jsx>{`
        @keyframes loading-bar {
          0% {
            width: 0%;
            margin-left: 0%;
          }
          50% {
            width: 60%;
            margin-left: 20%;
          }
          100% {
            width: 0%;
            margin-left: 100%;
          }
        }
        .animate-loading-bar {
          animation: loading-bar 2s ease-in-out infinite;
        }
      `}</style>
    </div>
  );
}

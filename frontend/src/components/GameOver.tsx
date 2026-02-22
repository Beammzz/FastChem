"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { useAuth } from "@/components/AuthProvider";
import { submitScore } from "@/lib/api";
import { EndMatchResponse } from "@/types";

interface GameOverProps {
  score: number;
  totalAnswered: number;
  correctAnswers: number;
  totalQuestions: number;
  difficulty: string;
  bestCombo: number;
  matchResult?: EndMatchResponse | null;
  onPlayAgain: () => void;
  onHome: () => void;
}

export default function GameOver({
  score,
  totalAnswered,
  correctAnswers,
  totalQuestions,
  difficulty,
  bestCombo,
  matchResult,
  onPlayAgain,
  onHome,
}: GameOverProps) {
  const { user, refreshUser } = useAuth();
  const [scoreSaved, setScoreSaved] = useState(false);
  const [saving, setSaving] = useState(false);
  const submittedRef = useRef(false);
  const refreshedRef = useRef(false);

  // Use server-validated match result when available
  const finalScore = matchResult?.totalScore ?? score;
  const finalCorrect = matchResult?.correctAnswers ?? correctAnswers;
  const finalAnswered = matchResult?.totalAnswered ?? totalAnswered;
  const finalBestCombo = matchResult?.bestCombo ?? bestCombo;

  const accuracy =
    finalAnswered > 0 ? Math.round((finalCorrect / finalAnswered) * 100) : 0;

  const difficultyLabel =
    difficulty === "hard"
      ? "ยาก"
      : difficulty === "medium"
      ? "ปานกลาง"
      : "ง่าย";

  // Auto-submit score when user is logged in (skip if match API was used — already saved)
  useEffect(() => {
    if (matchResult) {
      // Match API already persisted the score on the backend
      setScoreSaved(true);
      if (!refreshedRef.current) {
        refreshedRef.current = true;
        refreshUser();
      }
      return;
    }
    if (user && finalScore > 0 && !submittedRef.current) {
      submittedRef.current = true;
      setSaving(true);
      submitScore({
        score: finalScore,
        totalAnswered: finalAnswered,
        correctAnswers: finalCorrect,
        difficulty,
        timeLimit: totalQuestions,
        timeSpent: 0,
      })
        .then(() => {
          setScoreSaved(true);
          if (!refreshedRef.current) {
            refreshedRef.current = true;
            refreshUser();
          }
        })
        .catch(() => {
          submittedRef.current = false;
        })
        .finally(() => setSaving(false));
    }
  }, [user, finalScore, finalAnswered, finalCorrect, difficulty, totalQuestions, refreshUser, matchResult]);

  const getMessage = () => {
    if (accuracy >= 90) return { emoji: "🧪", text: "อัจฉริยะเคมี!" };
    if (accuracy >= 70) return { emoji: "🔬", text: "นักวิทยาศาสตร์ที่ยอดเยี่ยม!" };
    if (accuracy >= 50) return { emoji: "⚗️", text: "ทำได้ดี!" };
    return { emoji: "📚", text: "ฝึกฝนต่อไปนะ!" };
  };

  const msg = getMessage();

  return (
    <div className="w-full max-w-md mx-auto animate-fade-in">
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-8 text-center space-y-6">
        <div>
          <div className="text-5xl mb-3">{msg.emoji}</div>
          <h2 className="text-2xl font-bold text-white">จบเกม</h2>
          <p className="text-gray-400 mt-1">{msg.text}</p>
        </div>

        {/* Score */}
        <div className="bg-violet-500/10 border border-violet-500/20 rounded-xl p-5">
          <p className="text-xs text-violet-400 uppercase tracking-wider font-medium mb-1">
            คะแนนที่ได้
          </p>
          <p className="text-5xl font-black text-violet-400 tabular-nums">
            {finalScore}
          </p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-4 gap-3">
          <div className="bg-white/5 rounded-xl p-3">
            <p className="text-[10px] text-gray-500 uppercase tracking-wider">คำถาม</p>
            <p className="text-xl font-bold text-white tabular-nums">{finalAnswered}</p>
          </div>
          <div className="bg-white/5 rounded-xl p-3">
            <p className="text-[10px] text-gray-500 uppercase tracking-wider">ถูกต้อง</p>
            <p className="text-xl font-bold text-violet-400 tabular-nums">{finalCorrect}</p>
          </div>
          <div className="bg-white/5 rounded-xl p-3">
            <p className="text-[10px] text-gray-500 uppercase tracking-wider">ความแม่นยำ</p>
            <p className="text-xl font-bold text-cyan-400 tabular-nums">{accuracy}%</p>
          </div>
          <div className="bg-white/5 rounded-xl p-3">
            <p className="text-[10px] text-gray-500 uppercase tracking-wider">Best Combo</p>
            <p className="text-xl font-bold text-amber-400 tabular-nums">{finalBestCombo}x</p>
          </div>
        </div>

        {/* Difficulty badge */}
        <div className="bg-white/5 rounded-xl px-4 py-2 text-center">
          <span className="text-xs text-gray-500">ระดับ: </span>
          <span className={`text-sm font-bold ${
            difficulty === "hard" ? "text-red-400" :
            difficulty === "medium" ? "text-amber-400" :
            "text-violet-400"
          }`}>{difficultyLabel}</span>
        </div>

        {/* Score saved indicator */}
        {user && scoreSaved && (
          <div className="bg-violet-500/10 border border-violet-500/20 rounded-xl px-4 py-2.5 flex items-center justify-center gap-2">
            <span className="text-violet-400 text-sm font-medium">✓ บันทึกคะแนนเรียบร้อยแล้ว</span>
          </div>
        )}
        {user && saving && (
          <div className="bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 flex items-center justify-center gap-2">
            <span className="text-gray-400 text-sm">กำลังบันทึกคะแนน...</span>
          </div>
        )}
        {!user && (
          <div className="bg-white/5 border border-white/10 rounded-xl px-4 py-2.5">
            <Link
              href="/login"
              className="text-sm text-violet-400 hover:underline"
            >
              เข้าสู่ระบบเพื่อบันทึกคะแนน →
            </Link>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-3 pt-2">
          <button
            onClick={onPlayAgain}
            className="flex-1 bg-violet-500 hover:bg-violet-600 text-white font-bold py-3 rounded-xl transition-all hover:scale-[1.01] active:scale-[0.99]"
          >
            เล่นอีกครั้ง
          </button>
          <button
            onClick={onHome}
            className="flex-1 bg-white/5 hover:bg-white/10 border border-white/10 text-gray-300 font-bold py-3 rounded-xl transition-all"
          >
            หน้าหลัก
          </button>
        </div>

        {/* Leaderboard link */}
        <div className="text-center">
          <Link
            href="/leaderboard"
            className="text-sm text-gray-400 hover:text-violet-400 transition-colors"
          >
            ดูกระดานผู้นำ →
          </Link>
        </div>
      </div>
    </div>
  );
}

"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Navbar from "@/components/Navbar";
import { fetchProfile, fetchRankedStats, fetchRankedHistory } from "@/lib/api";
import { ProfileResponse, RankedStats, RankedMatchHistoryEntry } from "@/types";
import { useAuth } from "@/components/AuthProvider";

export default function ProfileContent({ username }: { username: string }) {
  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [error, setError] = useState("");
  const { user } = useAuth();

  const isOwnProfile = user?.username?.toLowerCase() === username.toLowerCase();

  const [rankedStats, setRankedStats] = useState<RankedStats | null>(null);
  const [rankedHistory, setRankedHistory] = useState<RankedMatchHistoryEntry[]>(
    []
  );

  useEffect(() => {
    setLoading(true);
    setError("");
    fetchProfile(username, page)
      .then((data) => setProfile(data))
      .catch(() => setError("ไม่พบผู้ใช้นี้"))
      .finally(() => setLoading(false));
  }, [username, page]);

  // Fetch ranked stats if it's the user's own profile
  useEffect(() => {
    if (isOwnProfile && user) {
      fetchRankedStats()
        .then(setRankedStats)
        .catch(() => setRankedStats(null));
      fetchRankedHistory()
        .then(setRankedHistory)
        .catch(() => setRankedHistory([]));
    }
  }, [isOwnProfile, user]);

  const totalPages = profile ? Math.ceil(profile.total / 10) : 0;

  const formatTime = (seconds: number) => {
    if (seconds === 0) return "-";
    return seconds.toFixed(2);
  };

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString("th-TH", {
      day: "2-digit",
      month: "2-digit",
      year: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
        <Navbar />
        <div className="flex justify-center py-32">
          <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
        </div>
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
        <Navbar />
        <div className="max-w-4xl mx-auto px-4 sm:px-6 pt-20 text-center">
          <div className="text-6xl mb-4">😕</div>
          <h1 className="text-2xl font-bold mb-2">ไม่พบผู้ใช้นี้</h1>
          <p className="text-gray-400 mb-6">
            ผู้ใช้ &quot;{username}&quot; ไม่มีอยู่ในระบบ
          </p>
          <Link
            href="/"
            className="text-violet-400 hover:underline"
          >
            กลับหน้าหลัก
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-5xl mx-auto px-4 sm:px-6 pt-8 pb-16 pb-safe">
        {/* Profile Header & Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6 mb-8">
          {/* Profile Card */}
          <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
            <h2 className="text-sm text-gray-500 font-medium mb-4">โปรไฟล์</h2>
            <div className="flex flex-col items-center text-center">
              <div className="w-20 h-20 rounded-full bg-gradient-to-br from-violet-400 to-purple-400 flex items-center justify-center text-3xl font-bold text-white mb-4">
                {profile.user.username.charAt(0).toUpperCase()}
              </div>
              <h1 className="text-2xl font-bold text-white mb-1">
                {profile.user.username}
              </h1>
              {isOwnProfile && (
                <span className="text-xs text-violet-400/60 mb-3">(คุณ)</span>
              )}
            </div>
          </div>

          {/* Stats Cards */}
          <div className="md:col-span-2 grid grid-cols-2 gap-3 sm:gap-4">
            {/* High Score */}
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-6 flex items-center gap-3 sm:gap-4">
              <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-amber-500/10 flex items-center justify-center text-xl sm:text-2xl shrink-0">
                🏅
              </div>
              <div className="min-w-0">
                <p className="text-[10px] sm:text-xs text-gray-500 font-medium">คะแนนสูงสุดต่อเกม</p>
                <p className="text-xl sm:text-3xl font-black text-white">
                  {profile.stats.highScore}
                  <span className="text-[10px] sm:text-sm text-gray-500 font-normal ml-1">คะแนน</span>
                </p>
              </div>
            </div>

            {/* Average Time */}
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-6 flex items-center gap-3 sm:gap-4">
              <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-blue-500/10 flex items-center justify-center text-xl sm:text-2xl shrink-0">
                ⏱️
              </div>
              <div className="min-w-0">
                <p className="text-[10px] sm:text-xs text-gray-500 font-medium">เวลาเฉลี่ย</p>
                <p className="text-xl sm:text-3xl font-black text-white">
                  {formatTime(profile.stats.averageTime)}
                  <span className="text-[10px] sm:text-sm text-gray-500 font-normal ml-1">วินาที</span>
                </p>
              </div>
            </div>

            {/* Total Points */}
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-6 flex items-center gap-3 sm:gap-4">
              <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-violet-500/10 flex items-center justify-center text-xl sm:text-2xl shrink-0">
                🏆
              </div>
              <div className="min-w-0">
                <p className="text-[10px] sm:text-xs text-gray-500 font-medium">คะแนนรวม</p>
                <p className="text-xl sm:text-3xl font-black text-violet-400">
                  {profile.stats.totalPoints.toLocaleString()}
                  <span className="text-[10px] sm:text-sm text-gray-500 font-normal ml-1">คะแนน</span>
                </p>
              </div>
            </div>

            {/* Total Games */}
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-6 flex items-center gap-3 sm:gap-4">
              <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-purple-500/10 flex items-center justify-center text-xl sm:text-2xl shrink-0">
                🎮
              </div>
              <div className="min-w-0">
                <p className="text-[10px] sm:text-xs text-gray-500 font-medium">เล่นไปทั้งหมด</p>
                <p className="text-xl sm:text-3xl font-black text-white">
                  {profile.stats.totalGames}
                  <span className="text-[10px] sm:text-sm text-gray-500 font-normal ml-1">ครั้ง</span>
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Ranked Stats (own profile only) */}
        {isOwnProfile && rankedStats && (
          <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6 mb-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold text-white">⚔️ สถิติ Ranked</h2>
              <Link
                href="/play/ranked"
                className="text-xs text-blue-400 hover:underline"
              >
                เล่น Ranked →
              </Link>
            </div>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-2xl font-extrabold text-blue-400">
                  {rankedStats.rating}
                </div>
                <div className="text-xs text-gray-500 mt-1">Rating</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-extrabold text-violet-400">
                  {rankedStats.highestRating}
                </div>
                <div className="text-xs text-gray-500 mt-1">สูงสุด</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-extrabold text-green-400">
                  {rankedStats.rankedWins}
                </div>
                <div className="text-xs text-gray-500 mt-1">ชนะ</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-extrabold text-red-400">
                  {rankedStats.rankedLosses}
                </div>
                <div className="text-xs text-gray-500 mt-1">แพ้</div>
              </div>
            </div>

            {/* Ranked match history */}
            {rankedHistory.length > 0 && (
              <div className="mt-6 border-t border-white/5 pt-4">
                <h3 className="text-sm text-gray-500 font-medium mb-3">
                  ประวัติ Ranked ล่าสุด
                </h3>
                <div className="space-y-2">
                  {rankedHistory.slice(0, 5).map((match) => {
                    const opponent =
                      match.player1Name === profile?.user.username
                        ? match.player2Name
                        : match.player1Name;
                    const myScore =
                      match.player1Name === profile?.user.username
                        ? match.player1Score
                        : match.player2Score;
                    const theirScore =
                      match.player1Name === profile?.user.username
                        ? match.player2Score
                        : match.player1Score;
                    return (
                      <div
                        key={match.matchId}
                        className={`flex items-center justify-between px-4 py-2 rounded-lg ${
                          match.won
                            ? "bg-green-500/5 border border-green-500/10"
                            : "bg-red-500/5 border border-red-500/10"
                        }`}
                      >
                        <div className="flex items-center gap-2">
                          <span
                            className={`text-xs font-bold px-2 py-0.5 rounded ${
                              match.won
                                ? "bg-green-500/20 text-green-400"
                                : "bg-red-500/20 text-red-400"
                            }`}
                          >
                            {match.won ? "ชนะ" : "แพ้"}
                          </span>
                          <span className="text-sm text-gray-300">
                            vs{" "}
                            <Link
                              href={`/profile/${opponent}`}
                              className="text-white hover:underline"
                            >
                              {opponent}
                            </Link>
                          </span>
                        </div>
                        <div className="text-sm tabular-nums">
                          <span className="text-white font-bold">{myScore}</span>
                          <span className="text-gray-600 mx-1">-</span>
                          <span className="text-gray-400">{theirScore}</span>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Game History */}
        <div className="bg-[#12122a] border border-white/5 rounded-2xl overflow-hidden">
          <div className="flex items-center justify-between px-4 sm:px-6 py-4 border-b border-white/5">
            <h2 className="text-lg font-bold text-white">ประวัติการเล่น</h2>
            <p className="text-xs text-gray-500">
              แสดง {profile.history.length} จาก {profile.total} รายการ
            </p>
          </div>

          {profile.history.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <div className="text-4xl mb-4">🎮</div>
              <p className="text-lg font-medium">ยังไม่มีประวัติการเล่น</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full min-w-[480px]">
                  <thead>
                    <tr className="border-b border-white/5 text-xs text-gray-500 uppercase tracking-wider">
                      <th className="text-left px-4 sm:px-6 py-3">วัน-เวลา</th>
                      <th className="text-right px-4 sm:px-6 py-3">คะแนน</th>
                      <th className="text-right px-4 sm:px-6 py-3">ถูก/ทั้งหมด</th>
                      <th className="text-right px-4 sm:px-6 py-3">ความแม่นยำ</th>
                      <th className="text-right px-4 sm:px-6 py-3">เวลาที่ใช้</th>
                    </tr>
                  </thead>
                  <tbody>
                    {profile.history.map((game) => {
                      const accuracy =
                        game.totalAnswered > 0
                          ? Math.round(
                              (game.correctAnswers / game.totalAnswered) * 100
                            )
                          : 0;
                      return (
                        <tr
                          key={game.id}
                          className="border-b border-white/5 last:border-0 hover:bg-white/[0.02] transition-colors"
                        >
                          <td className="px-4 sm:px-6 py-3 text-sm text-gray-400 whitespace-nowrap">
                            {formatDate(game.playedAt)}
                          </td>
                          <td className="px-4 sm:px-6 py-3 text-right">
                            <span className="font-bold text-violet-400 tabular-nums">
                              {game.score}
                            </span>
                          </td>
                          <td className="px-4 sm:px-6 py-3 text-right text-sm text-gray-400 tabular-nums">
                            {game.correctAnswers}/{game.totalAnswered}
                          </td>
                          <td className="px-4 sm:px-6 py-3 text-right text-sm tabular-nums">
                            <span
                              className={
                                accuracy >= 70
                                  ? "text-violet-400"
                                  : accuracy >= 50
                                  ? "text-amber-400"
                                  : "text-red-400"
                              }
                            >
                              {accuracy}%
                            </span>
                          </td>
                          <td className="px-4 sm:px-6 py-3 text-right text-sm text-gray-400 tabular-nums whitespace-nowrap">
                            {game.timeSpent > 0
                              ? `${game.timeSpent.toFixed(2)} วินาที`
                              : `${game.timeLimit} วินาที`}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-center gap-2 px-6 py-4 border-t border-white/5">
                  <button
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-gray-400 hover:bg-white/10 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                  >
                    ← ก่อนหน้า
                  </button>
                  {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                    const pageNum = i + 1;
                    return (
                      <button
                        key={pageNum}
                        onClick={() => setPage(pageNum)}
                        className={`w-8 h-8 text-sm rounded-lg font-bold transition-all ${
                          page === pageNum
                            ? "bg-violet-500 text-white"
                            : "bg-white/5 text-gray-400 hover:bg-white/10"
                        }`}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                  <button
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                    className="px-3 py-1.5 text-sm rounded-lg bg-white/5 text-gray-400 hover:bg-white/10 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                  >
                    ถัดไป →
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}

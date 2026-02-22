"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Navbar from "@/components/Navbar";
import { fetchLeaderboard, fetchRankedLeaderboard } from "@/lib/api";
import { LeaderboardEntry, RankedLeaderboardEntry } from "@/types";
import { useAuth } from "@/components/AuthProvider";

type Tab = "casual" | "ranked";

export default function LeaderboardPage() {
  const [tab, setTab] = useState<Tab>("casual");
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [rankedEntries, setRankedEntries] = useState<RankedLeaderboardEntry[]>(
    []
  );
  const [loading, setLoading] = useState(true);
  const { user } = useAuth();

  useEffect(() => {
    setLoading(true);
    if (tab === "casual") {
      fetchLeaderboard()
        .then((data) => setEntries(data.entries))
        .catch(() => setEntries([]))
        .finally(() => setLoading(false));
    } else {
      fetchRankedLeaderboard()
        .then((data) => setRankedEntries(data))
        .catch(() => setRankedEntries([]))
        .finally(() => setLoading(false));
    }
  }, [tab]);

  const getRankDisplay = (rank: number) => {
    if (rank === 1) return "🥇";
    if (rank === 2) return "🥈";
    if (rank === 3) return "🥉";
    return `#${rank}`;
  };

  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 pt-8 sm:pt-12 pb-16 pb-safe">
        {/* Header */}
        <div className="text-center mb-10">
          <h1 className="text-4xl font-extrabold mb-2">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
              กระดานผู้นำ
            </span>
          </h1>
          <p className="text-gray-400">
            อันดับคะแนนรวมของผู้เล่น FastChem
          </p>
        </div>

        {/* Tabs */}
        <div className="flex justify-center gap-2 mb-6">
          <button
            onClick={() => setTab("casual")}
            className={`px-5 py-2 rounded-lg text-sm font-medium transition-all cursor-pointer ${
              tab === "casual"
                ? "bg-violet-500/20 text-violet-400 border border-violet-500/30"
                : "text-gray-500 hover:text-gray-300 border border-transparent"
            }`}
          >
            🧪 ฝึกซ้อม
          </button>
          <button
            onClick={() => setTab("ranked")}
            className={`px-5 py-2 rounded-lg text-sm font-medium transition-all cursor-pointer ${
              tab === "ranked"
                ? "bg-blue-500/20 text-blue-400 border border-blue-500/30"
                : "text-gray-500 hover:text-gray-300 border border-transparent"
            }`}
          >
            ⚔️ Ranked
          </button>
        </div>

        {/* Table */}
        <div className="bg-[#12122a] border border-white/5 rounded-2xl overflow-hidden">
          {loading ? (
            <div className="flex justify-center py-16">
              <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
            </div>
          ) : tab === "casual" ? (
            /* Casual leaderboard */
            entries.length === 0 ? (
              <div className="text-center py-16 text-gray-500">
                <div className="text-4xl mb-4">🏆</div>
                <p className="text-lg font-medium">ยังไม่มีคะแนน</p>
                <p className="text-sm mt-1">เป็นคนแรกที่เล่นและทำสถิติ!</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
              <table className="w-full min-w-[400px]">
                <thead>
                  <tr className="border-b border-white/5 text-xs text-gray-500 uppercase tracking-wider">
                    <th className="text-left px-4 sm:px-6 py-4 w-12 sm:w-16">อันดับ</th>
                    <th className="text-left px-4 sm:px-6 py-4">ผู้เล่น</th>
                    <th className="text-right px-4 sm:px-6 py-4">คะแนนรวม</th>
                    <th className="text-right px-4 sm:px-6 py-4 hidden sm:table-cell">จำนวนเกม</th>
                  </tr>
                </thead>
                <tbody>
                  {entries.map((entry) => (
                    <tr
                      key={entry.userId}
                      className={`border-b border-white/5 last:border-0 transition-colors ${
                        user && entry.userId === user.id
                          ? "bg-violet-500/5"
                          : "hover:bg-white/[0.02]"
                      }`}
                    >
                      <td className="px-4 sm:px-6 py-4 font-bold text-lg">
                        {getRankDisplay(entry.rank)}
                      </td>
                      <td className="px-4 sm:px-6 py-4">
                        <Link
                          href={`/profile/${entry.username}`}
                          className={`font-semibold hover:underline truncate block max-w-[120px] sm:max-w-none ${
                            user && entry.userId === user.id
                              ? "text-violet-400"
                              : "text-white"
                          }`}
                        >
                          {entry.username}
                          {user && entry.userId === user.id && (
                            <span className="ml-2 text-xs text-violet-400/60">
                              (คุณ)
                            </span>
                          )}
                        </Link>
                      </td>
                      <td className="px-4 sm:px-6 py-4 text-right">
                        <span className="text-base sm:text-lg font-bold text-violet-400 tabular-nums">
                          {entry.totalPoints.toLocaleString()}
                        </span>
                        <span className="text-xs text-gray-500 ml-1 hidden sm:inline">คะแนน</span>
                      </td>
                      <td className="px-4 sm:px-6 py-4 text-right hidden sm:table-cell">
                        <span className="text-sm text-gray-400 tabular-nums">
                          {entry.totalGames} ครั้ง
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              </div>
            )
          ) : (
            /* Ranked leaderboard */
            rankedEntries.length === 0 ? (
              <div className="text-center py-16 text-gray-500">
                <div className="text-4xl mb-4">⚔️</div>
                <p className="text-lg font-medium">ยังไม่มี Ranked</p>
                <p className="text-sm mt-1">เป็นคนแรกที่เล่น Ranked และทำสถิติ!</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
              <table className="w-full min-w-[400px]">
                <thead>
                  <tr className="border-b border-white/5 text-xs text-gray-500 uppercase tracking-wider">
                    <th className="text-left px-4 sm:px-6 py-4 w-12 sm:w-16">อันดับ</th>
                    <th className="text-left px-4 sm:px-6 py-4">ผู้เล่น</th>
                    <th className="text-right px-4 sm:px-6 py-4">Rating</th>
                    <th className="text-right px-4 sm:px-6 py-4 hidden sm:table-cell">สถิติ</th>
                    <th className="text-right px-4 sm:px-6 py-4 hidden md:table-cell">สูงสุด</th>
                  </tr>
                </thead>
                <tbody>
                  {rankedEntries.map((entry) => {
                    const total = entry.wins + entry.losses;
                    const winRate =
                      total > 0
                        ? Math.round((entry.wins / total) * 100)
                        : 0;
                    return (
                      <tr
                        key={entry.userId}
                        className={`border-b border-white/5 last:border-0 transition-colors ${
                          user && entry.userId === user.id
                            ? "bg-blue-500/5"
                            : "hover:bg-white/[0.02]"
                        }`}
                      >
                        <td className="px-6 py-4 font-bold text-lg">
                          {getRankDisplay(entry.rank)}
                        </td>
                        <td className="px-6 py-4">
                          <Link
                            href={`/profile/${entry.username}`}
                            className={`font-semibold hover:underline ${
                              user && entry.userId === user.id
                                ? "text-blue-400"
                                : "text-white"
                            }`}
                          >
                            {entry.username}
                            {user && entry.userId === user.id && (
                              <span className="ml-2 text-xs text-blue-400/60">
                                (คุณ)
                              </span>
                            )}
                          </Link>
                        </td>
                        <td className="px-6 py-4 text-right">
                          <span className="text-lg font-bold text-blue-400 tabular-nums">
                            {entry.rating}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-right hidden sm:table-cell">
                          <span className="text-sm tabular-nums">
                            <span className="text-green-400">{entry.wins}W</span>
                            <span className="text-gray-600 mx-1">/</span>
                            <span className="text-red-400">{entry.losses}L</span>
                          </span>
                          <span className="text-xs text-gray-500 ml-2">
                            ({winRate}%)
                          </span>
                        </td>
                        <td className="px-6 py-4 text-right hidden md:table-cell">
                          <span className="text-sm text-violet-400 tabular-nums">
                            {entry.highestRating}
                          </span>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
              </div>
            )
          )}
        </div>
      </div>
    </div>
  );
}

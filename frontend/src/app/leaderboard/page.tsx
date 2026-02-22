"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Navbar from "@/components/Navbar";
import { fetchLeaderboard } from "@/lib/api";
import { LeaderboardEntry } from "@/types";
import { useAuth } from "@/components/AuthProvider";

export default function LeaderboardPage() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const { user } = useAuth();

  useEffect(() => {
    setLoading(true);
    fetchLeaderboard()
      .then((data) => setEntries(data.entries))
      .catch(() => setEntries([]))
      .finally(() => setLoading(false));
  }, []);

  const getRankDisplay = (rank: number) => {
    if (rank === 1) return "🥇";
    if (rank === 2) return "🥈";
    if (rank === 3) return "🥉";
    return `#${rank}`;
  };

  return (
    <div className="min-h-screen bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-4xl mx-auto px-6 pt-12 pb-16">
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

        {/* Table */}
        <div className="bg-[#12122a] border border-white/5 rounded-2xl overflow-hidden">
          {loading ? (
            <div className="flex justify-center py-16">
              <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
            </div>
          ) : entries.length === 0 ? (
            <div className="text-center py-16 text-gray-500">
              <div className="text-4xl mb-4">🏆</div>
              <p className="text-lg font-medium">ยังไม่มีคะแนน</p>
              <p className="text-sm mt-1">เป็นคนแรกที่เล่นและทำสถิติ!</p>
            </div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/5 text-xs text-gray-500 uppercase tracking-wider">
                  <th className="text-left px-6 py-4 w-16">อันดับ</th>
                  <th className="text-left px-6 py-4">ผู้เล่น</th>
                  <th className="text-right px-6 py-4">คะแนนรวม</th>
                  <th className="text-right px-6 py-4 hidden sm:table-cell">จำนวนเกม</th>
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
                    <td className="px-6 py-4 font-bold text-lg">
                      {getRankDisplay(entry.rank)}
                    </td>
                    <td className="px-6 py-4">
                      <Link
                        href={`/profile/${entry.username}`}
                        className={`font-semibold hover:underline ${
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
                    <td className="px-6 py-4 text-right">
                      <span className="text-lg font-bold text-violet-400 tabular-nums">
                        {entry.totalPoints.toLocaleString()}
                      </span>
                      <span className="text-xs text-gray-500 ml-1">คะแนน</span>
                    </td>
                    <td className="px-6 py-4 text-right hidden sm:table-cell">
                      <span className="text-sm text-gray-400 tabular-nums">
                        {entry.totalGames} ครั้ง
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}

"use client";

import Link from "next/link";
import { AdminLeaderboards } from "@/types";
import { Empty, Panel, Spinner, Table, formatNumber } from "./ui";

interface LeaderboardPanelProps {
  boards: AdminLeaderboards | null;
  loading: boolean;
}

// Read live, not from the public board's 30-second cache — an operator
// checking a correction needs the current number, not the cached one.
export default function LeaderboardPanel({ boards, loading }: LeaderboardPanelProps) {
  if (loading) return <Spinner />;

  return (
    <div className="grid lg:grid-cols-2 gap-4">
      <Panel title="กระดานผู้นำ (ฝึกซ้อม)" subtitle="เรียงตามคะแนนรวม · อ่านสด ไม่ผ่านแคช">
        {!boards || boards.casual.length === 0 ? (
          <Empty icon="🏆" text="ยังไม่มีคะแนน" />
        ) : (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3 w-12">#</th>
                <th className="text-left px-3 py-3">ผู้เล่น</th>
                <th className="text-right px-3 py-3">คะแนนรวม</th>
                <th className="text-right px-4 sm:px-5 py-3">เกม</th>
              </>
            }
          >
            {boards.casual.map((e) => (
              <tr key={e.userId} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 text-gray-500 tabular-nums">{e.rank}</td>
                <td className="px-3 py-2.5">
                  <Link
                    href={`/profile/${e.username}/`}
                    className="text-gray-100 hover:text-violet-400"
                  >
                    {e.username}
                  </Link>
                </td>
                <td className="px-3 py-2.5 text-right tabular-nums text-violet-300">
                  {formatNumber(e.totalPoints)}
                </td>
                <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                  {e.totalGames}
                </td>
              </tr>
            ))}
          </Table>
        )}
      </Panel>

      <Panel title="กระดานผู้นำ (Ranked)" subtitle="เรียงตาม rating">
        {!boards || boards.ranked.length === 0 ? (
          <Empty icon="⚔️" text="ยังไม่มีผู้เล่น Ranked" />
        ) : (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3 w-12">#</th>
                <th className="text-left px-3 py-3">ผู้เล่น</th>
                <th className="text-right px-3 py-3">Rating</th>
                <th className="text-right px-3 py-3">W/L</th>
                <th className="text-right px-4 sm:px-5 py-3">สูงสุด</th>
              </>
            }
          >
            {boards.ranked.map((e) => (
              <tr key={e.userId} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 text-gray-500 tabular-nums">{e.rank}</td>
                <td className="px-3 py-2.5">
                  <Link
                    href={`/profile/${e.username}/`}
                    className="text-gray-100 hover:text-violet-400"
                  >
                    {e.username}
                  </Link>
                </td>
                <td className="px-3 py-2.5 text-right tabular-nums text-blue-300">{e.rating}</td>
                <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">
                  <span className="text-emerald-400">{e.wins}</span>
                  <span className="text-gray-600">/</span>
                  <span className="text-red-400">{e.losses}</span>
                </td>
                <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                  {e.highestRating}
                </td>
              </tr>
            ))}
          </Table>
        )}
      </Panel>
    </div>
  );
}

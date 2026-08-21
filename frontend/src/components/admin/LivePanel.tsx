"use client";

import { AdminLive, AdminMatchesResponse } from "@/types";
import {
  Empty,
  Panel,
  Pill,
  Spinner,
  Stat,
  Table,
  formatDateTime,
  formatDuration,
  formatNumber,
} from "./ui";

interface LivePanelProps {
  live: AdminLive | null;
  matches: AdminMatchesResponse | null;
  loading: boolean;
  autoRefresh: boolean;
  onAutoRefresh: (on: boolean) => void;
}

const DIFFICULTY_LABELS: Record<string, string> = {
  easy: "ง่าย",
  medium: "ปานกลาง",
  hard: "ยาก",
};

export default function LivePanel({
  live,
  matches,
  loading,
  autoRefresh,
  onAutoRefresh,
}: LivePanelProps) {
  if (loading && !live) return <Spinner />;

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 lg:grid-cols-5 gap-3">
        <Stat label="ในคิว" value={live?.queue.length ?? 0} />
        <Stat label="แมตช์ที่กำลังแข่ง" value={live?.matches.length ?? 0} />
        <Stat label="ห้องที่เปิดอยู่" value={live?.rooms.length ?? 0} />
        <Stat label="เกมเดี่ยวที่ค้างอยู่" value={live?.soloMatches ?? 0} />
        <Stat label="คำถามที่รอคำตอบ" value={live?.storedQuestions ?? 0} />
      </div>

      <Panel
        title="คิว Ranked"
        subtitle="รอจับคู่ตาม rating"
        action={
          <label className="flex items-center gap-1.5 text-xs text-gray-400 cursor-pointer">
            <input
              type="checkbox"
              checked={autoRefresh}
              onChange={(e) => onAutoRefresh(e.target.checked)}
              className="accent-violet-500 cursor-pointer"
            />
            รีเฟรชอัตโนมัติ
          </label>
        }
      >
        {!live || live.queue.length === 0 ? (
          <Empty icon="🎯" text="ไม่มีใครอยู่ในคิว" />
        ) : (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3">ผู้เล่น</th>
                <th className="text-right px-3 py-3">Rating</th>
                <th className="text-right px-4 sm:px-5 py-3">รอมาแล้ว</th>
              </>
            }
          >
            {live.queue.map((q) => (
              <tr key={q.userId} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 text-gray-200">{q.username}</td>
                <td className="px-3 py-2.5 text-right tabular-nums">{q.rating}</td>
                <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                  {formatDuration(q.waitingSeconds)}
                </td>
              </tr>
            ))}
          </Table>
        )}
      </Panel>

      <Panel title="แมตช์ที่กำลังแข่ง" subtitle="Ranked และห้องส่วนตัว (matchId ติดลบคือห้อง)">
        {!live || live.matches.length === 0 ? (
          <Empty icon="⚔️" text="ยังไม่มีแมตช์ที่กำลังแข่ง" />
        ) : (
          <div className="space-y-3">
            {live.matches.map((m) => (
              <div key={m.matchId} className="border border-white/5 rounded-xl p-4">
                <div className="flex items-center justify-between text-xs text-gray-500 mb-3">
                  <span>
                    #{Math.abs(m.matchId)}{" "}
                    <Pill tone={m.matchId < 0 ? "blue" : "violet"}>
                      {m.matchId < 0 ? "ห้อง" : "ranked"}
                    </Pill>
                  </span>
                  <span>
                    ข้อที่ {m.question + 1} · {m.status} · เริ่ม {formatDateTime(m.createdAt)}
                  </span>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  {[m.player1, m.player2].map((p) => (
                    <div key={p.userId} className="bg-[#0a0a1a]/50 rounded-lg px-3 py-2">
                      <div className="flex items-center gap-2">
                        <span
                          className={`w-1.5 h-1.5 rounded-full ${
                            p.connected ? "bg-emerald-400" : "bg-red-400"
                          }`}
                          title={p.connected ? "เชื่อมต่ออยู่" : "หลุดการเชื่อมต่อ"}
                        />
                        <span className="text-sm text-gray-200 truncate">{p.username}</span>
                        <span className="text-[11px] text-gray-600">{p.rating}</span>
                      </div>
                      <div className="text-xs text-gray-500 mt-1 tabular-nums">
                        {formatNumber(p.totalScore)} คะแนน · ตอบ {p.answered} · ถูก {p.correct} ·
                        คอมโบ {p.combo}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </Panel>

      <Panel title="ห้องส่วนตัว" subtitle="รหัส 6 หลักที่ยังไม่หมดอายุ">
        {!live || live.rooms.length === 0 ? (
          <Empty icon="🚪" text="ไม่มีห้องที่เปิดอยู่" />
        ) : (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3">รหัส</th>
                <th className="text-left px-3 py-3">เจ้าของห้อง</th>
                <th className="text-left px-3 py-3">ผู้เข้าร่วม</th>
                <th className="text-left px-3 py-3">สถานะ</th>
                <th className="text-right px-4 sm:px-5 py-3">สร้างเมื่อ</th>
              </>
            }
          >
            {live.rooms.map((r) => (
              <tr key={r.code} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 font-mono text-violet-300">{r.code}</td>
                <td className="px-3 py-2.5 text-gray-200">{r.hostName}</td>
                <td className="px-3 py-2.5 text-gray-400">{r.guestName || "—"}</td>
                <td className="px-3 py-2.5">
                  <Pill tone={r.status === "active" ? "green" : "muted"}>{r.status}</Pill>
                </td>
                <td className="px-4 sm:px-5 py-2.5 text-right text-xs text-gray-500">
                  {formatDateTime(r.createdAt)}
                </td>
              </tr>
            ))}
          </Table>
        )}
      </Panel>

      <Panel title="แมตช์ Ranked ล่าสุด" subtitle="ผลที่บันทึกลงฐานข้อมูลแล้ว">
        {!matches || matches.ranked.length === 0 ? (
          <Empty icon="🏁" text="ยังไม่มีแมตช์ที่จบ" />
        ) : (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3">แมตช์</th>
                <th className="text-left px-3 py-3">ผู้เล่น</th>
                <th className="text-right px-3 py-3">คะแนน</th>
                <th className="text-left px-3 py-3">ผู้ชนะ</th>
                <th className="text-right px-4 sm:px-5 py-3">จบเมื่อ</th>
              </>
            }
          >
            {matches.ranked.map((m) => (
              <tr key={m.matchId} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 text-gray-500">#{m.matchId}</td>
                <td className="px-3 py-2.5 text-gray-200">
                  {m.player1} <span className="text-gray-600">vs</span> {m.player2}
                </td>
                <td className="px-3 py-2.5 text-right tabular-nums">
                  {m.player1Score} : {m.player2Score}
                </td>
                <td className="px-3 py-2.5">
                  {m.winner ? (
                    <span className="text-emerald-400">{m.winner}</span>
                  ) : (
                    <span className="text-gray-600">เสมอ/ไม่จบ</span>
                  )}
                </td>
                <td className="px-4 sm:px-5 py-2.5 text-right text-xs text-gray-500">
                  {formatDateTime(m.finishedAt)}
                </td>
              </tr>
            ))}
          </Table>
        )}
      </Panel>

      <div className="grid lg:grid-cols-2 gap-4">
        <Panel title="แมตช์เดี่ยวล่าสุด" subtitle="เกมที่เซิร์ฟเวอร์ให้คะแนนเอง">
          {!matches || matches.solo.length === 0 ? (
            <Empty icon="🧪" text="ยังไม่มีแมตช์เดี่ยว" />
          ) : (
            <Table
              head={
                <>
                  <th className="text-left px-4 sm:px-5 py-3">ผู้เล่น</th>
                  <th className="text-left px-3 py-3">ระดับ</th>
                  <th className="text-right px-3 py-3">คะแนน</th>
                  <th className="text-right px-4 sm:px-5 py-3">ข้อ</th>
                </>
              }
            >
              {matches.solo.map((m) => (
                <tr key={m.matchId} className="border-b border-white/5 last:border-0">
                  <td className="px-4 sm:px-5 py-2.5 text-gray-200">{m.username}</td>
                  <td className="px-3 py-2.5">
                    <Pill tone="muted">{DIFFICULTY_LABELS[m.difficulty] ?? m.difficulty}</Pill>
                  </td>
                  <td className="px-3 py-2.5 text-right tabular-nums">
                    {formatNumber(m.totalScore)}
                    <span className="text-[11px] text-gray-600 ml-1">x{m.bestCombo}</span>
                  </td>
                  <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                    {m.attempts}
                  </td>
                </tr>
              ))}
            </Table>
          )}
        </Panel>

        <Panel title="คะแนนที่ส่งล่าสุด" subtitle="จากโหมดฝึกซ้อม">
          {!matches || matches.scores.length === 0 ? (
            <Empty icon="📝" text="ยังไม่มีคะแนน" />
          ) : (
            <Table
              head={
                <>
                  <th className="text-left px-4 sm:px-5 py-3">ผู้เล่น</th>
                  <th className="text-right px-3 py-3">คะแนน</th>
                  <th className="text-right px-3 py-3">ถูก/ตอบ</th>
                  <th className="text-right px-4 sm:px-5 py-3">เมื่อ</th>
                </>
              }
            >
              {matches.scores.map((s) => (
                <tr key={s.id} className="border-b border-white/5 last:border-0">
                  <td className="px-4 sm:px-5 py-2.5 text-gray-200">{s.username}</td>
                  <td className="px-3 py-2.5 text-right tabular-nums">{formatNumber(s.score)}</td>
                  <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">
                    {s.correctAnswers}/{s.totalAnswered}
                  </td>
                  <td className="px-4 sm:px-5 py-2.5 text-right text-xs text-gray-500">
                    {formatDateTime(s.playedAt)}
                  </td>
                </tr>
              ))}
            </Table>
          )}
        </Panel>
      </div>
    </div>
  );
}

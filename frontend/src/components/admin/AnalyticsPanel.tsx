"use client";

import { AdminAnalytics, AdminOverview } from "@/types";
import {
  BarList,
  Columns,
  Empty,
  Panel,
  Pill,
  SERIES,
  Stat,
  Table,
  formatBytes,
  formatDuration,
  formatNumber,
} from "./ui";

interface AnalyticsPanelProps {
  overview: AdminOverview | null;
  analytics: AdminAnalytics | null;
  days: number;
  onDaysChange: (days: number) => void;
}

const DIFFICULTY_LABELS: Record<string, string> = {
  easy: "ง่าย",
  medium: "ปานกลาง",
  hard: "ยาก",
};

// The daily series is drawn as four small multiples rather than one chart with
// four lines: the measures have different units and scales, and a second y-axis
// is never the answer.
export default function AnalyticsPanel({
  overview,
  analytics,
  days,
  onDaysChange,
}: AnalyticsPanelProps) {
  const daily = analytics?.daily ?? [];
  const shortDate = (iso: string) => iso.slice(5).replace("-", "/");

  return (
    <div className="space-y-4">
      {overview && (
        <>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <Stat label="ผู้เล่นทั้งหมด" value={overview.users} hint={`ใหม่ 7 วัน +${overview.newUsers7d}`} />
            <Stat
              label="เล่นแล้ว"
              value={overview.games}
              hint={`24 ชม. ${formatNumber(overview.games24h)} เกม`}
            />
            <Stat label="ข้อที่ตอบไปแล้ว" value={overview.questionsAnswered} />
            <Stat
              label="ความแม่นยำเฉลี่ย"
              value={`${overview.averageAccuracy.toFixed(1)}%`}
              hint={`คะแนนรวม ${formatNumber(overview.totalPoints)}`}
            />
          </div>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <Stat label="แมตช์ Ranked" value={overview.rankedMatches} />
            <Stat label="แมตช์เดี่ยว" value={overview.soloMatches} />
            <Stat
              label="ผู้เล่นที่ยังเล่นอยู่ (7 วัน)"
              value={overview.activePlayers7d}
              tone="good"
            />
            <Stat
              label="Anti-cheat findings"
              value={overview.findings}
              hint={`24 ชม. ${formatNumber(overview.findings24h)}`}
              tone={overview.findings24h > 0 ? "warn" : "default"}
            />
          </div>
        </>
      )}

      <Panel
        title="กิจกรรมรายวัน"
        subtitle={`ช่วง ${days} วันล่าสุด (UTC)`}
        action={
          <div className="flex gap-1">
            {[7, 30, 90].map((d) => (
              <button
                key={d}
                onClick={() => onDaysChange(d)}
                className={`px-2.5 py-1 rounded-md text-xs font-medium transition-colors cursor-pointer ${
                  days === d
                    ? "bg-violet-500/20 text-violet-300 border border-violet-500/30"
                    : "text-gray-500 hover:text-gray-300 border border-transparent"
                }`}
              >
                {d}ว
              </button>
            ))}
          </div>
        }
      >
        {daily.length === 0 ? (
          <Empty icon="📈" text="ยังไม่มีกิจกรรม" />
        ) : (
          <div className="grid sm:grid-cols-2 gap-6">
            {[
              { title: "เกมที่จบ", color: SERIES[0], pick: (d: (typeof daily)[0]) => d.games },
              { title: "ผู้เล่นไม่ซ้ำ", color: SERIES[1], pick: (d: (typeof daily)[0]) => d.players },
              { title: "แมตช์ Ranked", color: SERIES[2], pick: (d: (typeof daily)[0]) => d.rankedMatches },
              { title: "Findings", color: SERIES[3], pick: (d: (typeof daily)[0]) => d.findings },
            ].map((series) => (
              <div key={series.title}>
                <div className="flex items-center gap-2 mb-2">
                  <span
                    className="w-2.5 h-2.5 rounded-sm"
                    style={{ backgroundColor: series.color }}
                  />
                  <span className="text-xs text-gray-300">{series.title}</span>
                </div>
                <Columns
                  color={series.color}
                  points={daily.map((d) => ({ label: shortDate(d.date), value: series.pick(d) }))}
                />
              </div>
            ))}
          </div>
        )}
      </Panel>

      <div className="grid lg:grid-cols-2 gap-4">
        <Panel title="ตามระดับความยาก" subtitle="จำนวนเกมและความแม่นยำเฉลี่ย">
          <BarList
            items={(analytics?.difficulty ?? []).map((b) => ({
              label: DIFFICULTY_LABELS[b.label] ?? b.label,
              count: b.count,
              note: `${b.value.toFixed(0)}% ถูก`,
            }))}
            unit=" เกม"
          />
        </Panel>

        <Panel title="เวลาที่เล่น" subtitle="จำนวนเกมต่อชั่วโมง (UTC)">
          {analytics ? (
            <Columns
              points={analytics.hours.map((h) => ({ label: `${h.label}:00`, value: h.count }))}
              unit=" เกม"
            />
          ) : (
            <Empty icon="🕐" text="ยังไม่มีข้อมูล" />
          )}
        </Panel>

        <Panel title="การกระจาย Rating" subtitle="ผู้เล่นที่แข่ง Ranked แล้ว">
          {analytics && analytics.ratings.length > 0 ? (
            <Columns
              points={analytics.ratings.map((r) => ({ label: r.label, value: r.count }))}
              unit=" คน"
            />
          ) : (
            <Empty icon="⚔️" text="ยังไม่มีผู้เล่น Ranked" />
          )}
        </Panel>

        <Panel title="กฎที่จับได้บ่อย" subtitle="จำนวน findings ต่อกฎ">
          <BarList
            items={(analytics?.rules ?? []).map((r) => ({ label: r.label, count: r.count }))}
            unit=" ครั้ง"
          />
        </Panel>
      </div>

      <Panel
        title="สถิติรายหัวข้อ"
        subtitle="วัดจากผลของ Ranked ซึ่งเป็นตารางเดียวที่บันทึกหัวข้อของแต่ละข้อ"
      >
        {analytics && analytics.topics.length > 0 ? (
          <Table
            head={
              <>
                <th className="text-left px-4 sm:px-5 py-3">หัวข้อ</th>
                <th className="text-left px-3 py-3">ระดับ</th>
                <th className="text-right px-3 py-3">ตอบ</th>
                <th className="text-right px-3 py-3">ถูก</th>
                <th className="text-left px-3 py-3 w-40">ความแม่นยำ</th>
                <th className="text-right px-4 sm:px-5 py-3">เวลาเฉลี่ย</th>
              </>
            }
          >
            {analytics.topics.map((t) => (
              <tr key={`${t.topic}-${t.difficulty}`} className="border-b border-white/5 last:border-0">
                <td className="px-4 sm:px-5 py-2.5 text-gray-200">{t.topic}</td>
                <td className="px-3 py-2.5">
                  <Pill tone="muted">{DIFFICULTY_LABELS[t.difficulty] ?? t.difficulty}</Pill>
                </td>
                <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">{t.attempts}</td>
                <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">{t.correct}</td>
                <td className="px-3 py-2.5">
                  <div className="flex items-center gap-2">
                    <div className="flex-1 h-2 bg-white/[0.04] rounded-full overflow-hidden">
                      <div
                        className="h-full rounded-full"
                        style={{ width: `${t.accuracy}%`, backgroundColor: SERIES[0] }}
                      />
                    </div>
                    <span className="tabular-nums text-xs text-gray-400 w-10 text-right">
                      {t.accuracy.toFixed(0)}%
                    </span>
                  </div>
                </td>
                <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                  {t.averageTime.toFixed(1)} วิ
                </td>
              </tr>
            ))}
          </Table>
        ) : (
          <Empty icon="🧪" text="ยังไม่มีผลจากแมตช์ Ranked" />
        )}
      </Panel>

      {overview && (
        <Panel title="เซิร์ฟเวอร์" subtitle="สถานะกระบวนการที่กำลังรันอยู่">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
            <Stat label="Uptime" value={formatDuration(overview.uptimeSeconds)} />
            <Stat label="ขนาดฐานข้อมูล" value={formatBytes(overview.databaseBytes)} />
            <Stat label="Goroutines" value={overview.goroutines} />
            <Stat label="หัวข้อคำถาม" value={overview.topics} />
          </div>
        </Panel>
      )}
    </div>
  );
}

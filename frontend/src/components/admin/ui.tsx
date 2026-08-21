"use client";

import { ReactNode } from "react";

// Shared pieces of the admin console: panels, stat tiles, the two chart forms,
// and the formatters every panel needs. Named exports rather than one default
// component per file, because none of these is a screen on its own.

// SERIES is the categorical order. It is fixed: a chart assigns SERIES[0] to
// its first category and counts up, never cycling and never re-assigning when a
// filter changes the number of series. Validated against the #12122a panel
// surface for lightness, chroma, contrast and colour-vision separation — do not
// substitute a step by eye.
export const SERIES = ["#8b5cf6", "#d97706", "#0891b2", "#e11d48", "#15803d"];

// One hue for magnitude. Length carries the value; colour carries nothing.
export const MAGNITUDE = "#8b5cf6";

export function formatNumber(n: number): string {
  return n.toLocaleString("en-US");
}

export function formatBytes(bytes: number): string {
  if (bytes <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** i).toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

export function formatDuration(seconds: number): string {
  if (seconds < 60) return `${Math.round(seconds)} วิ`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)} นาที`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} ชม.`;
  return `${Math.floor(seconds / 86400)} วัน`;
}

export function formatDateTime(value: string | null): string {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "—" : date.toLocaleString("th-TH");
}

interface PanelProps {
  title: string;
  subtitle?: string;
  action?: ReactNode;
  children: ReactNode;
}

export function Panel({ title, subtitle, action, children }: PanelProps) {
  return (
    <section className="bg-[#12122a] border border-white/5 rounded-2xl overflow-hidden">
      <header className="flex items-start justify-between gap-4 px-4 sm:px-5 py-4 border-b border-white/5">
        <div>
          <h2 className="font-semibold text-sm sm:text-base">{title}</h2>
          {subtitle && <p className="text-xs text-gray-500 mt-0.5">{subtitle}</p>}
        </div>
        {action}
      </header>
      <div className="p-4 sm:p-5">{children}</div>
    </section>
  );
}

interface StatProps {
  label: string;
  value: string | number;
  hint?: string;
  tone?: "default" | "warn" | "good";
}

export function Stat({ label, value, hint, tone = "default" }: StatProps) {
  const toneClass =
    tone === "warn" ? "text-amber-400" : tone === "good" ? "text-emerald-400" : "text-white";
  return (
    <div className="bg-[#12122a] border border-white/5 rounded-xl px-4 py-3">
      <div className="text-[11px] uppercase tracking-wider text-gray-500">{label}</div>
      <div className={`text-xl sm:text-2xl font-bold tabular-nums mt-1 ${toneClass}`}>
        {typeof value === "number" ? formatNumber(value) : value}
      </div>
      {hint && <div className="text-xs text-gray-500 mt-0.5">{hint}</div>}
    </div>
  );
}

export function Pill({
  children,
  tone = "muted",
}: {
  children: ReactNode;
  tone?: "muted" | "violet" | "blue" | "green" | "amber" | "red";
}) {
  const tones: Record<string, string> = {
    muted: "bg-white/5 text-gray-400 border-white/10",
    violet: "bg-violet-500/10 text-violet-300 border-violet-500/20",
    blue: "bg-blue-500/10 text-blue-300 border-blue-500/20",
    green: "bg-emerald-500/10 text-emerald-300 border-emerald-500/20",
    amber: "bg-amber-500/10 text-amber-300 border-amber-500/20",
    red: "bg-red-500/10 text-red-300 border-red-500/20",
  };
  return (
    <span className={`inline-block px-2 py-0.5 rounded-md border text-[11px] font-medium ${tones[tone]}`}>
      {children}
    </span>
  );
}

export function Empty({ icon, text }: { icon: string; text: string }) {
  return (
    <div className="text-center py-10 text-gray-500">
      <div className="text-3xl mb-2">{icon}</div>
      <p className="text-sm">{text}</p>
    </div>
  );
}

export function Spinner() {
  return (
    <div className="flex justify-center py-16">
      <div className="w-8 h-8 border-2 border-violet-400 border-t-transparent rounded-full animate-spin" />
    </div>
  );
}

interface BarListItem {
  label: string;
  count: number;
  note?: string;
  color?: string;
}

// BarList ranks categories by magnitude. Every row is direct-labelled, so
// identity never rests on colour alone.
export function BarList({ items, unit }: { items: BarListItem[]; unit?: string }) {
  const max = Math.max(1, ...items.map((i) => i.count));
  if (items.length === 0) return <Empty icon="📊" text="ยังไม่มีข้อมูล" />;
  return (
    <ul className="space-y-2.5">
      {items.map((item, i) => (
        <li key={item.label} title={`${item.label}: ${formatNumber(item.count)}${unit ?? ""}`}>
          <div className="flex items-baseline justify-between gap-3 text-xs mb-1">
            <span className="text-gray-300 truncate">{item.label}</span>
            <span className="tabular-nums text-gray-400 shrink-0">
              {formatNumber(item.count)}
              {unit}
              {item.note && <span className="text-gray-600 ml-1.5">{item.note}</span>}
            </span>
          </div>
          <div className="h-2 bg-white/[0.04] rounded-full overflow-hidden">
            <div
              className="h-full rounded-full"
              style={{
                width: `${Math.max(2, (item.count / max) * 100)}%`,
                backgroundColor: item.color ?? SERIES[i % SERIES.length],
              }}
            />
          </div>
        </li>
      ))}
    </ul>
  );
}

interface ColumnsProps {
  points: { label: string; value: number }[];
  color?: string;
  unit?: string;
}

// Columns is the time/ordinal form: one bar per bucket, one series, a recessive
// baseline, and only the peak and latest values labelled — never every point.
export function Columns({ points, color = MAGNITUDE, unit = "" }: ColumnsProps) {
  const max = Math.max(1, ...points.map((p) => p.value));
  const peakIndex = points.findIndex((p) => p.value === max);
  const last = points.length - 1;
  return (
    <div>
      <div className="flex items-end gap-[2px] h-24 border-b border-white/10">
        {points.map((p) => (
          <div
            key={p.label}
            title={`${p.label}: ${formatNumber(p.value)}${unit}`}
            className="flex-1 min-w-0 relative flex items-end"
            style={{ height: "100%" }}
          >
            <div
              className="w-full rounded-t"
              style={{
                height: `${p.value === 0 ? 1 : Math.max(3, (p.value / max) * 100)}%`,
                backgroundColor: p.value === 0 ? "rgba(255,255,255,0.08)" : color,
              }}
            />
          </div>
        ))}
      </div>
      <div className="flex justify-between text-[11px] text-gray-500 mt-1.5 tabular-nums">
        <span>{points[0]?.label ?? ""}</span>
        <span className="text-gray-400">
          สูงสุด {formatNumber(max)}
          {unit} · ล่าสุด {formatNumber(points[last]?.value ?? 0)}
          {unit}
        </span>
        <span>{points[last]?.label ?? ""}</span>
      </div>
      <span className="sr-only">
        {points.map((p) => `${p.label}: ${p.value}`).join(", ")} — สูงสุดที่{" "}
        {points[peakIndex]?.label}
      </span>
    </div>
  );
}

// A compact table shell. Panels supply their own headers and rows; this only
// owns the scroll container, so a wide table scrolls itself instead of the page.
export function Table({ head, children }: { head: ReactNode; children: ReactNode }) {
  return (
    <div className="overflow-x-auto -mx-4 sm:-mx-5">
      <table className="w-full min-w-[560px] text-sm">
        <thead>
          <tr className="text-[11px] text-gray-500 uppercase tracking-wider border-b border-white/5">
            {head}
          </tr>
        </thead>
        <tbody>{children}</tbody>
      </table>
    </div>
  );
}

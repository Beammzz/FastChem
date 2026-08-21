"use client";

import { useState } from "react";
import { AdminQuestionPreview, AdminTopic } from "@/types";
import { categoryBadgeClass, categoryLabel } from "@/data/categories";
import { Empty, Panel, Spinner } from "./ui";

interface QuestionsPanelProps {
  topics: AdminTopic[];
  loading: boolean;
  onPreview: (category: string) => Promise<AdminQuestionPreview>;
}

const DIFFICULTY_LABELS: Record<string, string> = {
  easy: "ง่าย",
  medium: "ปานกลาง",
  hard: "ยาก",
};

// The preview is generated fresh on the server and thrown away — it is never
// stored and never scored, which is why it may show the answer.
export default function QuestionsPanel({ topics, loading, onPreview }: QuestionsPanelProps) {
  const [preview, setPreview] = useState<AdminQuestionPreview | null>(null);
  const [selected, setSelected] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const load = async (category: string) => {
    setSelected(category);
    setBusy(true);
    setError("");
    try {
      setPreview(await onPreview(category));
    } catch (e) {
      setPreview(null);
      setError(e instanceof Error ? e.message : "สุ่มคำถามไม่สำเร็จ");
    } finally {
      setBusy(false);
    }
  };

  const byDifficulty = (difficulty: string) => topics.filter((t) => t.difficulty === difficulty);

  if (loading) return <Spinner />;

  return (
    <div className="space-y-4">
      <Panel
        title="คลังคำถาม"
        subtitle={`${topics.length} หัวข้อที่เครื่องสุ่มคำถามรองรับ — กดเพื่อสุ่มตัวอย่างพร้อมเฉลย`}
      >
        {topics.length === 0 ? (
          <Empty icon="📚" text="ยังไม่มีหัวข้อ" />
        ) : (
          <div className="space-y-4">
            {["easy", "medium", "hard"].map((difficulty) => (
              <div key={difficulty}>
                <div className="text-xs text-gray-500 mb-2">
                  {DIFFICULTY_LABELS[difficulty]} · {byDifficulty(difficulty).length} หัวข้อ
                </div>
                <div className="flex flex-wrap gap-2">
                  {byDifficulty(difficulty).map((t) => (
                    <button
                      key={t.category}
                      onClick={() => load(t.category)}
                      className={`px-2.5 py-1 rounded-lg border text-xs transition-all cursor-pointer ${categoryBadgeClass(
                        t.category
                      )} ${selected === t.category ? "ring-1 ring-violet-400/50" : ""}`}
                      title={t.category}
                    >
                      {categoryLabel(t.category)}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </Panel>

      <Panel
        title="ตัวอย่างคำถาม"
        subtitle={selected ? selected : "เลือกหัวข้อด้านบน"}
        action={
          selected && (
            <button
              onClick={() => load(selected)}
              disabled={busy}
              className="px-3 py-1.5 rounded-lg text-xs font-medium bg-violet-500/20 border border-violet-500/30 text-violet-300 hover:bg-violet-500/30 cursor-pointer disabled:opacity-40"
            >
              สุ่มใหม่
            </button>
          )
        }
      >
        {busy ? (
          <Spinner />
        ) : error ? (
          <p className="text-sm text-red-400">{error}</p>
        ) : !preview ? (
          <Empty icon="❓" text="ยังไม่ได้เลือกหัวข้อ" />
        ) : (
          <div>
            <p className="text-base text-gray-100 mb-4">{preview.question}</p>
            <ul className="space-y-2">
              {preview.choices.map((choice, i) => (
                <li
                  key={choice}
                  className={`px-3 py-2 rounded-lg border text-sm ${
                    i === preview.correctIndex
                      ? "border-emerald-500/40 bg-emerald-500/10 text-emerald-300"
                      : "border-white/5 text-gray-400"
                  }`}
                >
                  {choice}
                  {i === preview.correctIndex && <span className="ml-2 text-xs">← เฉลย</span>}
                </li>
              ))}
            </ul>
            <div className="text-xs text-gray-500 mt-4">
              {DIFFICULTY_LABELS[preview.difficulty] ?? preview.difficulty} · จับเวลา{" "}
              {preview.timeLimit} วินาที
            </div>
          </div>
        )}
      </Panel>
    </div>
  );
}

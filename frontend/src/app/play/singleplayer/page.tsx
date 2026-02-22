"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Navbar from "@/components/Navbar";

const QUESTION_COUNT_OPTIONS = [
  { label: "5 ข้อ", value: 5 },
  { label: "10 ข้อ", value: 10 },
  { label: "15 ข้อ", value: 15 },
  { label: "20 ข้อ", value: 20 },
  { label: "30 ข้อ", value: 30 },
];

const DIFFICULTY_OPTIONS = [
  {
    label: "ง่าย",
    value: "easy",
    desc: "โครงสร้างอะตอม & ออกซิเดชัน",
    color: "violet",
  },
  {
    label: "ปานกลาง",
    value: "medium",
    desc: "โมลคอนเซ็ปต์ (มวล↔โมล↔อนุภาค)",
    color: "amber",
  },
  {
    label: "ยาก",
    value: "hard",
    desc: "การเจือจาง, เตรียมสารละลาย, จุดเยือกแข็ง",
    color: "red",
  },
  {
    label: "กำหนดเอง",
    value: "custom",
    desc: "เลือกหัวข้อเองจากตั้งค่าขั้นสูง",
    color: "cyan",
  },
];

const ALL_CATEGORIES = [
  {
    id: "atomic_structure",
    label: "โครงสร้างอะตอม",
    desc: "โปรตอน, นิวตรอน, อิเล็กตรอน",
    icon: "⚛️",
    difficulty: "easy",
  },
  {
    id: "oxidation_number",
    label: "เลขออกซิเดชัน",
    desc: "วิเคราะห์สารประกอบ",
    icon: "🔢",
    difficulty: "easy",
  },
  {
    id: "mole_concept",
    label: "โมลคอนเซ็ปต์",
    desc: "มวล↔โมล↔อนุภาค",
    icon: "⚖️",
    difficulty: "medium",
  },
  {
    id: "dilution",
    label: "การเจือจาง",
    desc: "M₁V₁ = M₂V₂",
    icon: "🧪",
    difficulty: "hard",
  },
  {
    id: "preparing_solution",
    label: "เตรียมสารละลาย",
    desc: "mass = M × V × MW",
    icon: "🧂",
    difficulty: "hard",
  },
  {
    id: "freezing_point",
    label: "จุดเยือกแข็ง",
    desc: "ΔTf = iKfm",
    icon: "❄️",
    difficulty: "hard",
  },
];

export default function SinglePlayerPage() {
  const [selectedQuestions, setSelectedQuestions] = useState(10);
  const [selectedDifficulty, setSelectedDifficulty] = useState("easy");
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const router = useRouter();

  const isCustom = selectedDifficulty === "custom";

  const toggleCategory = (id: string) => {
    setSelectedCategories((prev) =>
      prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id]
    );
  };

  const selectAll = () => {
    setSelectedCategories(ALL_CATEGORIES.map((c) => c.id));
  };

  const deselectAll = () => {
    setSelectedCategories([]);
  };

  const startGame = () => {
    const params = new URLSearchParams();
    params.set("questions", String(selectedQuestions));

    if (isCustom && selectedCategories.length > 0) {
      params.set("difficulty", "custom");
      params.set("categories", selectedCategories.join(","));
    } else {
      params.set("difficulty", selectedDifficulty === "custom" ? "easy" : selectedDifficulty);
    }

    router.push(`/game?${params.toString()}`);
  };

  const canStart = !isCustom || selectedCategories.length > 0;

  return (
    <div className="min-h-screen bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-2xl mx-auto px-4 sm:px-6 pt-6 sm:pt-8 pb-16">
        {/* Back link */}
        <Link
          href="/play"
          className="inline-flex items-center text-sm text-gray-400 hover:text-white transition-colors mb-6"
        >
          ← กลับเลือกโหมด
        </Link>

        {/* Header */}
        <div className="flex items-center gap-3 mb-8">
          <div className="w-12 h-12 rounded-xl bg-violet-500/10 flex items-center justify-center text-2xl">
            🧪
          </div>
          <div>
            <h1 className="text-2xl font-bold">เล่นคนเดียว</h1>
            <p className="text-sm text-gray-400">ตั้งค่าเกมของคุณ</p>
          </div>
        </div>

        <div className="space-y-6">
          {/* Question Count Selection */}
          <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
            <p className="text-xs text-gray-500 uppercase tracking-wider font-medium mb-4">
              📝 จำนวนคำถาม
            </p>
            <div className="flex flex-wrap gap-2">
              {QUESTION_COUNT_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  onClick={() => setSelectedQuestions(opt.value)}
                  className={`flex-1 min-w-[60px] py-2 sm:py-2.5 rounded-lg text-xs sm:text-sm font-bold transition-all ${
                    selectedQuestions === opt.value
                      ? "bg-violet-500 text-white shadow-md shadow-violet-500/20"
                      : "bg-white/5 text-gray-400 hover:bg-white/10 hover:text-gray-200"
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>

          {/* Difficulty Selection */}
          <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6">
            <p className="text-xs text-gray-500 uppercase tracking-wider font-medium mb-4">
              📊 ระดับความยาก
            </p>
            <div className="grid grid-cols-2 gap-3">
              {DIFFICULTY_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  onClick={() => {
                    setSelectedDifficulty(opt.value);
                    if (opt.value === "custom") {
                      setShowAdvanced(true);
                    }
                  }}
                  className={`text-left p-4 rounded-xl border transition-all ${
                    selectedDifficulty === opt.value
                      ? opt.value === "easy"
                        ? "bg-violet-500/10 border-violet-500/40 ring-1 ring-violet-500/20"
                        : opt.value === "medium"
                        ? "bg-amber-500/10 border-amber-500/40 ring-1 ring-amber-500/20"
                        : opt.value === "hard"
                        ? "bg-red-500/10 border-red-500/40 ring-1 ring-red-500/20"
                        : "bg-cyan-500/10 border-cyan-500/40 ring-1 ring-cyan-500/20"
                      : "bg-white/[0.02] border-white/5 hover:border-white/10"
                  }`}
                >
                  <div className="flex items-center gap-2 mb-1">
                    <span
                      className={`text-sm font-bold ${
                        selectedDifficulty === opt.value
                          ? opt.value === "easy"
                            ? "text-violet-400"
                            : opt.value === "medium"
                            ? "text-amber-400"
                            : opt.value === "hard"
                            ? "text-red-400"
                            : "text-cyan-400"
                          : "text-white"
                      }`}
                    >
                      {opt.label}
                    </span>
                  </div>
                  <p className="text-xs text-gray-500">{opt.desc}</p>
                </button>
              ))}
            </div>
          </div>

          {/* Advanced Settings (Category Picker) */}
          {!isCustom && (
            <button
              onClick={() => {
                setShowAdvanced(!showAdvanced);
                if (!showAdvanced) {
                  setSelectedDifficulty("custom");
                }
              }}
              className="w-full text-left bg-[#12122a] border border-white/5 rounded-2xl p-4 hover:border-white/10 transition-all flex items-center justify-between"
            >
              <div className="flex items-center gap-2">
                <span className="text-sm">⚙️</span>
                <span className="text-sm text-gray-400">ตั้งค่าขั้นสูง</span>
              </div>
              <span className="text-xs text-gray-500">เลือกหัวข้อเอง →</span>
            </button>
          )}

          {isCustom && (
            <div className="bg-[#12122a] border border-white/5 rounded-2xl p-6 animate-fade-in">
              <div className="flex items-center justify-between mb-4">
                <p className="text-xs text-gray-500 uppercase tracking-wider font-medium">
                  ⚙️ เลือกหัวข้อ
                </p>
                <div className="flex gap-2">
                  <button
                    onClick={selectAll}
                    className="text-xs text-violet-400 hover:text-violet-300 transition-colors"
                  >
                    เลือกทั้งหมด
                  </button>
                  <span className="text-gray-600">|</span>
                  <button
                    onClick={deselectAll}
                    className="text-xs text-gray-400 hover:text-gray-300 transition-colors"
                  >
                    ล้าง
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-1 gap-2">
                {ALL_CATEGORIES.map((cat) => {
                  const isSelected = selectedCategories.includes(cat.id);
                  const diffColor =
                    cat.difficulty === "easy"
                      ? "text-violet-400"
                      : cat.difficulty === "medium"
                      ? "text-amber-400"
                      : "text-red-400";
                  const diffLabel =
                    cat.difficulty === "easy"
                      ? "ง่าย"
                      : cat.difficulty === "medium"
                      ? "ปานกลาง"
                      : "ยาก";
                  return (
                    <button
                      key={cat.id}
                      onClick={() => toggleCategory(cat.id)}
                      className={`flex items-center gap-3 p-3 rounded-xl border transition-all text-left ${
                        isSelected
                          ? "bg-violet-500/10 border-violet-500/30"
                          : "bg-white/[0.02] border-white/5 hover:border-white/10"
                      }`}
                    >
                      <span
                        className={`w-5 h-5 rounded flex items-center justify-center text-xs border transition-all ${
                          isSelected
                            ? "bg-violet-500 border-violet-500 text-white"
                            : "border-white/20 text-transparent"
                        }`}
                      >
                        ✓
                      </span>
                      <span className="text-lg">{cat.icon}</span>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span
                            className={`text-sm font-medium ${
                              isSelected ? "text-white" : "text-gray-300"
                            }`}
                          >
                            {cat.label}
                          </span>
                          <span
                            className={`text-[10px] px-1.5 py-0.5 rounded-full ${diffColor} bg-white/5`}
                          >
                            {diffLabel}
                          </span>
                        </div>
                        <p className="text-xs text-gray-500 truncate">
                          {cat.desc}
                        </p>
                      </div>
                    </button>
                  );
                })}
              </div>

              {selectedCategories.length === 0 && (
                <p className="text-xs text-red-400/70 mt-3 text-center">
                  กรุณาเลือกอย่างน้อย 1 หัวข้อ
                </p>
              )}
            </div>
          )}

          {/* Start Button */}
          <button
            onClick={startGame}
            disabled={!canStart}
            className="w-full bg-violet-500 hover:bg-violet-600 disabled:opacity-40 disabled:cursor-not-allowed text-white font-bold text-lg py-4 rounded-xl transition-all hover:scale-[1.01] active:scale-[0.99] shadow-lg shadow-violet-500/15"
          >
            เริ่มเกม
          </button>

          {/* Summary */}
          <div className="text-center text-xs text-gray-500 space-x-3">
            <span>📝 {selectedQuestions} ข้อ</span>
            <span>·</span>
            <span>
              📊{" "}
              {isCustom
                ? `${selectedCategories.length} หัวข้อ`
                : DIFFICULTY_OPTIONS.find((d) => d.value === selectedDifficulty)
                    ?.label}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

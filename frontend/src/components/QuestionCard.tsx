"use client";

import { Question } from "@/types";

interface QuestionCardProps {
  question: Question;
  selectedIndex: number | null;
  isCorrect: boolean | null;
  scoreEarned: number;
  speedBonus: number;
  comboMultiplier?: number;
  onAnswer: (index: number) => void;
}

export default function QuestionCard({
  question,
  selectedIndex,
  isCorrect,
  scoreEarned,
  speedBonus,
  comboMultiplier = 1,
  onAnswer,
}: QuestionCardProps) {
  const getCategoryLabel = (cat: string) => {
    switch (cat) {
      case "atomic_structure":
        return "โครงสร้างอะตอม";
      case "oxidation_number":
        return "เลขออกซิเดชัน";
      case "state_of_matter":
        return "สถานะของสาร";
      case "mole_concept":
        return "โมลคอนเซ็ปต์";
      case "dilution":
        return "การเจือจาง";
      case "preparing_solution":
        return "เตรียมสารละลาย";
      case "freezing_point":
        return "จุดเยือกแข็ง";
      default:
        return cat;
    }
  };

  const getCategoryColor = (cat: string) => {
    switch (cat) {
      case "atomic_structure":
        return "bg-cyan-500/10 text-cyan-400 border-cyan-500/20";
      case "oxidation_number":
        return "bg-purple-500/10 text-purple-400 border-purple-500/20";
      case "state_of_matter":
        return "bg-teal-500/10 text-teal-400 border-teal-500/20";
      case "mole_concept":
        return "bg-amber-500/10 text-amber-400 border-amber-500/20";
      case "dilution":
        return "bg-rose-500/10 text-rose-400 border-rose-500/20";
      case "preparing_solution":
        return "bg-orange-500/10 text-orange-400 border-orange-500/20";
      case "freezing_point":
        return "bg-sky-500/10 text-sky-400 border-sky-500/20";
      default:
        return "bg-gray-500/10 text-gray-400 border-gray-500/20";
    }
  };

  const getChoiceStyle = (index: number) => {
    if (selectedIndex === null) {
      return "bg-[#12122a] border-white/5 hover:border-violet-500/40 hover:bg-violet-500/5 text-white cursor-pointer";
    }
    if (index === question.correctIndex) {
      return "bg-violet-500/10 border-violet-500/50 text-violet-300";
    }
    if (index === selectedIndex && !isCorrect) {
      return "bg-red-500/10 border-red-500/50 text-red-300";
    }
    return "bg-[#12122a] border-white/5 text-gray-600";
  };

  const getChoicePrefix = (index: number) => {
    if (selectedIndex === null) return "";
    if (index === question.correctIndex) return "✓ ";
    if (index === selectedIndex && !isCorrect) return "✗ ";
    return "";
  };

  return (
    <div className="w-full animate-fade-in">
      {/* Category badge */}
      <div className="flex justify-center mb-6">
        <span
          className={`inline-block px-4 py-1.5 rounded-full text-xs font-semibold border tracking-wide uppercase ${getCategoryColor(
            question.category
          )}`}
        >
          {getCategoryLabel(question.category)}
        </span>
      </div>

      {/* Question text */}
      <div className="bg-[#12122a] border border-white/5 rounded-2xl p-4 sm:p-8 mb-4 sm:mb-6">
        <p className="text-base sm:text-xl md:text-2xl font-semibold text-white text-center leading-relaxed">
          {question.question}
        </p>
      </div>

      {/* Answer choices */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 sm:gap-3">
        {question.choices.map((choice, index) => (
          <button
            key={index}
            onClick={() => onAnswer(index)}
            disabled={selectedIndex !== null}
            className={`
              border rounded-xl px-4 py-3 sm:px-6 sm:py-5 text-sm sm:text-lg font-bold
              transition-all duration-150 min-h-[44px]
              ${selectedIndex === null ? "hover:scale-[1.02] active:scale-[0.98]" : ""}
              ${getChoiceStyle(index)}
            `}
          >
            <span className="text-xs sm:text-sm opacity-40 mr-2">
              {String.fromCharCode(65 + index)}
            </span>
            {getChoicePrefix(index)}
            {choice}
          </button>
        ))}
      </div>

      {/* Feedback */}
      {selectedIndex !== null && (
        <div
          className={`mt-5 text-center animate-fade-in ${
            isCorrect ? "text-violet-400" : "text-red-400"
          }`}
        >
          {selectedIndex === -1 ? (
            <div className="text-base font-bold">⏰ หมดเวลา!</div>
          ) : isCorrect ? (
            <div>
              <div className="text-base font-bold">
                ถูกต้อง! +{scoreEarned} คะแนน
              </div>
              <div className="flex items-center justify-center gap-3 mt-1">
                {speedBonus > 0 && (
                  <span className="text-xs text-violet-300">
                    ⚡ Speed +{speedBonus}
                  </span>
                )}
                {comboMultiplier > 1 && (
                  <span className="text-xs text-amber-400 font-bold">
                    🔥 x{comboMultiplier} Combo!
                  </span>
                )}
              </div>
            </div>
          ) : (
            <div className="text-base font-bold">ผิด</div>
          )}
        </div>
      )}
    </div>
  );
}

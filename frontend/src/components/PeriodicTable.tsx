"use client";

import { useState } from "react";
import {
  ELEMENTS,
  CATEGORY_COLORS,
  CATEGORY_LABELS,
  Element,
  ElementCategory,
} from "@/data/periodicTable";

interface PeriodicTableProps {
  /** When true, the toggle button is visible */
  enabled?: boolean;
}

export default function PeriodicTable({ enabled = true }: PeriodicTableProps) {
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = useState<Element | null>(null);

  if (!enabled) return null;

  // Build a grid: rows = periods 1-7 + lanthanide (8) + actinide (9)
  // Total 9 display rows, 18 columns each
  const grid: (Element | null)[][] = Array.from({ length: 9 }, () =>
    Array(18).fill(null)
  );
  for (const el of ELEMENTS) {
    const rowIdx = el.period - 1; // period 1→row 0, ... period 9→row 8
    grid[rowIdx][el.group - 1] = el;
  }

  return (
    <>
      {/* Toggle button */}
      <button
        onClick={() => setOpen(!open)}
        className={`fixed bottom-24 right-3 sm:bottom-4 sm:right-4 z-50 w-11 h-11 sm:w-12 sm:h-12 rounded-full flex items-center justify-center text-base sm:text-lg font-bold shadow-lg transition-all
          ${
            open
              ? "bg-violet-600 text-white shadow-violet-500/30"
              : "bg-[#1a1a2e] border border-white/10 text-violet-400 hover:border-violet-500/40 hover:bg-violet-500/10"
          }`}
        title="ตารางธาตุ"
      >
        🧪
      </button>

      {/* Overlay */}
      {open && (
        <div className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm flex items-center justify-center p-2 sm:p-4">
          <div className="bg-[#0e0e20] border border-white/10 rounded-2xl w-full max-w-4xl max-h-[90vh] sm:max-h-[85vh] overflow-auto shadow-2xl">
            {/* Header */}
            <div className="sticky top-0 bg-[#0e0e20] border-b border-white/5 px-4 py-2.5 sm:px-5 sm:py-3 flex items-center justify-between z-10">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                🧪 ตารางธาตุ
              </h3>
              <button
                onClick={() => {
                  setOpen(false);
                  setSelected(null);
                }}
                className="text-gray-500 hover:text-white text-xl leading-none p-1"
              >
                ✕
              </button>
            </div>

            {/* Selected element detail */}
            {selected && (
              <div className="px-4 py-2.5 sm:px-5 sm:py-3 border-b border-white/5">
                <div className="flex items-center gap-3 sm:gap-4">
                  <div
                    className={`w-12 h-12 sm:w-16 sm:h-16 rounded-xl flex flex-col items-center justify-center border ${
                      CATEGORY_COLORS[selected.category]
                    }`}
                  >
                    <span className="text-[8px] sm:text-[10px] text-white/70">
                      {selected.atomicNumber}
                    </span>
                    <span className="text-base sm:text-xl font-bold text-white">
                      {selected.symbol}
                    </span>
                    <span className="text-[7px] sm:text-[9px] text-white/60">
                      {selected.atomicMass}
                    </span>
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="text-white font-bold text-sm sm:text-lg truncate">
                      {selected.name}{" "}
                      <span className="text-gray-400 font-normal text-xs sm:text-sm">
                        ({selected.nameTH})
                      </span>
                    </div>
                    <div className="text-[10px] sm:text-xs text-gray-500 mt-0.5">
                      เลขอะตอม {selected.atomicNumber} · มวลอะตอม{" "}
                      {selected.atomicMass} · กลุ่ม {selected.group} · คาบ{" "}
                      {selected.period}
                    </div>
                    <div className="mt-1">
                      <span
                        className={`inline-block text-[10px] px-2 py-0.5 rounded-full border ${
                          CATEGORY_COLORS[selected.category]
                        } text-white`}
                      >
                        {CATEGORY_LABELS[selected.category]}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Grid */}
            <div className="p-2 sm:p-4">
              <div className="w-full">
                {/* Group numbers header */}
                <div className="grid grid-cols-18 gap-[1px] sm:gap-[2px] mb-1">
                  {Array.from({ length: 18 }, (_, i) => (
                    <div
                      key={i}
                      className="text-[7px] sm:text-[9px] text-gray-600 text-center"
                    >
                      {i + 1}
                    </div>
                  ))}
                </div>

                {/* Element rows (periods 1-7) */}
                {grid.slice(0, 7).map((row, rowIdx) => (
                  <div key={rowIdx} className="grid grid-cols-18 gap-[1px] sm:gap-[2px] mb-[1px] sm:mb-[2px]">
                    {row.map((el, colIdx) => {
                      if (!el) {
                        // In period 6 (row 5) group 3, show lanthanide placeholder
                        if (rowIdx === 5 && colIdx === 2) {
                          return (
                            <div
                              key={colIdx}
                              className="aspect-square rounded-sm border flex items-center justify-center bg-pink-600/30 border-pink-500/20"
                              title="Lanthanides 57-71"
                            >
                              <span className="text-[5px] sm:text-[7px] text-pink-300 font-bold">La-Lu</span>
                            </div>
                          );
                        }
                        // In period 7 (row 6) group 3, show actinide placeholder
                        if (rowIdx === 6 && colIdx === 2) {
                          return (
                            <div
                              key={colIdx}
                              className="aspect-square rounded-sm border flex items-center justify-center bg-rose-600/30 border-rose-500/20"
                              title="Actinides 89-103"
                            >
                              <span className="text-[5px] sm:text-[7px] text-rose-300 font-bold">Ac-Lr</span>
                            </div>
                          );
                        }
                        return <div key={colIdx} className="aspect-square" />;
                      }
                      const isSelected = selected?.atomicNumber === el.atomicNumber;
                      return (
                        <button
                          key={colIdx}
                          onClick={() =>
                            setSelected(isSelected ? null : el)
                          }
                          className={`aspect-square rounded-sm border flex flex-col items-center justify-center transition-all text-white cursor-pointer active:scale-95 sm:hover:scale-110 sm:hover:z-10 ${
                            CATEGORY_COLORS[el.category]
                          } ${
                            isSelected
                              ? "ring-2 ring-white scale-110 z-10"
                              : ""
                          }`}
                          title={`${el.name} (${el.nameTH})`}
                        >
                          <span className="text-[5px] sm:text-[7px] leading-none opacity-70">
                            {el.atomicNumber}
                          </span>
                          <span className="text-[8px] sm:text-[11px] font-bold leading-tight">
                            {el.symbol}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                ))}

                {/* Separator */}
                <div className="h-2 sm:h-3" />

                {/* Lanthanide row (period 8) */}
                <div className="grid grid-cols-18 gap-[1px] sm:gap-[2px] mb-[1px] sm:mb-[2px]">
                  {grid[7].map((el, colIdx) => {
                    // Lanthanides occupy groups 3-17 (cols 2-16)
                    if (!el) {
                      return <div key={colIdx} className="aspect-square" />;
                    }
                    const isSelected = selected?.atomicNumber === el.atomicNumber;
                    return (
                      <button
                        key={colIdx}
                        onClick={() => setSelected(isSelected ? null : el)}
                        className={`aspect-square rounded-sm border flex flex-col items-center justify-center transition-all text-white cursor-pointer active:scale-95 sm:hover:scale-110 sm:hover:z-10 ${
                          CATEGORY_COLORS[el.category]
                        } ${isSelected ? "ring-2 ring-white scale-110 z-10" : ""}`}
                        title={`${el.name} (${el.nameTH})`}
                      >
                        <span className="text-[5px] sm:text-[7px] leading-none opacity-70">{el.atomicNumber}</span>
                        <span className="text-[8px] sm:text-[11px] font-bold leading-tight">{el.symbol}</span>
                      </button>
                    );
                  })}
                </div>

                {/* Actinide row (period 9) */}
                <div className="grid grid-cols-18 gap-[1px] sm:gap-[2px] mb-[1px] sm:mb-[2px]">
                  {grid[8].map((el, colIdx) => {
                    if (!el) {
                      return <div key={colIdx} className="aspect-square" />;
                    }
                    const isSelected = selected?.atomicNumber === el.atomicNumber;
                    return (
                      <button
                        key={colIdx}
                        onClick={() => setSelected(isSelected ? null : el)}
                        className={`aspect-square rounded-sm border flex flex-col items-center justify-center transition-all text-white cursor-pointer active:scale-95 sm:hover:scale-110 sm:hover:z-10 ${
                          CATEGORY_COLORS[el.category]
                        } ${isSelected ? "ring-2 ring-white scale-110 z-10" : ""}`}
                        title={`${el.name} (${el.nameTH})`}
                      >
                        <span className="text-[5px] sm:text-[7px] leading-none opacity-70">{el.atomicNumber}</span>
                        <span className="text-[8px] sm:text-[11px] font-bold leading-tight">{el.symbol}</span>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Legend */}
              <div className="flex flex-wrap gap-1.5 sm:gap-2 mt-3 sm:mt-4 justify-center">
                {(Object.keys(CATEGORY_LABELS) as ElementCategory[]).map(
                  (cat) => (
                    <div
                      key={cat}
                      className="flex items-center gap-1 sm:gap-1.5 text-[8px] sm:text-[10px] text-gray-400"
                    >
                      <div
                        className={`w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-sm border ${CATEGORY_COLORS[cat]}`}
                      />
                      {CATEGORY_LABELS[cat]}
                    </div>
                  )
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

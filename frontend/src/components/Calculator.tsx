"use client";

import { useState, useCallback } from "react";

export default function Calculator() {
  const [open, setOpen] = useState(false);
  const [display, setDisplay] = useState("0");
  const [expression, setExpression] = useState("");
  const [newNumber, setNewNumber] = useState(true);

  const constants: { label: string; value: string; desc: string }[] = [
    { label: "Nₐ", value: "6.022e23", desc: "เลขอาโวกาโดร" },
    { label: "R", value: "8.314", desc: "ค่าคงที่แก๊ส (J/mol·K)" },
    { label: "F", value: "96485", desc: "ค่าคงที่ฟาราเดย์ (C/mol)" },
    { label: "Kf", value: "1.86", desc: "H₂O cryoscopic (°C·kg/mol)" },
  ];

  const handleNumber = useCallback(
    (num: string) => {
      if (newNumber) {
        setDisplay(num === "." ? "0." : num);
        setNewNumber(false);
      } else {
        if (num === "." && display.includes(".")) return;
        setDisplay((prev) => (prev === "0" && num !== "." ? num : prev + num));
      }
    },
    [newNumber, display]
  );

  const handleOperator = useCallback(
    (op: string) => {
      setExpression((prev) => prev + display + " " + op + " ");
      setNewNumber(true);
    },
    [display]
  );

  const handleEquals = useCallback(() => {
    try {
      const fullExpr = expression + display;
      // Safe eval using Function constructor with only math operations
      const sanitized = fullExpr.replace(/[^0-9+\-*/().e ]/g, "");
      if (!sanitized) return;
      const result = new Function("return " + sanitized)();
      const formatted =
        typeof result === "number"
          ? Number.isInteger(result) && Math.abs(result) < 1e15
            ? result.toString()
            : result.toPrecision(10).replace(/\.?0+$/, "")
          : "Error";
      setDisplay(formatted);
      setExpression("");
      setNewNumber(true);
    } catch {
      setDisplay("Error");
      setExpression("");
      setNewNumber(true);
    }
  }, [expression, display]);

  const handleClear = useCallback(() => {
    setDisplay("0");
    setExpression("");
    setNewNumber(true);
  }, []);

  const handleBackspace = useCallback(() => {
    if (display.length > 1) {
      setDisplay((prev) => prev.slice(0, -1));
    } else {
      setDisplay("0");
      setNewNumber(true);
    }
  }, [display]);

  const handleConstant = useCallback((value: string) => {
    setDisplay(value);
    setNewNumber(false);
  }, []);

  const handleParenthesis = useCallback(
    (paren: string) => {
      if (paren === "(") {
        setExpression((prev) => prev + "(");
        setNewNumber(true);
      } else {
        setExpression((prev) => prev + display + ")");
        setNewNumber(true);
      }
    },
    [display]
  );

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="fixed bottom-20 left-3 sm:bottom-4 sm:left-4 z-50 w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center text-base sm:text-lg font-bold shadow-lg bg-[#1a1a2e] border border-white/10 text-blue-400 hover:border-blue-500/40 hover:bg-blue-500/10 transition-all"
        title="เครื่องคิดเลข"
      >
        🔢
      </button>
    );
  }

  return (
    <div className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm flex items-end sm:items-center justify-center p-2 sm:p-4">
      <div className="bg-[#0e0e20] border border-white/10 rounded-2xl w-full max-w-sm shadow-2xl">
        {/* Header */}
        <div className="border-b border-white/5 px-5 py-3 flex items-center justify-between">
          <h3 className="text-sm font-bold text-white flex items-center gap-2">
            🔢 เครื่องคิดเลข
          </h3>
          <button
            onClick={() => setOpen(false)}
            className="text-gray-500 hover:text-white text-xl leading-none"
          >
            ✕
          </button>
        </div>

        {/* Display */}
        <div className="px-4 pt-3">
          <div className="bg-[#1a1a2e] rounded-xl p-4 border border-white/5">
            <div className="text-xs text-gray-500 h-5 text-right font-mono truncate">
              {expression}
            </div>
            <div className="text-2xl font-bold text-white text-right font-mono truncate">
              {display}
            </div>
          </div>
        </div>

        {/* Constants */}
        <div className="px-4 pt-3">
          <div className="flex gap-2">
            {constants.map((c) => (
              <button
                key={c.label}
                onClick={() => handleConstant(c.value)}
                className="flex-1 bg-violet-500/10 border border-violet-500/20 text-violet-300 rounded-lg py-1.5 text-xs font-bold hover:bg-violet-500/20 transition-colors"
                title={c.desc}
              >
                {c.label}
              </button>
            ))}
          </div>
        </div>

        {/* Buttons */}
        <div className="p-4 grid grid-cols-4 gap-2">
          {/* Row 1 */}
          <button
            onClick={handleClear}
            className="bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl py-3 font-bold text-sm hover:bg-red-500/20 transition-colors"
          >
            AC
          </button>
          <button
            onClick={() => handleParenthesis("(")}
            className="bg-white/5 border border-white/10 text-gray-300 rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
          >
            (
          </button>
          <button
            onClick={() => handleParenthesis(")")}
            className="bg-white/5 border border-white/10 text-gray-300 rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
          >
            )
          </button>
          <button
            onClick={() => handleOperator("/")}
            className="bg-blue-500/10 border border-blue-500/20 text-blue-400 rounded-xl py-3 font-bold text-sm hover:bg-blue-500/20 transition-colors"
          >
            ÷
          </button>

          {/* Row 2 */}
          {["7", "8", "9"].map((n) => (
            <button
              key={n}
              onClick={() => handleNumber(n)}
              className="bg-white/5 border border-white/10 text-white rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
            >
              {n}
            </button>
          ))}
          <button
            onClick={() => handleOperator("*")}
            className="bg-blue-500/10 border border-blue-500/20 text-blue-400 rounded-xl py-3 font-bold text-sm hover:bg-blue-500/20 transition-colors"
          >
            ×
          </button>

          {/* Row 3 */}
          {["4", "5", "6"].map((n) => (
            <button
              key={n}
              onClick={() => handleNumber(n)}
              className="bg-white/5 border border-white/10 text-white rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
            >
              {n}
            </button>
          ))}
          <button
            onClick={() => handleOperator("-")}
            className="bg-blue-500/10 border border-blue-500/20 text-blue-400 rounded-xl py-3 font-bold text-sm hover:bg-blue-500/20 transition-colors"
          >
            −
          </button>

          {/* Row 4 */}
          {["1", "2", "3"].map((n) => (
            <button
              key={n}
              onClick={() => handleNumber(n)}
              className="bg-white/5 border border-white/10 text-white rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
            >
              {n}
            </button>
          ))}
          <button
            onClick={() => handleOperator("+")}
            className="bg-blue-500/10 border border-blue-500/20 text-blue-400 rounded-xl py-3 font-bold text-sm hover:bg-blue-500/20 transition-colors"
          >
            +
          </button>

          {/* Row 5 */}
          <button
            onClick={() => handleNumber("0")}
            className="bg-white/5 border border-white/10 text-white rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
          >
            0
          </button>
          <button
            onClick={() => handleNumber(".")}
            className="bg-white/5 border border-white/10 text-white rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
          >
            .
          </button>
          <button
            onClick={handleBackspace}
            className="bg-white/5 border border-white/10 text-gray-400 rounded-xl py-3 font-bold text-sm hover:bg-white/10 transition-colors"
          >
            ⌫
          </button>
          <button
            onClick={handleEquals}
            className="bg-violet-500 text-white rounded-xl py-3 font-bold text-sm hover:bg-violet-600 transition-colors"
          >
            =
          </button>
        </div>
      </div>
    </div>
  );
}

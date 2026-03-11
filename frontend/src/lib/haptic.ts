/**
 * haptic.ts — Thin wrapper around the Web Vibration API.
 * Silently no-ops on browsers / devices that don't support it.
 */

type HapticType = "correct" | "wrong" | "timeout";

const patterns: Record<HapticType, number | number[]> = {
  /** Short double-pulse — satisfying "ding ding" feel */
  correct: [40, 30, 80],
  /** Single firm buzz */
  wrong: [120],
  /** Two short pulses indicating time ran out */
  timeout: [60, 40, 60],
};

export function triggerHaptic(type: HapticType): void {
  if (typeof navigator === "undefined" || !navigator.vibrate) return;
  try {
    navigator.vibrate(patterns[type]);
  } catch {
    // Vibration API can throw in some sandboxed environments
  }
}

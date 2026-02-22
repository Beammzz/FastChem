"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/components/AuthProvider";
import Navbar from "@/components/Navbar";

export default function LoginPage() {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { login, register } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      if (mode === "login") {
        await login(username, password);
      } else {
        await register(username, password);
      }
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "เกิดข้อผิดพลาด");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="flex items-center justify-center px-6 pt-20 pb-16">
        <div className="w-full max-w-sm">
          {/* Mode tabs */}
          <div className="flex bg-white/5 rounded-xl p-1 mb-8">
            <button
              onClick={() => {
                setMode("login");
                setError("");
              }}
              className={`flex-1 py-2.5 rounded-lg text-sm font-bold transition-all ${
                mode === "login"
                  ? "bg-violet-500 text-white shadow-md"
                  : "text-gray-400 hover:text-white"
              }`}
            >
              เข้าสู่ระบบ
            </button>
            <button
              onClick={() => {
                setMode("register");
                setError("");
              }}
              className={`flex-1 py-2.5 rounded-lg text-sm font-bold transition-all ${
                mode === "register"
                  ? "bg-violet-500 text-white shadow-md"
                  : "text-gray-400 hover:text-white"
              }`}
            >
              สมัครสมาชิก
            </button>
          </div>

          <div className="bg-[#12122a] border border-white/5 rounded-2xl p-8">
            <div className="text-center mb-6">
              <div className="text-4xl mb-2">⚗️</div>
              <h1 className="text-2xl font-bold">
                {mode === "login" ? "ยินดีต้อนรับกลับ" : "สร้างบัญชี"}
              </h1>
              <p className="text-gray-400 text-sm mt-1">
                {mode === "login"
                  ? "เข้าสู่ระบบเพื่อติดตามคะแนน"
                  : "สมัครเพื่อบันทึกคะแนน & แข่งขัน"}
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs text-gray-500 uppercase tracking-wider font-medium mb-2">
                  ชื่อผู้ใช้
                </label>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:border-violet-500/50 focus:ring-1 focus:ring-violet-500/20 transition-all"
                  placeholder="กรอกชื่อผู้ใช้"
                  minLength={3}
                  maxLength={20}
                  required
                  autoComplete="username"
                />
              </div>

              <div>
                <label className="block text-xs text-gray-500 uppercase tracking-wider font-medium mb-2">
                  รหัสผ่าน
                </label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:border-violet-500/50 focus:ring-1 focus:ring-violet-500/20 transition-all"
                  placeholder="กรอกรหัสผ่าน"
                  minLength={6}
                  required
                  autoComplete={
                    mode === "login" ? "current-password" : "new-password"
                  }
                />
              </div>

              {error && (
                <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-sm rounded-xl px-4 py-3">
                  {error}
                </div>
              )}

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-violet-500 hover:bg-violet-600 disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold py-3.5 rounded-xl transition-all hover:scale-[1.01] active:scale-[0.99] shadow-lg shadow-violet-500/15"
              >
                {loading
                  ? "กำลังโหลด..."
                  : mode === "login"
                  ? "เข้าสู่ระบบ"
                  : "สร้างบัญชี"}
              </button>
            </form>

            <p className="text-center text-gray-500 text-sm mt-6">
              {mode === "login" ? (
                <>
                  ยังไม่มีบัญชี?{" "}
                  <button
                    onClick={() => {
                      setMode("register");
                      setError("");
                    }}
                    className="text-violet-400 hover:underline"
                  >
                    สมัครสมาชิก
                  </button>
                </>
              ) : (
                <>
                  มีบัญชีอยู่แล้ว?{" "}
                  <button
                    onClick={() => {
                      setMode("login");
                      setError("");
                    }}
                    className="text-violet-400 hover:underline"
                  >
                    เข้าสู่ระบบ
                  </button>
                </>
              )}
            </p>
          </div>

          <p className="text-center text-gray-600 text-xs mt-6">
            <Link href="/" className="hover:text-gray-400 transition-colors">
              ← เล่นต่อโดยไม่ต้องเข้าสู่ระบบ
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}

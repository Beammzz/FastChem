"use client";

import Link from "next/link";
import { useAuth } from "@/components/AuthProvider";

export default function Navbar() {
  const { user, loading, logout } = useAuth();

  return (
    <nav className="border-b border-white/5 bg-[#0a0a1a]/80 backdrop-blur-md sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2 text-xl font-bold">
          <span className="text-2xl">⚗️</span>
          <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
            FastChem
          </span>
        </Link>
        <div className="flex items-center gap-6 text-sm text-gray-400">
          <Link href="/" className="hover:text-white transition-colors">
            หน้าหลัก
          </Link>
          <Link
            href="/leaderboard"
            className="hover:text-white transition-colors"
          >
            กระดานผู้นำ
          </Link>
          {loading ? (
            <div className="w-16 h-8 bg-white/5 rounded-lg animate-pulse" />
          ) : user ? (
            <div className="flex items-center gap-4">
              <Link
                href={`/profile/${user.username}`}
                className="text-violet-400 font-medium hover:text-violet-300 transition-colors"
              >
                {user.username}
              </Link>
              <button
                onClick={logout}
                className="text-gray-500 hover:text-red-400 transition-colors"
              >
                ออกจากระบบ
              </button>
            </div>
          ) : (
            <Link
              href="/login"
              className="bg-violet-500/10 border border-violet-500/20 text-violet-400 hover:bg-violet-500/20 px-4 py-1.5 rounded-lg font-medium transition-all"
            >
              เข้าสู่ระบบ
            </Link>
          )}
        </div>
      </div>
    </nav>
  );
}

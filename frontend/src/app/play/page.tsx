"use client";

import { useRouter } from "next/navigation";
import Navbar from "@/components/Navbar";

export default function PlayPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-[#0a0a1a] text-white">
      <Navbar />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 pt-8 sm:pt-12 pb-16">
        {/* Header */}
        <div className="text-center mb-8 sm:mb-10">
          <h1 className="text-2xl sm:text-4xl font-extrabold mb-2">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-purple-400">
              เลือกโหมดเกม
            </span>
          </h1>
          <p className="text-gray-400">เลือกรูปแบบการเล่นที่ต้องการ</p>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {/* Single Player */}
          <button
            onClick={() => router.push("/play/singleplayer")}
            className="text-left bg-[#12122a] border border-white/5 rounded-2xl p-5 sm:p-8 hover:border-violet-500/40 transition-all group cursor-pointer"
          >
            <div className="flex items-center gap-3 mb-4">
              <div className="w-12 h-12 rounded-xl bg-violet-500/10 flex items-center justify-center text-2xl group-hover:scale-110 transition-transform">
                🧪
              </div>
              <div>
                <h2 className="text-xl font-bold text-white">เล่นคนเดียว</h2>
                <p className="text-xs text-gray-500">Single Player</p>
              </div>
            </div>
            <p className="text-gray-400 text-sm leading-relaxed mb-6">
              ฝึกเคมีด้วยคำถามที่สร้างอัตโนมัติ เลือกระดับความยาก ตั้งเวลาเอง
              และปรับแต่งหัวข้อได้ตามต้องการ
            </p>
            <div className="flex flex-wrap gap-2 mb-6">
              <span className="text-xs bg-violet-500/10 text-violet-400 px-2.5 py-1 rounded-full">
                3 ระดับความยาก
              </span>
              <span className="text-xs bg-violet-500/10 text-violet-400 px-2.5 py-1 rounded-full">
                เลือกหัวข้อได้
              </span>
              <span className="text-xs bg-violet-500/10 text-violet-400 px-2.5 py-1 rounded-full">
                ตั้งเวลาเอง
              </span>
            </div>
            <div className="flex items-center text-violet-400 text-sm font-medium group-hover:gap-2 transition-all gap-1">
              <span>เลือกโหมดนี้</span>
              <span>→</span>
            </div>
          </button>

          {/* Online 1vs1 */}
          <button
            onClick={() => router.push("/play/ranked")}
            className="text-left bg-[#12122a] border border-white/5 rounded-2xl p-5 sm:p-8 hover:border-blue-500/40 transition-all group cursor-pointer"
          >
            <div className="flex items-center gap-3 mb-4">
              <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center text-2xl group-hover:scale-110 transition-transform">
                ⚔️
              </div>
              <div>
                <h2 className="text-xl font-bold text-white">ออนไลน์ 1vs1</h2>
                <p className="text-xs text-gray-500">Ranked Multiplayer</p>
              </div>
            </div>
            <p className="text-gray-400 text-sm leading-relaxed mb-6">
              ท้าทายผู้เล่นคนอื่นแบบเรียลไทม์ แข่งขันตัวต่อตัว
              และไต่อันดับ Rating ด้วยระบบ ELO
            </p>
            <div className="flex flex-wrap gap-2 mb-6">
              <span className="text-xs bg-blue-500/10 text-blue-400 px-2.5 py-1 rounded-full">
                จับคู่เรียลไทม์
              </span>
              <span className="text-xs bg-blue-500/10 text-blue-400 px-2.5 py-1 rounded-full">
                ระบบ ELO Rating
              </span>
              <span className="text-xs bg-blue-500/10 text-blue-400 px-2.5 py-1 rounded-full">
                10 ข้อ 3 ระดับ
              </span>
            </div>
            <div className="flex items-center text-blue-400 text-sm font-medium group-hover:gap-2 transition-all gap-1">
              <span>เลือกโหมดนี้</span>
              <span>→</span>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}

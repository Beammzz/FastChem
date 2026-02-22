"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";
import Navbar from "@/components/Navbar";

export default function Home() {
  const router = useRouter();

  return (
    <div className="min-h-screen-safe bg-[#0a0a1a] text-white">
      <Navbar />

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Background glow */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-violet-500/10 rounded-full blur-[120px] pointer-events-none" />

        <div className="max-w-6xl mx-auto px-4 sm:px-6 pt-12 sm:pt-20 pb-12 sm:pb-16 text-center relative">
          <h1 className="text-3xl sm:text-5xl md:text-7xl font-extrabold mb-4">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-violet-400 via-purple-400 to-fuchsia-500">
              ฝึกเคมี
            </span>
            <br />
            <span className="text-white/90">สำหรับทุกคน</span>
          </h1>
          <p className="text-base sm:text-lg md:text-xl text-gray-400 max-w-2xl mx-auto mt-4 leading-relaxed">
            ฝึกโครงสร้างอะตอม เลขออกซิเดชัน โมลคอนเซ็ปต์ และอื่นๆ
            ด้วยคำถามที่สร้างอัตโนมัติ — เร็ว ฟรี!
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4 mt-8 px-4 sm:px-0">
            <button
              onClick={() => router.push("/play")}
              className="w-full sm:w-auto bg-violet-500 hover:bg-violet-600 text-white font-bold text-lg px-8 py-3.5 rounded-xl transition-all hover:scale-[1.02] active:scale-[0.98] shadow-lg shadow-violet-500/20 min-h-[44px]"
            >
              เล่นเลย
            </button>
            <Link
              href="#topics"
              className="w-full sm:w-auto text-center border border-white/10 hover:border-white/20 text-gray-300 hover:text-white font-medium px-8 py-3.5 rounded-xl transition-all min-h-[44px] flex items-center justify-center"
            >
              เรียนรู้เพิ่มเติม
            </Link>
          </div>
        </div>
      </section>

      {/* Features / Topics */}
      <section id="topics" className="max-w-6xl mx-auto px-4 sm:px-6 pb-16 pt-12">
        <h3 className="text-xs text-gray-500 uppercase tracking-widest font-medium text-center mb-8">
          หัวข้อที่ครอบคลุม
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            { icon: "⚛️", label: "โปรตอน", desc: "เลขอะตอม" },
            { icon: "🔵", label: "นิวตรอน", desc: "มวล − เลขอะตอม" },
            { icon: "⚡", label: "อิเล็กตรอน", desc: "ประจุอะตอมเป็นกลาง" },
            { icon: "🔢", label: "ออกซิเดชัน", desc: "วิเคราะห์สารประกอบ" },
            { icon: "⚖️", label: "โมล", desc: "มวล↔โมล↔อนุภาค" },
            { icon: "🧪", label: "การเจือจาง", desc: "M₁V₁ = M₂V₂" },
            { icon: "🧂", label: "เตรียมสารละลาย", desc: "จากของแข็ง" },
            { icon: "❄️", label: "จุดเยือกแข็ง", desc: "ΔTf = iKfm" },
          ].map((topic) => (
            <div
              key={topic.label}
              className="bg-[#12122a] border border-white/5 rounded-xl p-5 text-center hover:border-violet-500/20 transition-all"
            >
              <div className="text-3xl mb-2">{topic.icon}</div>
              <div className="text-sm font-bold text-white">{topic.label}</div>
              <div className="text-xs text-gray-500 mt-1">{topic.desc}</div>
            </div>
          ))}
        </div>
      </section>

      {/* Rules */}
      <section className="max-w-6xl mx-auto px-4 sm:px-6 pb-16">
        <div className="bg-[#12122a] border border-white/5 rounded-2xl p-5 sm:p-8">
          <h3 className="text-lg font-bold text-white mb-4">วิธีเล่น</h3>
          <div className="grid md:grid-cols-4 gap-6">
            {[
              { step: "01", text: "เลือกโหมดเกมและระดับความยาก" },
              { step: "02", text: "ตอบคำถามเคมีแบบปรนัย" },
              { step: "03", text: "ได้รับคะแนนสำหรับคำตอบที่ถูกต้อง" },
              { step: "04", text: "ทำคะแนนให้มากที่สุดก่อนหมดเวลา" },
            ].map((item) => (
              <div key={item.step} className="flex gap-3">
                <span className="text-2xl font-black text-violet-500/30">{item.step}</span>
                <p className="text-sm text-gray-400 leading-relaxed">{item.text}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-white/5 py-6 sm:py-8 pb-safe">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-2">
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <span>⚗️</span>
            <span>FastChem v1.0</span>
          </div>
          <p className="text-xs text-gray-600">
            © 2026 FastChem. สงวนลิขสิทธิ์
          </p>
        </div>
      </footer>
    </div>
  );
}

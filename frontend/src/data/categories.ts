/**
 * Question categories, mirroring the topic registry in
 * `backend/internal/services/problem.go`.
 *
 * The `id` values are a cross-stack contract: the backend stamps them onto
 * every question and the single-player page sends them back as filters, so a
 * rename is a two-sided change. Adding a topic on the server without adding it
 * here leaves its questions rendering a raw snake_case id as their label.
 *
 * Chapters follow สาระเคมี in หลักสูตรแกนกลางการศึกษาขั้นพื้นฐาน พ.ศ. 2551
 * (ฉบับปรับปรุง พ.ศ. 2560), which is the curriculum ม.4–ม.6 is taught from.
 */

export type CategoryDifficulty = "easy" | "medium" | "hard";

export interface CategoryMeta {
  id: string;
  label: string;
  desc: string;
  icon: string;
  difficulty: CategoryDifficulty;
  chapter: string;
  /** Tailwind classes for the category chip shown above a question. */
  badgeClass: string;
}

/** Chapter headings in curriculum order; the picker groups by these. */
export const CHAPTERS = [
  "พื้นฐาน",
  "บทที่ 2 อะตอมและสมบัติของธาตุ",
  "บทที่ 3 พันธะเคมี",
  "บทที่ 4 โมลและสูตรเคมี",
  "บทที่ 5 สารละลาย",
  "บทที่ 6 ปริมาณสัมพันธ์",
  "บทที่ 7 แก๊ส",
  "บทที่ 8 อัตราการเกิดปฏิกิริยาเคมี",
  "บทที่ 9 สมดุลเคมี",
  "บทที่ 10 กรด–เบส",
  "บทที่ 11 เคมีไฟฟ้า",
  "บทที่ 12 เคมีอินทรีย์",
  "บทที่ 13 พอลิเมอร์",
] as const;

export const CATEGORIES: CategoryMeta[] = [
  {
    id: "state_of_matter",
    label: "สถานะของสาร",
    desc: "ของแข็ง / ของเหลว / แก๊ส / สารละลาย (aq)",
    icon: "🧊",
    difficulty: "easy",
    chapter: "พื้นฐาน",
    badgeClass: "bg-teal-500/10 text-teal-400 border-teal-500/20",
  },
  {
    id: "atomic_structure",
    label: "โครงสร้างอะตอม",
    desc: "โปรตอน, นิวตรอน, อิเล็กตรอน, ไอโซโทป",
    icon: "⚛️",
    difficulty: "easy",
    chapter: "บทที่ 2 อะตอมและสมบัติของธาตุ",
    badgeClass: "bg-cyan-500/10 text-cyan-400 border-cyan-500/20",
  },
  {
    id: "electron_config",
    label: "การจัดเรียงอิเล็กตรอน",
    desc: "ระดับพลังงาน, เวเลนซ์อิเล็กตรอน, หมู่และคาบ",
    icon: "🌀",
    difficulty: "easy",
    chapter: "บทที่ 2 อะตอมและสมบัติของธาตุ",
    badgeClass: "bg-blue-500/10 text-blue-400 border-blue-500/20",
  },
  {
    id: "bond_type",
    label: "ชนิดพันธะเคมี",
    desc: "ไอออนิก / โคเวเลนต์มีขั้ว-ไม่มีขั้ว / โลหะ",
    icon: "🔗",
    difficulty: "easy",
    chapter: "บทที่ 3 พันธะเคมี",
    badgeClass: "bg-indigo-500/10 text-indigo-400 border-indigo-500/20",
  },
  {
    id: "molecular_shape",
    label: "รูปร่างโมเลกุล",
    desc: "VSEPR: เส้นตรง, มุมงอ, ทรงสี่หน้า",
    icon: "🔷",
    difficulty: "easy",
    chapter: "บทที่ 3 พันธะเคมี",
    badgeClass: "bg-violet-500/10 text-violet-400 border-violet-500/20",
  },
  {
    id: "mole_concept",
    label: "โมลคอนเซ็ปต์",
    desc: "มวล↔โมล↔อนุภาค",
    icon: "⚖️",
    difficulty: "medium",
    chapter: "บทที่ 4 โมลและสูตรเคมี",
    badgeClass: "bg-amber-500/10 text-amber-400 border-amber-500/20",
  },
  {
    id: "molar_mass",
    label: "มวลต่อโมล",
    desc: "คำนวณมวลโมเลกุลจากสูตรเคมี",
    icon: "🧮",
    difficulty: "medium",
    chapter: "บทที่ 4 โมลและสูตรเคมี",
    badgeClass: "bg-yellow-500/10 text-yellow-400 border-yellow-500/20",
  },
  {
    id: "percent_composition",
    label: "ร้อยละโดยมวล",
    desc: "สัดส่วนมวลของธาตุในสารประกอบ",
    icon: "📊",
    difficulty: "medium",
    chapter: "บทที่ 4 โมลและสูตรเคมี",
    badgeClass: "bg-lime-500/10 text-lime-400 border-lime-500/20",
  },
  {
    id: "concentration",
    label: "ความเข้มข้นสารละลาย",
    desc: "โมลาร์, โมแลล, ร้อยละโดยมวล, ppm",
    icon: "🥤",
    difficulty: "medium",
    chapter: "บทที่ 5 สารละลาย",
    badgeClass: "bg-emerald-500/10 text-emerald-400 border-emerald-500/20",
  },
  {
    id: "dilution",
    label: "การเจือจาง",
    desc: "M₁V₁ = M₂V₂",
    icon: "🧪",
    difficulty: "hard",
    chapter: "บทที่ 5 สารละลาย",
    badgeClass: "bg-rose-500/10 text-rose-400 border-rose-500/20",
  },
  {
    id: "preparing_solution",
    label: "เตรียมสารละลาย",
    desc: "mass = M × V × MW",
    icon: "🧂",
    difficulty: "hard",
    chapter: "บทที่ 5 สารละลาย",
    badgeClass: "bg-orange-500/10 text-orange-400 border-orange-500/20",
  },
  {
    id: "freezing_point",
    label: "จุดเยือกแข็ง",
    desc: "ΔTf = i·Kf·m",
    icon: "❄️",
    difficulty: "hard",
    chapter: "บทที่ 5 สารละลาย",
    badgeClass: "bg-sky-500/10 text-sky-400 border-sky-500/20",
  },
  {
    id: "stoichiometry",
    label: "ปริมาณสัมพันธ์",
    desc: "อัตราส่วนโมลจากสมการเคมี",
    icon: "⚗️",
    difficulty: "medium",
    chapter: "บทที่ 6 ปริมาณสัมพันธ์",
    badgeClass: "bg-green-500/10 text-green-400 border-green-500/20",
  },
  {
    id: "limiting_reagent",
    label: "สารกำหนดปริมาณ",
    desc: "สารที่หมดก่อน และร้อยละผลได้",
    icon: "🎯",
    difficulty: "hard",
    chapter: "บทที่ 6 ปริมาณสัมพันธ์",
    badgeClass: "bg-red-500/10 text-red-400 border-red-500/20",
  },
  {
    id: "gas_law",
    label: "กฎของแก๊ส",
    desc: "บอยล์, ชาร์ล, กฎรวม, ปริมาตรที่ STP",
    icon: "🎈",
    difficulty: "medium",
    chapter: "บทที่ 7 แก๊ส",
    badgeClass: "bg-fuchsia-500/10 text-fuchsia-400 border-fuchsia-500/20",
  },
  {
    id: "ideal_gas",
    label: "แก๊สอุดมคติ",
    desc: "PV = nRT",
    icon: "💨",
    difficulty: "hard",
    chapter: "บทที่ 7 แก๊ส",
    badgeClass: "bg-pink-500/10 text-pink-400 border-pink-500/20",
  },
  {
    id: "reaction_rate",
    label: "อัตราการเกิดปฏิกิริยา",
    desc: "อัตราเฉลี่ย และอันดับของปฏิกิริยา",
    icon: "⏱️",
    difficulty: "hard",
    chapter: "บทที่ 8 อัตราการเกิดปฏิกิริยาเคมี",
    badgeClass: "bg-orange-500/10 text-orange-300 border-orange-500/20",
  },
  {
    id: "equilibrium",
    label: "สมดุลเคมี",
    desc: "ค่าคงที่สมดุล K และการเขียนนิพจน์",
    icon: "⚖",
    difficulty: "hard",
    chapter: "บทที่ 9 สมดุลเคมี",
    badgeClass: "bg-amber-500/10 text-amber-300 border-amber-500/20",
  },
  {
    id: "acid_base",
    label: "กรด–เบส",
    desc: "pH, pOH, Ka และกรดแก่-กรดอ่อน",
    icon: "🍋",
    difficulty: "hard",
    chapter: "บทที่ 10 กรด–เบส",
    badgeClass: "bg-rose-500/10 text-rose-300 border-rose-500/20",
  },
  {
    id: "oxidation_number",
    label: "เลขออกซิเดชัน",
    desc: "วิเคราะห์สารประกอบ",
    icon: "🔢",
    difficulty: "easy",
    chapter: "บทที่ 11 เคมีไฟฟ้า",
    badgeClass: "bg-purple-500/10 text-purple-400 border-purple-500/20",
  },
  {
    id: "electrochemistry",
    label: "เคมีไฟฟ้า",
    desc: "E°cell และความแรงของตัวออกซิไดซ์-รีดิวซ์",
    icon: "🔋",
    difficulty: "hard",
    chapter: "บทที่ 11 เคมีไฟฟ้า",
    badgeClass: "bg-purple-500/10 text-purple-300 border-purple-500/20",
  },
  {
    id: "functional_group",
    label: "หมู่ฟังก์ชัน",
    desc: "ประเภทสารอินทรีย์และหมู่ฟังก์ชัน",
    icon: "🧬",
    difficulty: "easy",
    chapter: "บทที่ 12 เคมีอินทรีย์",
    badgeClass: "bg-teal-500/10 text-teal-300 border-teal-500/20",
  },
  {
    id: "polymer",
    label: "พอลิเมอร์",
    desc: "มอนอเมอร์ และแบบเติม-แบบควบแน่น",
    icon: "🪢",
    difficulty: "easy",
    chapter: "บทที่ 13 พอลิเมอร์",
    badgeClass: "bg-sky-500/10 text-sky-300 border-sky-500/20",
  },
];

const BY_ID = new Map(CATEGORIES.map((c) => [c.id, c]));

/**
 * Thai label for a category id, falling back to the raw id so an unknown
 * category renders something rather than an empty chip.
 */
export function categoryLabel(id: string): string {
  return BY_ID.get(id)?.label ?? id;
}

export function categoryBadgeClass(id: string): string {
  return (
    BY_ID.get(id)?.badgeClass ??
    "bg-gray-500/10 text-gray-400 border-gray-500/20"
  );
}

export function categoriesForChapter(chapter: string): CategoryMeta[] {
  return CATEGORIES.filter((c) => c.chapter === chapter);
}

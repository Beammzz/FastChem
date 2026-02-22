// Periodic table data for the periodic table component.
// Frontend-only — not fetched from the backend.
// All 118 elements included.

export interface Element {
  atomicNumber: number;
  symbol: string;
  name: string;
  nameTH: string;
  atomicMass: number;
  category: ElementCategory;
  group: number; // column (1-18)
  period: number; // row (1-7), 8 = lanthanide row, 9 = actinide row
}

export type ElementCategory =
  | "nonmetal"
  | "noble-gas"
  | "alkali-metal"
  | "alkaline-earth"
  | "metalloid"
  | "halogen"
  | "transition-metal"
  | "post-transition"
  | "lanthanide"
  | "actinide";

export const CATEGORY_COLORS: Record<ElementCategory, string> = {
  "nonmetal": "bg-green-600/80 border-green-500/40",
  "noble-gas": "bg-purple-600/80 border-purple-500/40",
  "alkali-metal": "bg-red-600/80 border-red-500/40",
  "alkaline-earth": "bg-orange-600/80 border-orange-500/40",
  "metalloid": "bg-teal-600/80 border-teal-500/40",
  "halogen": "bg-cyan-600/80 border-cyan-500/40",
  "transition-metal": "bg-blue-600/80 border-blue-500/40",
  "post-transition": "bg-indigo-600/80 border-indigo-500/40",
  "lanthanide": "bg-pink-600/80 border-pink-500/40",
  "actinide": "bg-rose-600/80 border-rose-500/40",
};

export const CATEGORY_LABELS: Record<ElementCategory, string> = {
  "nonmetal": "อโลหะ",
  "noble-gas": "แก๊สเฉื่อย",
  "alkali-metal": "โลหะแอลคาไล",
  "alkaline-earth": "โลหะแอลคาไลน์เอิร์ท",
  "metalloid": "กึ่งโลหะ",
  "halogen": "แฮโลเจน",
  "transition-metal": "โลหะทรานซิชัน",
  "post-transition": "โลหะหลังทรานซิชัน",
  "lanthanide": "แลนทาไนด์",
  "actinide": "แอกทิไนด์",
};

// All 118 elements positioned in the standard periodic table grid
// Lanthanides (57-71) use period 8 (displayed as separate row)
// Actinides (89-103) use period 9 (displayed as separate row)
export const ELEMENTS: Element[] = [
  // ── Period 1 ──
  { atomicNumber: 1, symbol: "H", name: "Hydrogen", nameTH: "ไฮโดรเจน", atomicMass: 1.008, category: "nonmetal", group: 1, period: 1 },
  { atomicNumber: 2, symbol: "He", name: "Helium", nameTH: "ฮีเลียม", atomicMass: 4.003, category: "noble-gas", group: 18, period: 1 },

  // ── Period 2 ──
  { atomicNumber: 3, symbol: "Li", name: "Lithium", nameTH: "ลิเทียม", atomicMass: 6.941, category: "alkali-metal", group: 1, period: 2 },
  { atomicNumber: 4, symbol: "Be", name: "Beryllium", nameTH: "เบริลเลียม", atomicMass: 9.012, category: "alkaline-earth", group: 2, period: 2 },
  { atomicNumber: 5, symbol: "B", name: "Boron", nameTH: "โบรอน", atomicMass: 10.81, category: "metalloid", group: 13, period: 2 },
  { atomicNumber: 6, symbol: "C", name: "Carbon", nameTH: "คาร์บอน", atomicMass: 12.01, category: "nonmetal", group: 14, period: 2 },
  { atomicNumber: 7, symbol: "N", name: "Nitrogen", nameTH: "ไนโตรเจน", atomicMass: 14.01, category: "nonmetal", group: 15, period: 2 },
  { atomicNumber: 8, symbol: "O", name: "Oxygen", nameTH: "ออกซิเจน", atomicMass: 16.00, category: "nonmetal", group: 16, period: 2 },
  { atomicNumber: 9, symbol: "F", name: "Fluorine", nameTH: "ฟลูออรีน", atomicMass: 19.00, category: "halogen", group: 17, period: 2 },
  { atomicNumber: 10, symbol: "Ne", name: "Neon", nameTH: "นีออน", atomicMass: 20.18, category: "noble-gas", group: 18, period: 2 },

  // ── Period 3 ──
  { atomicNumber: 11, symbol: "Na", name: "Sodium", nameTH: "โซเดียม", atomicMass: 22.99, category: "alkali-metal", group: 1, period: 3 },
  { atomicNumber: 12, symbol: "Mg", name: "Magnesium", nameTH: "แมกนีเซียม", atomicMass: 24.31, category: "alkaline-earth", group: 2, period: 3 },
  { atomicNumber: 13, symbol: "Al", name: "Aluminium", nameTH: "อะลูมิเนียม", atomicMass: 26.98, category: "post-transition", group: 13, period: 3 },
  { atomicNumber: 14, symbol: "Si", name: "Silicon", nameTH: "ซิลิคอน", atomicMass: 28.09, category: "metalloid", group: 14, period: 3 },
  { atomicNumber: 15, symbol: "P", name: "Phosphorus", nameTH: "ฟอสฟอรัส", atomicMass: 30.97, category: "nonmetal", group: 15, period: 3 },
  { atomicNumber: 16, symbol: "S", name: "Sulfur", nameTH: "กำมะถัน", atomicMass: 32.07, category: "nonmetal", group: 16, period: 3 },
  { atomicNumber: 17, symbol: "Cl", name: "Chlorine", nameTH: "คลอรีน", atomicMass: 35.45, category: "halogen", group: 17, period: 3 },
  { atomicNumber: 18, symbol: "Ar", name: "Argon", nameTH: "อาร์กอน", atomicMass: 39.95, category: "noble-gas", group: 18, period: 3 },

  // ── Period 4 ──
  { atomicNumber: 19, symbol: "K", name: "Potassium", nameTH: "โพแทสเซียม", atomicMass: 39.10, category: "alkali-metal", group: 1, period: 4 },
  { atomicNumber: 20, symbol: "Ca", name: "Calcium", nameTH: "แคลเซียม", atomicMass: 40.08, category: "alkaline-earth", group: 2, period: 4 },
  { atomicNumber: 21, symbol: "Sc", name: "Scandium", nameTH: "สแกนเดียม", atomicMass: 44.96, category: "transition-metal", group: 3, period: 4 },
  { atomicNumber: 22, symbol: "Ti", name: "Titanium", nameTH: "ไทเทเนียม", atomicMass: 47.87, category: "transition-metal", group: 4, period: 4 },
  { atomicNumber: 23, symbol: "V", name: "Vanadium", nameTH: "วาเนเดียม", atomicMass: 50.94, category: "transition-metal", group: 5, period: 4 },
  { atomicNumber: 24, symbol: "Cr", name: "Chromium", nameTH: "โครเมียม", atomicMass: 52.00, category: "transition-metal", group: 6, period: 4 },
  { atomicNumber: 25, symbol: "Mn", name: "Manganese", nameTH: "แมงกานีส", atomicMass: 54.94, category: "transition-metal", group: 7, period: 4 },
  { atomicNumber: 26, symbol: "Fe", name: "Iron", nameTH: "เหล็ก", atomicMass: 55.85, category: "transition-metal", group: 8, period: 4 },
  { atomicNumber: 27, symbol: "Co", name: "Cobalt", nameTH: "โคบอลต์", atomicMass: 58.93, category: "transition-metal", group: 9, period: 4 },
  { atomicNumber: 28, symbol: "Ni", name: "Nickel", nameTH: "นิกเกิล", atomicMass: 58.69, category: "transition-metal", group: 10, period: 4 },
  { atomicNumber: 29, symbol: "Cu", name: "Copper", nameTH: "ทองแดง", atomicMass: 63.55, category: "transition-metal", group: 11, period: 4 },
  { atomicNumber: 30, symbol: "Zn", name: "Zinc", nameTH: "สังกะสี", atomicMass: 65.38, category: "transition-metal", group: 12, period: 4 },
  { atomicNumber: 31, symbol: "Ga", name: "Gallium", nameTH: "แกลเลียม", atomicMass: 69.72, category: "post-transition", group: 13, period: 4 },
  { atomicNumber: 32, symbol: "Ge", name: "Germanium", nameTH: "เจอร์เมเนียม", atomicMass: 72.63, category: "metalloid", group: 14, period: 4 },
  { atomicNumber: 33, symbol: "As", name: "Arsenic", nameTH: "สารหนู", atomicMass: 74.92, category: "metalloid", group: 15, period: 4 },
  { atomicNumber: 34, symbol: "Se", name: "Selenium", nameTH: "ซีลีเนียม", atomicMass: 78.97, category: "nonmetal", group: 16, period: 4 },
  { atomicNumber: 35, symbol: "Br", name: "Bromine", nameTH: "โบรมีน", atomicMass: 79.90, category: "halogen", group: 17, period: 4 },
  { atomicNumber: 36, symbol: "Kr", name: "Krypton", nameTH: "คริปทอน", atomicMass: 83.80, category: "noble-gas", group: 18, period: 4 },

  // ── Period 5 ──
  { atomicNumber: 37, symbol: "Rb", name: "Rubidium", nameTH: "รูบิเดียม", atomicMass: 85.47, category: "alkali-metal", group: 1, period: 5 },
  { atomicNumber: 38, symbol: "Sr", name: "Strontium", nameTH: "สทรอนเชียม", atomicMass: 87.62, category: "alkaline-earth", group: 2, period: 5 },
  { atomicNumber: 39, symbol: "Y", name: "Yttrium", nameTH: "อิตเทรียม", atomicMass: 88.91, category: "transition-metal", group: 3, period: 5 },
  { atomicNumber: 40, symbol: "Zr", name: "Zirconium", nameTH: "เซอร์โคเนียม", atomicMass: 91.22, category: "transition-metal", group: 4, period: 5 },
  { atomicNumber: 41, symbol: "Nb", name: "Niobium", nameTH: "ไนโอเบียม", atomicMass: 92.91, category: "transition-metal", group: 5, period: 5 },
  { atomicNumber: 42, symbol: "Mo", name: "Molybdenum", nameTH: "โมลิบดีนัม", atomicMass: 95.95, category: "transition-metal", group: 6, period: 5 },
  { atomicNumber: 43, symbol: "Tc", name: "Technetium", nameTH: "เทคนีเชียม", atomicMass: 98, category: "transition-metal", group: 7, period: 5 },
  { atomicNumber: 44, symbol: "Ru", name: "Ruthenium", nameTH: "รูทีเนียม", atomicMass: 101.1, category: "transition-metal", group: 8, period: 5 },
  { atomicNumber: 45, symbol: "Rh", name: "Rhodium", nameTH: "โรเดียม", atomicMass: 102.9, category: "transition-metal", group: 9, period: 5 },
  { atomicNumber: 46, symbol: "Pd", name: "Palladium", nameTH: "แพลเลเดียม", atomicMass: 106.4, category: "transition-metal", group: 10, period: 5 },
  { atomicNumber: 47, symbol: "Ag", name: "Silver", nameTH: "เงิน", atomicMass: 107.9, category: "transition-metal", group: 11, period: 5 },
  { atomicNumber: 48, symbol: "Cd", name: "Cadmium", nameTH: "แคดเมียม", atomicMass: 112.4, category: "transition-metal", group: 12, period: 5 },
  { atomicNumber: 49, symbol: "In", name: "Indium", nameTH: "อินเดียม", atomicMass: 114.8, category: "post-transition", group: 13, period: 5 },
  { atomicNumber: 50, symbol: "Sn", name: "Tin", nameTH: "ดีบุก", atomicMass: 118.7, category: "post-transition", group: 14, period: 5 },
  { atomicNumber: 51, symbol: "Sb", name: "Antimony", nameTH: "พลวง", atomicMass: 121.8, category: "metalloid", group: 15, period: 5 },
  { atomicNumber: 52, symbol: "Te", name: "Tellurium", nameTH: "เทลลูเรียม", atomicMass: 127.6, category: "metalloid", group: 16, period: 5 },
  { atomicNumber: 53, symbol: "I", name: "Iodine", nameTH: "ไอโอดีน", atomicMass: 126.9, category: "halogen", group: 17, period: 5 },
  { atomicNumber: 54, symbol: "Xe", name: "Xenon", nameTH: "ซีนอน", atomicMass: 131.3, category: "noble-gas", group: 18, period: 5 },

  // ── Period 6 ──
  { atomicNumber: 55, symbol: "Cs", name: "Caesium", nameTH: "ซีเซียม", atomicMass: 132.9, category: "alkali-metal", group: 1, period: 6 },
  { atomicNumber: 56, symbol: "Ba", name: "Barium", nameTH: "แบเรียม", atomicMass: 137.3, category: "alkaline-earth", group: 2, period: 6 },
  // La–Lu → lanthanide row (period 8 for display)
  { atomicNumber: 57, symbol: "La", name: "Lanthanum", nameTH: "แลนทานัม", atomicMass: 138.9, category: "lanthanide", group: 3, period: 8 },
  { atomicNumber: 58, symbol: "Ce", name: "Cerium", nameTH: "ซีเรียม", atomicMass: 140.1, category: "lanthanide", group: 4, period: 8 },
  { atomicNumber: 59, symbol: "Pr", name: "Praseodymium", nameTH: "เพรซีโอดิเมียม", atomicMass: 140.9, category: "lanthanide", group: 5, period: 8 },
  { atomicNumber: 60, symbol: "Nd", name: "Neodymium", nameTH: "นีโอดิเมียม", atomicMass: 144.2, category: "lanthanide", group: 6, period: 8 },
  { atomicNumber: 61, symbol: "Pm", name: "Promethium", nameTH: "โพรมีเทียม", atomicMass: 145, category: "lanthanide", group: 7, period: 8 },
  { atomicNumber: 62, symbol: "Sm", name: "Samarium", nameTH: "ซาแมเรียม", atomicMass: 150.4, category: "lanthanide", group: 8, period: 8 },
  { atomicNumber: 63, symbol: "Eu", name: "Europium", nameTH: "ยูโรเพียม", atomicMass: 152.0, category: "lanthanide", group: 9, period: 8 },
  { atomicNumber: 64, symbol: "Gd", name: "Gadolinium", nameTH: "แกโดลิเนียม", atomicMass: 157.3, category: "lanthanide", group: 10, period: 8 },
  { atomicNumber: 65, symbol: "Tb", name: "Terbium", nameTH: "เทอร์เบียม", atomicMass: 158.9, category: "lanthanide", group: 11, period: 8 },
  { atomicNumber: 66, symbol: "Dy", name: "Dysprosium", nameTH: "ดิสโพรเซียม", atomicMass: 162.5, category: "lanthanide", group: 12, period: 8 },
  { atomicNumber: 67, symbol: "Ho", name: "Holmium", nameTH: "โฮลเมียม", atomicMass: 164.9, category: "lanthanide", group: 13, period: 8 },
  { atomicNumber: 68, symbol: "Er", name: "Erbium", nameTH: "เออร์เบียม", atomicMass: 167.3, category: "lanthanide", group: 14, period: 8 },
  { atomicNumber: 69, symbol: "Tm", name: "Thulium", nameTH: "ทูเลียม", atomicMass: 168.9, category: "lanthanide", group: 15, period: 8 },
  { atomicNumber: 70, symbol: "Yb", name: "Ytterbium", nameTH: "อิตเทอร์เบียม", atomicMass: 173.0, category: "lanthanide", group: 16, period: 8 },
  { atomicNumber: 71, symbol: "Lu", name: "Lutetium", nameTH: "ลูทีเชียม", atomicMass: 175.0, category: "lanthanide", group: 17, period: 8 },
  // Hf onward in period 6
  { atomicNumber: 72, symbol: "Hf", name: "Hafnium", nameTH: "แฮฟเนียม", atomicMass: 178.5, category: "transition-metal", group: 4, period: 6 },
  { atomicNumber: 73, symbol: "Ta", name: "Tantalum", nameTH: "แทนทาลัม", atomicMass: 180.9, category: "transition-metal", group: 5, period: 6 },
  { atomicNumber: 74, symbol: "W", name: "Tungsten", nameTH: "ทังสเตน", atomicMass: 183.8, category: "transition-metal", group: 6, period: 6 },
  { atomicNumber: 75, symbol: "Re", name: "Rhenium", nameTH: "รีเนียม", atomicMass: 186.2, category: "transition-metal", group: 7, period: 6 },
  { atomicNumber: 76, symbol: "Os", name: "Osmium", nameTH: "ออสเมียม", atomicMass: 190.2, category: "transition-metal", group: 8, period: 6 },
  { atomicNumber: 77, symbol: "Ir", name: "Iridium", nameTH: "อิริเดียม", atomicMass: 192.2, category: "transition-metal", group: 9, period: 6 },
  { atomicNumber: 78, symbol: "Pt", name: "Platinum", nameTH: "แพลทินัม", atomicMass: 195.1, category: "transition-metal", group: 10, period: 6 },
  { atomicNumber: 79, symbol: "Au", name: "Gold", nameTH: "ทองคำ", atomicMass: 197.0, category: "transition-metal", group: 11, period: 6 },
  { atomicNumber: 80, symbol: "Hg", name: "Mercury", nameTH: "ปรอท", atomicMass: 200.6, category: "transition-metal", group: 12, period: 6 },
  { atomicNumber: 81, symbol: "Tl", name: "Thallium", nameTH: "แทลเลียม", atomicMass: 204.4, category: "post-transition", group: 13, period: 6 },
  { atomicNumber: 82, symbol: "Pb", name: "Lead", nameTH: "ตะกั่ว", atomicMass: 207.2, category: "post-transition", group: 14, period: 6 },
  { atomicNumber: 83, symbol: "Bi", name: "Bismuth", nameTH: "บิสมัท", atomicMass: 209.0, category: "post-transition", group: 15, period: 6 },
  { atomicNumber: 84, symbol: "Po", name: "Polonium", nameTH: "พอโลเนียม", atomicMass: 209, category: "post-transition", group: 16, period: 6 },
  { atomicNumber: 85, symbol: "At", name: "Astatine", nameTH: "แอสทาทีน", atomicMass: 210, category: "halogen", group: 17, period: 6 },
  { atomicNumber: 86, symbol: "Rn", name: "Radon", nameTH: "เรดอน", atomicMass: 222, category: "noble-gas", group: 18, period: 6 },

  // ── Period 7 ──
  { atomicNumber: 87, symbol: "Fr", name: "Francium", nameTH: "แฟรนเซียม", atomicMass: 223, category: "alkali-metal", group: 1, period: 7 },
  { atomicNumber: 88, symbol: "Ra", name: "Radium", nameTH: "เรเดียม", atomicMass: 226, category: "alkaline-earth", group: 2, period: 7 },
  // Ac–Lr → actinide row (period 9 for display)
  { atomicNumber: 89, symbol: "Ac", name: "Actinium", nameTH: "แอกทิเนียม", atomicMass: 227, category: "actinide", group: 3, period: 9 },
  { atomicNumber: 90, symbol: "Th", name: "Thorium", nameTH: "ทอเรียม", atomicMass: 232.0, category: "actinide", group: 4, period: 9 },
  { atomicNumber: 91, symbol: "Pa", name: "Protactinium", nameTH: "โพรแทกทิเนียม", atomicMass: 231.0, category: "actinide", group: 5, period: 9 },
  { atomicNumber: 92, symbol: "U", name: "Uranium", nameTH: "ยูเรเนียม", atomicMass: 238.0, category: "actinide", group: 6, period: 9 },
  { atomicNumber: 93, symbol: "Np", name: "Neptunium", nameTH: "เนปทูเนียม", atomicMass: 237, category: "actinide", group: 7, period: 9 },
  { atomicNumber: 94, symbol: "Pu", name: "Plutonium", nameTH: "พลูโทเนียม", atomicMass: 244, category: "actinide", group: 8, period: 9 },
  { atomicNumber: 95, symbol: "Am", name: "Americium", nameTH: "อะเมริเซียม", atomicMass: 243, category: "actinide", group: 9, period: 9 },
  { atomicNumber: 96, symbol: "Cm", name: "Curium", nameTH: "คูเรียม", atomicMass: 247, category: "actinide", group: 10, period: 9 },
  { atomicNumber: 97, symbol: "Bk", name: "Berkelium", nameTH: "เบอร์คีเลียม", atomicMass: 247, category: "actinide", group: 11, period: 9 },
  { atomicNumber: 98, symbol: "Cf", name: "Californium", nameTH: "แคลิฟอร์เนียม", atomicMass: 251, category: "actinide", group: 12, period: 9 },
  { atomicNumber: 99, symbol: "Es", name: "Einsteinium", nameTH: "ไอน์สไตเนียม", atomicMass: 252, category: "actinide", group: 13, period: 9 },
  { atomicNumber: 100, symbol: "Fm", name: "Fermium", nameTH: "เฟอร์เมียม", atomicMass: 257, category: "actinide", group: 14, period: 9 },
  { atomicNumber: 101, symbol: "Md", name: "Mendelevium", nameTH: "เมนเดลีเวียม", atomicMass: 258, category: "actinide", group: 15, period: 9 },
  { atomicNumber: 102, symbol: "No", name: "Nobelium", nameTH: "โนเบเลียม", atomicMass: 259, category: "actinide", group: 16, period: 9 },
  { atomicNumber: 103, symbol: "Lr", name: "Lawrencium", nameTH: "ลอว์เรนเซียม", atomicMass: 266, category: "actinide", group: 17, period: 9 },
  // Rf onward in period 7
  { atomicNumber: 104, symbol: "Rf", name: "Rutherfordium", nameTH: "รัทเทอร์ฟอร์เดียม", atomicMass: 267, category: "transition-metal", group: 4, period: 7 },
  { atomicNumber: 105, symbol: "Db", name: "Dubnium", nameTH: "ดุบเนียม", atomicMass: 268, category: "transition-metal", group: 5, period: 7 },
  { atomicNumber: 106, symbol: "Sg", name: "Seaborgium", nameTH: "ซีบอร์เกียม", atomicMass: 269, category: "transition-metal", group: 6, period: 7 },
  { atomicNumber: 107, symbol: "Bh", name: "Bohrium", nameTH: "โบห์เรียม", atomicMass: 270, category: "transition-metal", group: 7, period: 7 },
  { atomicNumber: 108, symbol: "Hs", name: "Hassium", nameTH: "ฮัสเซียม", atomicMass: 277, category: "transition-metal", group: 8, period: 7 },
  { atomicNumber: 109, symbol: "Mt", name: "Meitnerium", nameTH: "ไมต์เนเรียม", atomicMass: 278, category: "transition-metal", group: 9, period: 7 },
  { atomicNumber: 110, symbol: "Ds", name: "Darmstadtium", nameTH: "ดาร์มสตัดเทียม", atomicMass: 281, category: "transition-metal", group: 10, period: 7 },
  { atomicNumber: 111, symbol: "Rg", name: "Roentgenium", nameTH: "เรินต์เกเนียม", atomicMass: 282, category: "transition-metal", group: 11, period: 7 },
  { atomicNumber: 112, symbol: "Cn", name: "Copernicium", nameTH: "โคเปอร์นิเซียม", atomicMass: 285, category: "transition-metal", group: 12, period: 7 },
  { atomicNumber: 113, symbol: "Nh", name: "Nihonium", nameTH: "นิโฮเนียม", atomicMass: 286, category: "post-transition", group: 13, period: 7 },
  { atomicNumber: 114, symbol: "Fl", name: "Flerovium", nameTH: "ฟลีโรเวียม", atomicMass: 289, category: "post-transition", group: 14, period: 7 },
  { atomicNumber: 115, symbol: "Mc", name: "Moscovium", nameTH: "มอสโกเวียม", atomicMass: 290, category: "post-transition", group: 15, period: 7 },
  { atomicNumber: 116, symbol: "Lv", name: "Livermorium", nameTH: "ลิเวอร์มอเรียม", atomicMass: 293, category: "post-transition", group: 16, period: 7 },
  { atomicNumber: 117, symbol: "Ts", name: "Tennessine", nameTH: "เทนเนสซีน", atomicMass: 294, category: "halogen", group: 17, period: 7 },
  { atomicNumber: 118, symbol: "Og", name: "Oganesson", nameTH: "โอกาเนสสัน", atomicMass: 294, category: "noble-gas", group: 18, period: 7 },
];

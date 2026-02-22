// Periodic table data for the mini periodic table component.
// Frontend-only — not fetched from the backend.
// Only the first 36 elements are included (matches the quiz scope).

export interface Element {
  atomicNumber: number;
  symbol: string;
  name: string;
  nameTH: string;
  atomicMass: number;
  category: ElementCategory;
  group: number; // column (1-18)
  period: number; // row (1-4)
}

export type ElementCategory =
  | "nonmetal"
  | "noble-gas"
  | "alkali-metal"
  | "alkaline-earth"
  | "metalloid"
  | "halogen"
  | "transition-metal"
  | "post-transition";

export const CATEGORY_COLORS: Record<ElementCategory, string> = {
  "nonmetal": "bg-green-600/80 border-green-500/40",
  "noble-gas": "bg-purple-600/80 border-purple-500/40",
  "alkali-metal": "bg-red-600/80 border-red-500/40",
  "alkaline-earth": "bg-orange-600/80 border-orange-500/40",
  "metalloid": "bg-teal-600/80 border-teal-500/40",
  "halogen": "bg-cyan-600/80 border-cyan-500/40",
  "transition-metal": "bg-blue-600/80 border-blue-500/40",
  "post-transition": "bg-indigo-600/80 border-indigo-500/40",
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
};

// Elements 1–36 (H through Kr) positioned in the standard periodic table grid
export const ELEMENTS: Element[] = [
  // Period 1
  { atomicNumber: 1, symbol: "H", name: "Hydrogen", nameTH: "ไฮโดรเจน", atomicMass: 1.008, category: "nonmetal", group: 1, period: 1 },
  { atomicNumber: 2, symbol: "He", name: "Helium", nameTH: "ฮีเลียม", atomicMass: 4.003, category: "noble-gas", group: 18, period: 1 },
  // Period 2
  { atomicNumber: 3, symbol: "Li", name: "Lithium", nameTH: "ลิเทียม", atomicMass: 6.941, category: "alkali-metal", group: 1, period: 2 },
  { atomicNumber: 4, symbol: "Be", name: "Beryllium", nameTH: "เบริลเลียม", atomicMass: 9.012, category: "alkaline-earth", group: 2, period: 2 },
  { atomicNumber: 5, symbol: "B", name: "Boron", nameTH: "โบรอน", atomicMass: 10.81, category: "metalloid", group: 13, period: 2 },
  { atomicNumber: 6, symbol: "C", name: "Carbon", nameTH: "คาร์บอน", atomicMass: 12.01, category: "nonmetal", group: 14, period: 2 },
  { atomicNumber: 7, symbol: "N", name: "Nitrogen", nameTH: "ไนโตรเจน", atomicMass: 14.01, category: "nonmetal", group: 15, period: 2 },
  { atomicNumber: 8, symbol: "O", name: "Oxygen", nameTH: "ออกซิเจน", atomicMass: 16.00, category: "nonmetal", group: 16, period: 2 },
  { atomicNumber: 9, symbol: "F", name: "Fluorine", nameTH: "ฟลูออรีน", atomicMass: 19.00, category: "halogen", group: 17, period: 2 },
  { atomicNumber: 10, symbol: "Ne", name: "Neon", nameTH: "นีออน", atomicMass: 20.18, category: "noble-gas", group: 18, period: 2 },
  // Period 3
  { atomicNumber: 11, symbol: "Na", name: "Sodium", nameTH: "โซเดียม", atomicMass: 22.99, category: "alkali-metal", group: 1, period: 3 },
  { atomicNumber: 12, symbol: "Mg", name: "Magnesium", nameTH: "แมกนีเซียม", atomicMass: 24.31, category: "alkaline-earth", group: 2, period: 3 },
  { atomicNumber: 13, symbol: "Al", name: "Aluminium", nameTH: "อะลูมิเนียม", atomicMass: 26.98, category: "post-transition", group: 13, period: 3 },
  { atomicNumber: 14, symbol: "Si", name: "Silicon", nameTH: "ซิลิคอน", atomicMass: 28.09, category: "metalloid", group: 14, period: 3 },
  { atomicNumber: 15, symbol: "P", name: "Phosphorus", nameTH: "ฟอสฟอรัส", atomicMass: 30.97, category: "nonmetal", group: 15, period: 3 },
  { atomicNumber: 16, symbol: "S", name: "Sulfur", nameTH: "กำมะถัน", atomicMass: 32.07, category: "nonmetal", group: 16, period: 3 },
  { atomicNumber: 17, symbol: "Cl", name: "Chlorine", nameTH: "คลอรีน", atomicMass: 35.45, category: "halogen", group: 17, period: 3 },
  { atomicNumber: 18, symbol: "Ar", name: "Argon", nameTH: "อาร์กอน", atomicMass: 39.95, category: "noble-gas", group: 18, period: 3 },
  // Period 4
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
];

package services

import "github.com/takumi/fastchem/internal/models"

// Elements 1–30 for atomic structure questions.
//
// Isotopes lists mass numbers that actually occur — the stable ones, plus the
// few radioisotopes taught at this level (H-3, C-14, K-40). Neutron questions
// draw from this list rather than offsetting the atomic number by a random
// amount, which used to invent nuclides like hydrogen-7.
var elements = []models.Element{
	{Symbol: "H", Name: "Hydrogen", NameTH: "ไฮโดรเจน", AtomicNumber: 1, MolarMass: 1.008, Isotopes: []int{1, 2, 3}},
	{Symbol: "He", Name: "Helium", NameTH: "ฮีเลียม", AtomicNumber: 2, MolarMass: 4.003, Isotopes: []int{3, 4}},
	{Symbol: "Li", Name: "Lithium", NameTH: "ลิเทียม", AtomicNumber: 3, MolarMass: 6.941, Isotopes: []int{6, 7}},
	{Symbol: "Be", Name: "Beryllium", NameTH: "เบริลเลียม", AtomicNumber: 4, MolarMass: 9.012, Isotopes: []int{9}},
	{Symbol: "B", Name: "Boron", NameTH: "โบรอน", AtomicNumber: 5, MolarMass: 10.81, Isotopes: []int{10, 11}},
	{Symbol: "C", Name: "Carbon", NameTH: "คาร์บอน", AtomicNumber: 6, MolarMass: 12.01, Isotopes: []int{12, 13, 14}},
	{Symbol: "N", Name: "Nitrogen", NameTH: "ไนโตรเจน", AtomicNumber: 7, MolarMass: 14.01, Isotopes: []int{14, 15}},
	{Symbol: "O", Name: "Oxygen", NameTH: "ออกซิเจน", AtomicNumber: 8, MolarMass: 16.00, Isotopes: []int{16, 17, 18}},
	{Symbol: "F", Name: "Fluorine", NameTH: "ฟลูออรีน", AtomicNumber: 9, MolarMass: 19.00, Isotopes: []int{19}},
	{Symbol: "Ne", Name: "Neon", NameTH: "นีออน", AtomicNumber: 10, MolarMass: 20.18, Isotopes: []int{20, 21, 22}},
	{Symbol: "Na", Name: "Sodium", NameTH: "โซเดียม", AtomicNumber: 11, MolarMass: 22.99, Isotopes: []int{23}},
	{Symbol: "Mg", Name: "Magnesium", NameTH: "แมกนีเซียม", AtomicNumber: 12, MolarMass: 24.31, Isotopes: []int{24, 25, 26}},
	{Symbol: "Al", Name: "Aluminium", NameTH: "อะลูมิเนียม", AtomicNumber: 13, MolarMass: 26.98, Isotopes: []int{27}},
	{Symbol: "Si", Name: "Silicon", NameTH: "ซิลิคอน", AtomicNumber: 14, MolarMass: 28.09, Isotopes: []int{28, 29, 30}},
	{Symbol: "P", Name: "Phosphorus", NameTH: "ฟอสฟอรัส", AtomicNumber: 15, MolarMass: 30.97, Isotopes: []int{31}},
	{Symbol: "S", Name: "Sulfur", NameTH: "กำมะถัน", AtomicNumber: 16, MolarMass: 32.07, Isotopes: []int{32, 33, 34, 36}},
	{Symbol: "Cl", Name: "Chlorine", NameTH: "คลอรีน", AtomicNumber: 17, MolarMass: 35.45, Isotopes: []int{35, 37}},
	{Symbol: "Ar", Name: "Argon", NameTH: "อาร์กอน", AtomicNumber: 18, MolarMass: 39.95, Isotopes: []int{36, 38, 40}},
	{Symbol: "K", Name: "Potassium", NameTH: "โพแทสเซียม", AtomicNumber: 19, MolarMass: 39.10, Isotopes: []int{39, 40, 41}},
	{Symbol: "Ca", Name: "Calcium", NameTH: "แคลเซียม", AtomicNumber: 20, MolarMass: 40.08, Isotopes: []int{40, 42, 43, 44, 46, 48}},
	{Symbol: "Sc", Name: "Scandium", NameTH: "สแกนเดียม", AtomicNumber: 21, MolarMass: 44.96, Isotopes: []int{45}},
	{Symbol: "Ti", Name: "Titanium", NameTH: "ไทเทเนียม", AtomicNumber: 22, MolarMass: 47.87, Isotopes: []int{46, 47, 48, 49, 50}},
	{Symbol: "V", Name: "Vanadium", NameTH: "วาเนเดียม", AtomicNumber: 23, MolarMass: 50.94, Isotopes: []int{50, 51}},
	{Symbol: "Cr", Name: "Chromium", NameTH: "โครเมียม", AtomicNumber: 24, MolarMass: 52.00, Isotopes: []int{50, 52, 53, 54}},
	{Symbol: "Mn", Name: "Manganese", NameTH: "แมงกานีส", AtomicNumber: 25, MolarMass: 54.94, Isotopes: []int{55}},
	{Symbol: "Fe", Name: "Iron", NameTH: "เหล็ก", AtomicNumber: 26, MolarMass: 55.85, Isotopes: []int{54, 56, 57, 58}},
	{Symbol: "Co", Name: "Cobalt", NameTH: "โคบอลต์", AtomicNumber: 27, MolarMass: 58.93, Isotopes: []int{59}},
	{Symbol: "Ni", Name: "Nickel", NameTH: "นิกเกิล", AtomicNumber: 28, MolarMass: 58.69, Isotopes: []int{58, 60, 61, 62, 64}},
	{Symbol: "Cu", Name: "Copper", NameTH: "ทองแดง", AtomicNumber: 29, MolarMass: 63.55, Isotopes: []int{63, 65}},
	{Symbol: "Zn", Name: "Zinc", NameTH: "สังกะสี", AtomicNumber: 30, MolarMass: 65.38, Isotopes: []int{64, 66, 67, 68, 70}},
}

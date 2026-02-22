package services

import "github.com/takumi/fastchem/internal/models"

// Elements 1–30 for atomic structure questions
var elements = []models.Element{
	{Symbol: "H", Name: "Hydrogen", NameTH: "ไฮโดรเจน", AtomicNumber: 1, MolarMass: 1.008},
	{Symbol: "He", Name: "Helium", NameTH: "ฮีเลียม", AtomicNumber: 2, MolarMass: 4.003},
	{Symbol: "Li", Name: "Lithium", NameTH: "ลิเทียม", AtomicNumber: 3, MolarMass: 6.941},
	{Symbol: "Be", Name: "Beryllium", NameTH: "เบริลเลียม", AtomicNumber: 4, MolarMass: 9.012},
	{Symbol: "B", Name: "Boron", NameTH: "โบรอน", AtomicNumber: 5, MolarMass: 10.81},
	{Symbol: "C", Name: "Carbon", NameTH: "คาร์บอน", AtomicNumber: 6, MolarMass: 12.01},
	{Symbol: "N", Name: "Nitrogen", NameTH: "ไนโตรเจน", AtomicNumber: 7, MolarMass: 14.01},
	{Symbol: "O", Name: "Oxygen", NameTH: "ออกซิเจน", AtomicNumber: 8, MolarMass: 16.00},
	{Symbol: "F", Name: "Fluorine", NameTH: "ฟลูออรีน", AtomicNumber: 9, MolarMass: 19.00},
	{Symbol: "Ne", Name: "Neon", NameTH: "นีออน", AtomicNumber: 10, MolarMass: 20.18},
	{Symbol: "Na", Name: "Sodium", NameTH: "โซเดียม", AtomicNumber: 11, MolarMass: 22.99},
	{Symbol: "Mg", Name: "Magnesium", NameTH: "แมกนีเซียม", AtomicNumber: 12, MolarMass: 24.31},
	{Symbol: "Al", Name: "Aluminium", NameTH: "อะลูมิเนียม", AtomicNumber: 13, MolarMass: 26.98},
	{Symbol: "Si", Name: "Silicon", NameTH: "ซิลิคอน", AtomicNumber: 14, MolarMass: 28.09},
	{Symbol: "P", Name: "Phosphorus", NameTH: "ฟอสฟอรัส", AtomicNumber: 15, MolarMass: 30.97},
	{Symbol: "S", Name: "Sulfur", NameTH: "กำมะถัน", AtomicNumber: 16, MolarMass: 32.07},
	{Symbol: "Cl", Name: "Chlorine", NameTH: "คลอรีน", AtomicNumber: 17, MolarMass: 35.45},
	{Symbol: "Ar", Name: "Argon", NameTH: "อาร์กอน", AtomicNumber: 18, MolarMass: 39.95},
	{Symbol: "K", Name: "Potassium", NameTH: "โพแทสเซียม", AtomicNumber: 19, MolarMass: 39.10},
	{Symbol: "Ca", Name: "Calcium", NameTH: "แคลเซียม", AtomicNumber: 20, MolarMass: 40.08},
	{Symbol: "Sc", Name: "Scandium", NameTH: "สแกนเดียม", AtomicNumber: 21, MolarMass: 44.96},
	{Symbol: "Ti", Name: "Titanium", NameTH: "ไทเทเนียม", AtomicNumber: 22, MolarMass: 47.87},
	{Symbol: "V", Name: "Vanadium", NameTH: "วาเนเดียม", AtomicNumber: 23, MolarMass: 50.94},
	{Symbol: "Cr", Name: "Chromium", NameTH: "โครเมียม", AtomicNumber: 24, MolarMass: 52.00},
	{Symbol: "Mn", Name: "Manganese", NameTH: "แมงกานีส", AtomicNumber: 25, MolarMass: 54.94},
	{Symbol: "Fe", Name: "Iron", NameTH: "เหล็ก", AtomicNumber: 26, MolarMass: 55.85},
	{Symbol: "Co", Name: "Cobalt", NameTH: "โคบอลต์", AtomicNumber: 27, MolarMass: 58.93},
	{Symbol: "Ni", Name: "Nickel", NameTH: "นิกเกิล", AtomicNumber: 28, MolarMass: 58.69},
	{Symbol: "Cu", Name: "Copper", NameTH: "ทองแดง", AtomicNumber: 29, MolarMass: 63.55},
	{Symbol: "Zn", Name: "Zinc", NameTH: "สังกะสี", AtomicNumber: 30, MolarMass: 65.38},
}

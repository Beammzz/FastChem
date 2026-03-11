package services

import "github.com/takumi/fastchem/internal/models"

// compounds is used for oxidation-number questions.
// State field is also filled so the same slice can power state-of-matter questions.
var compounds = []models.Compound{
	// ─── Binary ionic / simple salts ───────────────────────────
	{
		Formula: "NaCl", NameTH: "โซเดียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	{
		Formula: "KCl", NameTH: "โพแทสเซียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	{
		Formula: "MgCl₂", NameTH: "แมกนีเซียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Mg", OxidationNumber: +2, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 2},
		},
	},
	{
		Formula: "CaCl₂", NameTH: "แคลเซียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 2},
		},
	},
	{
		Formula: "BaCl₂", NameTH: "แบเรียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ba", OxidationNumber: +2, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 2},
		},
	},
	{
		Formula: "FeCl₂", NameTH: "เหล็ก(II)คลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +2, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 2},
		},
	},
	{
		Formula: "FeCl₃", NameTH: "เหล็ก(III)คลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +3, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 3},
		},
	},
	{
		Formula: "AlCl₃", NameTH: "อะลูมิเนียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Al", OxidationNumber: +3, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 3},
		},
	},
	{
		Formula: "CuCl₂", NameTH: "คอปเปอร์(II)คลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Cu", OxidationNumber: +2, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 2},
		},
	},
	{
		Formula: "AgCl", NameTH: "ซิลเวอร์คลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ag", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	// ─── Oxides ────────────────────────────────────────────────
	{
		Formula: "H₂O", NameTH: "น้ำ", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "CO", NameTH: "คาร์บอนมอนอกไซด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "MgO", NameTH: "แมกนีเซียมออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Mg", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "CaO", NameTH: "แคลเซียมออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "Al₂O₃", NameTH: "อะลูมิเนียมออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Al", OxidationNumber: +3, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "Fe₂O₃", NameTH: "เหล็ก(III)ออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +3, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "FeO", NameTH: "เหล็ก(II)ออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "Fe₃O₄", NameTH: "เหล็กออกไซด์ผสม", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "ZnO", NameTH: "สังกะสีออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Zn", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "CuO", NameTH: "คอปเปอร์(II)ออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Cu", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "PbO₂", NameTH: "ตะกั่ว(IV)ออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Pb", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "MnO₂", NameTH: "แมงกานีส(IV)ออกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Mn", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "SO₂", NameTH: "ซัลเฟอร์ไดออกไซด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "S", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "SO₃", NameTH: "ซัลเฟอร์ไตรออกไซด์", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "NO₂", NameTH: "ไนโตรเจนไดออกไซด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "N₂O", NameTH: "ไดไนโตรเจนมอนอกไซด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: +1, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "Cl₂O₇", NameTH: "ไดคลอรีนเฮปตะออกไซด์", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "Cl", OxidationNumber: +7, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 7},
		},
	},
	// ─── Acids ─────────────────────────────────────────────────
	{
		Formula: "H₂SO₄", NameTH: "กรดซัลฟิวริก", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "HNO₃", NameTH: "กรดไนทริก", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "N", OxidationNumber: +5, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "H₃PO₄", NameTH: "กรดฟอสฟอริก", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 3},
			{Symbol: "P", OxidationNumber: +5, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "HClO₄", NameTH: "กรดเปอร์คลอริก", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: +7, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "HClO₃", NameTH: "กรดคลอริก", State: "aq",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: +5, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "HClO₂", NameTH: "กรดคลอรัส", State: "aq",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: +3, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "HClO", NameTH: "กรดไฮโปคลอรัส", State: "aq",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: +1, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "HCl", NameTH: "ไฮโดรเจนคลอไรด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	{
		Formula: "H₂S", NameTH: "ไฮโดรเจนซัลไฟด์", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "H₂O₂", NameTH: "ไฮโดรเจนเปอร์ออกไซด์", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "O", OxidationNumber: -1, Count: 2},
		},
	},
	// ─── Bases / hydroxides ────────────────────────────────────
	{
		Formula: "NaOH", NameTH: "โซเดียมไฮดรอกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 1},
		},
	},
	{
		Formula: "KOH", NameTH: "โพแทสเซียมไฮดรอกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 1},
		},
	},
	{
		Formula: "Ca(OH)₂", NameTH: "แคลเซียมไฮดรอกไซด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
			{Symbol: "H", OxidationNumber: +1, Count: 2},
		},
	},
	{
		Formula: "NH₃", NameTH: "แอมโมเนีย", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: -3, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 3},
		},
	},
	// ─── Sulfates / nitrates / carbonates / phosphates ─────────
	{
		Formula: "Na₂SO₄", NameTH: "โซเดียมซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "K₂SO₄", NameTH: "โพแทสเซียมซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "CuSO₄", NameTH: "คอปเปอร์(II)ซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Cu", OxidationNumber: +2, Count: 1},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "ZnSO₄", NameTH: "สังกะสีซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Zn", OxidationNumber: +2, Count: 1},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "MgSO₄", NameTH: "แมกนีเซียมซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Mg", OxidationNumber: +2, Count: 1},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "BaSO₄", NameTH: "แบเรียมซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ba", OxidationNumber: +2, Count: 1},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "CaCO₃", NameTH: "แคลเซียมคาร์บอเนต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 1},
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "Na₂CO₃", NameTH: "โซเดียมคาร์บอเนต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 2},
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "NaHCO₃", NameTH: "โซเดียมไบคาร์บอเนต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "AgNO₃", NameTH: "ซิลเวอร์ไนเตรต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ag", OxidationNumber: +1, Count: 1},
			{Symbol: "N", OxidationNumber: +5, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "Ca₃(PO₄)₂", NameTH: "แคลเซียมฟอสเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 3},
			{Symbol: "P", OxidationNumber: +5, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 8},
		},
	},
	// ─── Permanganate / dichromate ──────────────────────────────
	{
		Formula: "KMnO₄", NameTH: "โพแทสเซียมเปอร์แมงกาเนต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 1},
			{Symbol: "Mn", OxidationNumber: +7, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "K₂Cr₂O₇", NameTH: "โพแทสเซียมไดโครเมต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 2},
			{Symbol: "Cr", OxidationNumber: +6, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 7},
		},
	},
	{
		Formula: "K₂CrO₄", NameTH: "โพแทสเซียมโครเมต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 2},
			{Symbol: "Cr", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	// ─── Nitrogen / ammonium ────────────────────────────────────
	{
		Formula: "NH₄Cl", NameTH: "แอมโมเนียมคลอไรด์", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: -3, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 4},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	{
		Formula: "(NH₄)₂SO₄", NameTH: "แอมโมเนียมซัลเฟต", State: "s",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: -3, Count: 2},
			{Symbol: "H", OxidationNumber: +1, Count: 8},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	// ─── Organic / halogens ─────────────────────────────────────
	{
		Formula: "CCl₄", NameTH: "คาร์บอนเตตระคลอไรด์", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 4},
		},
	},
	{
		Formula: "CHCl₃", NameTH: "คลอโรฟอร์ม", State: "l",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: +2, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 3},
		},
	},
	{
		Formula: "CH₄", NameTH: "มีเทน", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: -4, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 4},
		},
	},
	{
		Formula: "C₂H₂", NameTH: "อะเซทิลีน", State: "g",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: -1, Count: 2},
			{Symbol: "H", OxidationNumber: +1, Count: 2},
		},
	},
}

// stateCompounds is a curated subset for state-of-matter questions.
// Each compound has an unambiguous standard state (25 °C, 1 atm): s, l, g, or aq.
var stateCompounds = []models.Compound{
	// ─── Solid (s) ──────────────────────────────────────────────
	{Formula: "NaCl", NameTH: "โซเดียมคลอไรด์", State: "s"},
	{Formula: "KCl", NameTH: "โพแทสเซียมคลอไรด์", State: "s"},
	{Formula: "MgCl₂", NameTH: "แมกนีเซียมคลอไรด์", State: "s"},
	{Formula: "CaCl₂", NameTH: "แคลเซียมคลอไรด์", State: "s"},
	{Formula: "BaCl₂", NameTH: "แบเรียมคลอไรด์", State: "s"},
	{Formula: "AlCl₃", NameTH: "อะลูมิเนียมคลอไรด์", State: "s"},
	{Formula: "FeCl₂", NameTH: "เหล็ก(II)คลอไรด์", State: "s"},
	{Formula: "FeCl₃", NameTH: "เหล็ก(III)คลอไรด์", State: "s"},
	{Formula: "AgCl", NameTH: "ซิลเวอร์คลอไรด์", State: "s"},
	{Formula: "CaCO₃", NameTH: "แคลเซียมคาร์บอเนต", State: "s"},
	{Formula: "Na₂CO₃", NameTH: "โซเดียมคาร์บอเนต", State: "s"},
	{Formula: "NaHCO₃", NameTH: "โซเดียมไบคาร์บอเนต", State: "s"},
	{Formula: "Fe₂O₃", NameTH: "เหล็ก(III)ออกไซด์", State: "s"},
	{Formula: "FeO", NameTH: "เหล็ก(II)ออกไซด์", State: "s"},
	{Formula: "Fe₃O₄", NameTH: "เหล็กออกไซด์ผสม", State: "s"},
	{Formula: "CuO", NameTH: "คอปเปอร์(II)ออกไซด์", State: "s"},
	{Formula: "MnO₂", NameTH: "แมงกานีส(IV)ออกไซด์", State: "s"},
	{Formula: "PbO₂", NameTH: "ตะกั่ว(IV)ออกไซด์", State: "s"},
	{Formula: "KMnO₄", NameTH: "โพแทสเซียมเปอร์แมงกาเนต", State: "s"},
	{Formula: "K₂Cr₂O₇", NameTH: "โพแทสเซียมไดโครเมต", State: "s"},
	{Formula: "K₂CrO₄", NameTH: "โพแทสเซียมโครเมต", State: "s"},
	{Formula: "NaOH", NameTH: "โซเดียมไฮดรอกไซด์", State: "s"},
	{Formula: "KOH", NameTH: "โพแทสเซียมไฮดรอกไซด์", State: "s"},
	{Formula: "Ca(OH)₂", NameTH: "แคลเซียมไฮดรอกไซด์", State: "s"},
	{Formula: "MgO", NameTH: "แมกนีเซียมออกไซด์", State: "s"},
	{Formula: "CaO", NameTH: "แคลเซียมออกไซด์", State: "s"},
	{Formula: "Al₂O₃", NameTH: "อะลูมิเนียมออกไซด์", State: "s"},
	{Formula: "ZnO", NameTH: "สังกะสีออกไซด์", State: "s"},
	{Formula: "SiO₂", NameTH: "ซิลิคอนไดออกไซด์", State: "s"},
	{Formula: "BaSO₄", NameTH: "แบเรียมซัลเฟต", State: "s"},
	{Formula: "CuSO₄", NameTH: "คอปเปอร์(II)ซัลเฟต", State: "s"},
	{Formula: "ZnSO₄", NameTH: "สังกะสีซัลเฟต", State: "s"},
	{Formula: "MgSO₄", NameTH: "แมกนีเซียมซัลเฟต", State: "s"},
	{Formula: "K₂SO₄", NameTH: "โพแทสเซียมซัลเฟต", State: "s"},
	{Formula: "Na₂SO₄", NameTH: "โซเดียมซัลเฟต", State: "s"},
	{Formula: "PbSO₄", NameTH: "ตะกั่ว(II)ซัลเฟต", State: "s"},
	{Formula: "AgNO₃", NameTH: "ซิลเวอร์ไนเตรต", State: "s"},
	{Formula: "KNO₃", NameTH: "โพแทสเซียมไนเตรต", State: "s"},
	{Formula: "NaNO₃", NameTH: "โซเดียมไนเตรต", State: "s"},
	{Formula: "Pb(NO₃)₂", NameTH: "ตะกั่ว(II)ไนเตรต", State: "s"},
	{Formula: "Ca₃(PO₄)₂", NameTH: "แคลเซียมฟอสเฟต", State: "s"},
	{Formula: "H₃PO₄", NameTH: "กรดฟอสฟอริก", State: "s"},
	{Formula: "NH₄Cl", NameTH: "แอมโมเนียมคลอไรด์", State: "s"},
	{Formula: "(NH₄)₂SO₄", NameTH: "แอมโมเนียมซัลเฟต", State: "s"},
	{Formula: "I₂", NameTH: "ไอโอดีน", State: "s"},
	{Formula: "S₈", NameTH: "กำมะถัน", State: "s"},
	{Formula: "P₄", NameTH: "ฟอสฟอรัสขาว", State: "s"},
	// ─── Liquid (l) ─────────────────────────────────────────────
	{Formula: "H₂O", NameTH: "น้ำ", State: "l"},
	{Formula: "H₂SO₄", NameTH: "กรดซัลฟิวริกเข้มข้น", State: "l"},
	{Formula: "HNO₃", NameTH: "กรดไนทริกเข้มข้น", State: "l"},
	{Formula: "H₂O₂", NameTH: "ไฮโดรเจนเปอร์ออกไซด์", State: "l"},
	{Formula: "HClO₄", NameTH: "กรดเปอร์คลอริก", State: "l"},
	{Formula: "SO₃", NameTH: "ซัลเฟอร์ไตรออกไซด์", State: "l"},
	{Formula: "Br₂", NameTH: "โบรมีน", State: "l"},
	{Formula: "Hg", NameTH: "ปรอท", State: "l"},
	{Formula: "CCl₄", NameTH: "คาร์บอนเตตระคลอไรด์", State: "l"},
	{Formula: "CHCl₃", NameTH: "คลอโรฟอร์ม", State: "l"},
	{Formula: "Cl₂O₇", NameTH: "ไดคลอรีนเฮปตะออกไซด์", State: "l"},
	{Formula: "C₆H₆", NameTH: "เบนซีน", State: "l"},
	{Formula: "C₂H₅OH", NameTH: "เอทานอล", State: "l"},
	{Formula: "CH₃OH", NameTH: "เมทานอล", State: "l"},
	// ─── Gas (g) ────────────────────────────────────────────────
	{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", State: "g"},
	{Formula: "CO", NameTH: "คาร์บอนมอนอกไซด์", State: "g"},
	{Formula: "CH₄", NameTH: "มีเทน", State: "g"},
	{Formula: "C₂H₄", NameTH: "เอทิลีน", State: "g"},
	{Formula: "C₂H₂", NameTH: "อะเซทิลีน", State: "g"},
	{Formula: "NH₃", NameTH: "แอมโมเนีย", State: "g"},
	{Formula: "HCl", NameTH: "ไฮโดรเจนคลอไรด์", State: "g"},
	{Formula: "HBr", NameTH: "ไฮโดรเจนโบรไมด์", State: "g"},
	{Formula: "HI", NameTH: "ไฮโดรเจนไอโอไดด์", State: "g"},
	{Formula: "HF", NameTH: "ไฮโดรเจนฟลูออไรด์", State: "g"},
	{Formula: "H₂S", NameTH: "ไฮโดรเจนซัลไฟด์", State: "g"},
	{Formula: "SO₂", NameTH: "ซัลเฟอร์ไดออกไซด์", State: "g"},
	{Formula: "NO₂", NameTH: "ไนโตรเจนไดออกไซด์", State: "g"},
	{Formula: "NO", NameTH: "ไนตริกออกไซด์", State: "g"},
	{Formula: "N₂O", NameTH: "ไดไนโตรเจนมอนอกไซด์", State: "g"},
	{Formula: "Cl₂", NameTH: "คลอรีน", State: "g"},
	{Formula: "F₂", NameTH: "ฟลูออรีน", State: "g"},
	{Formula: "O₂", NameTH: "ออกซิเจน", State: "g"},
	{Formula: "O₃", NameTH: "โอโซน", State: "g"},
	{Formula: "N₂", NameTH: "ไนโตรเจน", State: "g"},
	{Formula: "H₂", NameTH: "ไฮโดรเจน", State: "g"},
	// ─── Aqueous (aq) ───────────────────────────────────────────
	{Formula: "NaCl(aq)", NameTH: "สารละลายโซเดียมคลอไรด์", State: "aq"},
	{Formula: "KCl(aq)", NameTH: "สารละลายโพแทสเซียมคลอไรด์", State: "aq"},
	{Formula: "MgCl₂(aq)", NameTH: "สารละลายแมกนีเซียมคลอไรด์", State: "aq"},
	{Formula: "AlCl₃(aq)", NameTH: "สารละลายอะลูมิเนียมคลอไรด์", State: "aq"},
	{Formula: "NaOH(aq)", NameTH: "สารละลายโซเดียมไฮดรอกไซด์", State: "aq"},
	{Formula: "KOH(aq)", NameTH: "สารละลายโพแทสเซียมไฮดรอกไซด์", State: "aq"},
	{Formula: "HCl(aq)", NameTH: "สารละลายกรดเกลือ", State: "aq"},
	{Formula: "HNO₃(aq)", NameTH: "สารละลายกรดไนทริก", State: "aq"},
	{Formula: "H₂SO₄(aq)", NameTH: "สารละลายกรดซัลฟิวริก", State: "aq"},
	{Formula: "CH₃COOH(aq)", NameTH: "สารละลายกรดแอซีติก", State: "aq"},
	{Formula: "CuSO₄(aq)", NameTH: "สารละลายคอปเปอร์(II)ซัลเฟต", State: "aq"},
	{Formula: "ZnSO₄(aq)", NameTH: "สารละลายสังกะสีซัลเฟต", State: "aq"},
	{Formula: "FeSO₄(aq)", NameTH: "สารละลายเหล็ก(II)ซัลเฟต", State: "aq"},
	{Formula: "FeCl₃(aq)", NameTH: "สารละลายเหล็ก(III)คลอไรด์", State: "aq"},
	{Formula: "AgNO₃(aq)", NameTH: "สารละลายซิลเวอร์ไนเตรต", State: "aq"},
	{Formula: "Pb(NO₃)₂(aq)", NameTH: "สารละลายตะกั่ว(II)ไนเตรต", State: "aq"},
	{Formula: "Na₂CO₃(aq)", NameTH: "สารละลายโซเดียมคาร์บอเนต", State: "aq"},
	{Formula: "NaHCO₃(aq)", NameTH: "สารละลายโซเดียมไบคาร์บอเนต", State: "aq"},
	{Formula: "KMnO₄(aq)", NameTH: "สารละลายโพแทสเซียมเปอร์แมงกาเนต", State: "aq"},
	{Formula: "K₂Cr₂O₇(aq)", NameTH: "สารละลายโพแทสเซียมไดโครเมต", State: "aq"},
	{Formula: "NH₃(aq)", NameTH: "สารละลายแอมโมเนีย", State: "aq"},
	{Formula: "NH₄Cl(aq)", NameTH: "สารละลายแอมโมเนียมคลอไรด์", State: "aq"},
}

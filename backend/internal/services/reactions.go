package services

// Reaction reference data for บทที่ 6 ปริมาณสัมพันธ์ and บทที่ 9 สมดุลเคมี.

// reactionSpecies is one substance in a balanced equation. Coef is its
// coefficient as written in Equation — the mole ratio a stoichiometry question
// tests is the ratio of these numbers, so the two must never drift apart.
type reactionSpecies struct {
	Formula   string
	NameTH    string
	Coef      int
	MolarMass float64 // g/mol
}

// reaction is a balanced equation with its species split by side.
type reaction struct {
	Equation  string
	Reactants []reactionSpecies
	Products  []reactionSpecies
}

var reactions = []reaction{
	{
		Equation: "2H₂ + O₂ → 2H₂O",
		Reactants: []reactionSpecies{
			{Formula: "H₂", NameTH: "ไฮโดรเจน", Coef: 2, MolarMass: 2.02},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 1, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "H₂O", NameTH: "น้ำ", Coef: 2, MolarMass: 18.02},
		},
	},
	{
		Equation: "N₂ + 3H₂ → 2NH₃",
		Reactants: []reactionSpecies{
			{Formula: "N₂", NameTH: "ไนโตรเจน", Coef: 1, MolarMass: 28.02},
			{Formula: "H₂", NameTH: "ไฮโดรเจน", Coef: 3, MolarMass: 2.02},
		},
		Products: []reactionSpecies{
			{Formula: "NH₃", NameTH: "แอมโมเนีย", Coef: 2, MolarMass: 17.03},
		},
	},
	{
		Equation: "CH₄ + 2O₂ → CO₂ + 2H₂O",
		Reactants: []reactionSpecies{
			{Formula: "CH₄", NameTH: "มีเทน", Coef: 1, MolarMass: 16.04},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 2, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Coef: 1, MolarMass: 44.01},
			{Formula: "H₂O", NameTH: "น้ำ", Coef: 2, MolarMass: 18.02},
		},
	},
	{
		Equation: "2Na + Cl₂ → 2NaCl",
		Reactants: []reactionSpecies{
			{Formula: "Na", NameTH: "โซเดียม", Coef: 2, MolarMass: 22.99},
			{Formula: "Cl₂", NameTH: "คลอรีน", Coef: 1, MolarMass: 70.90},
		},
		Products: []reactionSpecies{
			{Formula: "NaCl", NameTH: "โซเดียมคลอไรด์", Coef: 2, MolarMass: 58.44},
		},
	},
	{
		Equation: "2Mg + O₂ → 2MgO",
		Reactants: []reactionSpecies{
			{Formula: "Mg", NameTH: "แมกนีเซียม", Coef: 2, MolarMass: 24.31},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 1, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "MgO", NameTH: "แมกนีเซียมออกไซด์", Coef: 2, MolarMass: 40.31},
		},
	},
	{
		Equation: "Zn + 2HCl → ZnCl₂ + H₂",
		Reactants: []reactionSpecies{
			{Formula: "Zn", NameTH: "สังกะสี", Coef: 1, MolarMass: 65.38},
			{Formula: "HCl", NameTH: "กรดไฮโดรคลอริก", Coef: 2, MolarMass: 36.46},
		},
		Products: []reactionSpecies{
			{Formula: "ZnCl₂", NameTH: "ซิงค์คลอไรด์", Coef: 1, MolarMass: 136.28},
			{Formula: "H₂", NameTH: "ไฮโดรเจน", Coef: 1, MolarMass: 2.02},
		},
	},
	{
		Equation: "4Fe + 3O₂ → 2Fe₂O₃",
		Reactants: []reactionSpecies{
			{Formula: "Fe", NameTH: "เหล็ก", Coef: 4, MolarMass: 55.85},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 3, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "Fe₂O₃", NameTH: "ไอร์ออน(III)ออกไซด์", Coef: 2, MolarMass: 159.70},
		},
	},
	{
		Equation: "2KClO₃ → 2KCl + 3O₂",
		Reactants: []reactionSpecies{
			{Formula: "KClO₃", NameTH: "โพแทสเซียมคลอเรต", Coef: 2, MolarMass: 122.55},
		},
		Products: []reactionSpecies{
			{Formula: "KCl", NameTH: "โพแทสเซียมคลอไรด์", Coef: 2, MolarMass: 74.55},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 3, MolarMass: 32.00},
		},
	},
	{
		Equation: "C₃H₈ + 5O₂ → 3CO₂ + 4H₂O",
		Reactants: []reactionSpecies{
			{Formula: "C₃H₈", NameTH: "โพรเพน", Coef: 1, MolarMass: 44.09},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 5, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Coef: 3, MolarMass: 44.01},
			{Formula: "H₂O", NameTH: "น้ำ", Coef: 4, MolarMass: 18.02},
		},
	},
	{
		Equation: "2H₂O₂ → 2H₂O + O₂",
		Reactants: []reactionSpecies{
			{Formula: "H₂O₂", NameTH: "ไฮโดรเจนเปอร์ออกไซด์", Coef: 2, MolarMass: 34.02},
		},
		Products: []reactionSpecies{
			{Formula: "H₂O", NameTH: "น้ำ", Coef: 2, MolarMass: 18.02},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 1, MolarMass: 32.00},
		},
	},
	{
		Equation: "CaCO₃ → CaO + CO₂",
		Reactants: []reactionSpecies{
			{Formula: "CaCO₃", NameTH: "แคลเซียมคาร์บอเนต", Coef: 1, MolarMass: 100.09},
		},
		Products: []reactionSpecies{
			{Formula: "CaO", NameTH: "แคลเซียมออกไซด์", Coef: 1, MolarMass: 56.08},
			{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Coef: 1, MolarMass: 44.01},
		},
	},
	{
		Equation: "4Al + 3O₂ → 2Al₂O₃",
		Reactants: []reactionSpecies{
			{Formula: "Al", NameTH: "อะลูมิเนียม", Coef: 4, MolarMass: 26.98},
			{Formula: "O₂", NameTH: "ออกซิเจน", Coef: 3, MolarMass: 32.00},
		},
		Products: []reactionSpecies{
			{Formula: "Al₂O₃", NameTH: "อะลูมิเนียมออกไซด์", Coef: 2, MolarMass: 101.96},
		},
	},
}

// twoReactantReactions indexes the reactions with exactly two reactants — the
// only ones a limiting-reagent question can be built from.
var twoReactantReactions = filterTwoReactant(reactions)

func filterTwoReactant(all []reaction) []reaction {
	out := make([]reaction, 0, len(all))
	for _, r := range all {
		if len(r.Reactants) == 2 {
			out = append(out, r)
		}
	}
	return out
}

// ─── Equilibrium ─────────────────────────────────────────────────

// equilibriumSpecies is one substance in an equilibrium expression. Only the
// coefficient matters: it becomes the exponent on that species' concentration.
type equilibriumSpecies struct {
	Formula string
	Coef    int
}

// equilibriumReaction is a homogeneous gas-phase equilibrium. All species are
// gases, so every one of them appears in Kc — no phase has to be dropped, and
// the question stays about the exponent rule rather than about which species
// to leave out.
type equilibriumReaction struct {
	Equation  string
	Reactants []equilibriumSpecies
	Products  []equilibriumSpecies
}

var equilibriumReactions = []equilibriumReaction{
	{
		Equation:  "N₂(g) + 3H₂(g) ⇌ 2NH₃(g)",
		Reactants: []equilibriumSpecies{{"N₂", 1}, {"H₂", 3}},
		Products:  []equilibriumSpecies{{"NH₃", 2}},
	},
	{
		Equation:  "H₂(g) + I₂(g) ⇌ 2HI(g)",
		Reactants: []equilibriumSpecies{{"H₂", 1}, {"I₂", 1}},
		Products:  []equilibriumSpecies{{"HI", 2}},
	},
	{
		Equation:  "2SO₂(g) + O₂(g) ⇌ 2SO₃(g)",
		Reactants: []equilibriumSpecies{{"SO₂", 2}, {"O₂", 1}},
		Products:  []equilibriumSpecies{{"SO₃", 2}},
	},
	{
		Equation:  "PCl₅(g) ⇌ PCl₃(g) + Cl₂(g)",
		Reactants: []equilibriumSpecies{{"PCl₅", 1}},
		Products:  []equilibriumSpecies{{"PCl₃", 1}, {"Cl₂", 1}},
	},
	{
		Equation:  "2NO₂(g) ⇌ N₂O₄(g)",
		Reactants: []equilibriumSpecies{{"NO₂", 2}},
		Products:  []equilibriumSpecies{{"N₂O₄", 1}},
	},
	{
		Equation:  "CO(g) + H₂O(g) ⇌ CO₂(g) + H₂(g)",
		Reactants: []equilibriumSpecies{{"CO", 1}, {"H₂O", 1}},
		Products:  []equilibriumSpecies{{"CO₂", 1}, {"H₂", 1}},
	},
	{
		Equation:  "2NO(g) + O₂(g) ⇌ 2NO₂(g)",
		Reactants: []equilibriumSpecies{{"NO", 2}, {"O₂", 1}},
		Products:  []equilibriumSpecies{{"NO₂", 2}},
	},
	{
		Equation:  "N₂(g) + O₂(g) ⇌ 2NO(g)",
		Reactants: []equilibriumSpecies{{"N₂", 1}, {"O₂", 1}},
		Products:  []equilibriumSpecies{{"NO", 2}},
	},
}

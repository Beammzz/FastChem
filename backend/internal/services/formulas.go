package services

// Formula reference data for บทที่ 4 โมลและสูตรเคมี — molecular mass and
// percent composition.

// atomicMasses is the g/mol table the formula questions quote. Values match
// elements.go, so a molar mass printed by a mole question and one printed by a
// formula question can never disagree.
//
// Only ever read by key: ranging over it would put map order into the seeded
// generation path.
var atomicMasses = map[string]float64{
	"H": 1.008, "He": 4.003, "Li": 6.941, "Be": 9.012, "B": 10.81,
	"C": 12.01, "N": 14.01, "O": 16.00, "F": 19.00, "Ne": 20.18,
	"Na": 22.99, "Mg": 24.31, "Al": 26.98, "Si": 28.09, "P": 30.97,
	"S": 32.07, "Cl": 35.45, "Ar": 39.95, "K": 39.10, "Ca": 40.08,
	"Cr": 52.00, "Mn": 54.94, "Fe": 55.85, "Cu": 63.55, "Zn": 65.38,
}

// formulaPart is one element in a formula, with how many atoms of it appear.
type formulaPart struct {
	Symbol string
	Count  int
}

// Mass is the contribution this part makes to the molar mass.
func (p formulaPart) Mass() float64 { return atomicMasses[p.Symbol] * float64(p.Count) }

// formulaCompound is a compound whose molar mass is derived from its parts
// rather than stored. A hand-typed molar mass that drifts from the subscripts
// makes both the answer and every distractor wrong at once.
type formulaCompound struct {
	Formula string
	NameTH  string
	Parts   []formulaPart
}

// MolarMass sums the parts, in g/mol.
func (c formulaCompound) MolarMass() float64 {
	total := 0.0
	for _, p := range c.Parts {
		total += p.Mass()
	}
	return total
}

// TotalAtoms is the number of atoms in one formula unit. Percent-composition
// questions need it for the "counted atoms instead of mass" misconception.
func (c formulaCompound) TotalAtoms() int {
	n := 0
	for _, p := range c.Parts {
		n += p.Count
	}
	return n
}

var formulaCompounds = []formulaCompound{
	{Formula: "H₂O", NameTH: "น้ำ", Parts: []formulaPart{{"H", 2}, {"O", 1}}},
	{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Parts: []formulaPart{{"C", 1}, {"O", 2}}},
	{Formula: "NaCl", NameTH: "โซเดียมคลอไรด์", Parts: []formulaPart{{"Na", 1}, {"Cl", 1}}},
	{Formula: "CaCO₃", NameTH: "แคลเซียมคาร์บอเนต", Parts: []formulaPart{{"Ca", 1}, {"C", 1}, {"O", 3}}},
	{Formula: "H₂SO₄", NameTH: "กรดซัลฟิวริก", Parts: []formulaPart{{"H", 2}, {"S", 1}, {"O", 4}}},
	{Formula: "NaOH", NameTH: "โซเดียมไฮดรอกไซด์", Parts: []formulaPart{{"Na", 1}, {"O", 1}, {"H", 1}}},
	{Formula: "NH₃", NameTH: "แอมโมเนีย", Parts: []formulaPart{{"N", 1}, {"H", 3}}},
	{Formula: "CH₄", NameTH: "มีเทน", Parts: []formulaPart{{"C", 1}, {"H", 4}}},
	{Formula: "C₆H₁₂O₆", NameTH: "กลูโคส", Parts: []formulaPart{{"C", 6}, {"H", 12}, {"O", 6}}},
	{Formula: "CaCl₂", NameTH: "แคลเซียมคลอไรด์", Parts: []formulaPart{{"Ca", 1}, {"Cl", 2}}},
	{Formula: "MgSO₄", NameTH: "แมกนีเซียมซัลเฟต", Parts: []formulaPart{{"Mg", 1}, {"S", 1}, {"O", 4}}},
	{Formula: "KMnO₄", NameTH: "โพแทสเซียมเปอร์แมงกาเนต", Parts: []formulaPart{{"K", 1}, {"Mn", 1}, {"O", 4}}},
	{Formula: "Al₂O₃", NameTH: "อะลูมิเนียมออกไซด์", Parts: []formulaPart{{"Al", 2}, {"O", 3}}},
	{Formula: "Fe₂O₃", NameTH: "ไอร์ออน(III)ออกไซด์", Parts: []formulaPart{{"Fe", 2}, {"O", 3}}},
	{Formula: "CuSO₄", NameTH: "คอปเปอร์(II)ซัลเฟต", Parts: []formulaPart{{"Cu", 1}, {"S", 1}, {"O", 4}}},
	{Formula: "HNO₃", NameTH: "กรดไนตริก", Parts: []formulaPart{{"H", 1}, {"N", 1}, {"O", 3}}},
	{Formula: "K₂Cr₂O₇", NameTH: "โพแทสเซียมไดโครเมต", Parts: []formulaPart{{"K", 2}, {"Cr", 2}, {"O", 7}}},
	{Formula: "C₂H₅OH", NameTH: "เอทานอล", Parts: []formulaPart{{"C", 2}, {"H", 6}, {"O", 1}}},
	{Formula: "CH₃COOH", NameTH: "กรดแอซีติก", Parts: []formulaPart{{"C", 2}, {"H", 4}, {"O", 2}}},
	{Formula: "(NH₄)₂SO₄", NameTH: "แอมโมเนียมซัลเฟต", Parts: []formulaPart{{"N", 2}, {"H", 8}, {"S", 1}, {"O", 4}}},
	{Formula: "Na₂CO₃", NameTH: "โซเดียมคาร์บอเนต", Parts: []formulaPart{{"Na", 2}, {"C", 1}, {"O", 3}}},
	{Formula: "ZnCl₂", NameTH: "ซิงค์คลอไรด์", Parts: []formulaPart{{"Zn", 1}, {"Cl", 2}}},
}

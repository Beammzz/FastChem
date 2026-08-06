package services

// Acid–base and electrochemistry reference data for บทที่ 10 กรด–เบส and
// บทที่ 11 เคมีไฟฟ้า.

// kWater is the ion-product constant of water at 25 °C, which fixes
// pH + pOH = 14.00.
const kWater = 1.0e-14

// weakAcid carries the Ka a weak-acid pH question needs. Values are the 25 °C
// figures the IPST book tabulates.
type weakAcid struct {
	NameTH  string
	Formula string
	Ka      float64
}

var weakAcids = []weakAcid{
	{NameTH: "กรดแอซีติก", Formula: "CH₃COOH", Ka: 1.8e-5},
	{NameTH: "กรดเมทาโนอิก", Formula: "HCOOH", Ka: 1.8e-4},
	{NameTH: "กรดไฮโดรฟลูออริก", Formula: "HF", Ka: 6.8e-4},
	{NameTH: "กรดไฮโดรไซยานิก", Formula: "HCN", Ka: 6.2e-10},
	{NameTH: "กรดไนตรัส", Formula: "HNO₂", Ka: 4.5e-4},
	{NameTH: "กรดไฮโปคลอรัส", Formula: "HClO", Ka: 3.0e-8},
	{NameTH: "กรดเบนโซอิก", Formula: "C₆H₅COOH", Ka: 6.3e-5},
	{NameTH: "กรดโพรพาโนอิก", Formula: "CH₃CH₂COOH", Ka: 1.3e-5},
}

// strongAcid dissociates completely, so [H₃O⁺] follows straight from the
// concentration and the number of ionisable protons.
type strongAcid struct {
	NameTH  string
	Formula string
	Protons int
}

var strongAcids = []strongAcid{
	{NameTH: "กรดไฮโดรคลอริก", Formula: "HCl", Protons: 1},
	{NameTH: "กรดไนตริก", Formula: "HNO₃", Protons: 1},
	{NameTH: "กรดไฮโดรโบรมิก", Formula: "HBr", Protons: 1},
	{NameTH: "กรดเปอร์คลอริก", Formula: "HClO₄", Protons: 1},
	{NameTH: "กรดซัลฟิวริก", Formula: "H₂SO₄", Protons: 2},
}

// halfCell is a standard reduction half-reaction and its E° at 25 °C, 1 M,
// 1 atm, against the standard hydrogen electrode.
//
// Oxidised and Reduced name the two sides of Reaction. Which side a question
// offers depends on what it asks: an oxidising agent is the species that gets
// reduced, and a reducing agent is the species that gets oxidised, so a
// question about agents cannot be answered from the half-reaction string alone.
type halfCell struct {
	Reaction string
	Oxidised string // the species on the left, which acts as the oxidising agent
	Reduced  string // the species on the right, which acts as the reducing agent
	E0       float64
}

// standardPotentials runs from the strongest reducing agent up to the
// strongest oxidising agent, matching the direction the IPST table is printed
// in. Order is not load-bearing for correctness — the topic reads by index —
// but it keeps the table readable next to the book.
var standardPotentials = []halfCell{
	{Reaction: "Li⁺ + e⁻ → Li", Oxidised: "Li⁺", Reduced: "Li", E0: -3.04},
	{Reaction: "K⁺ + e⁻ → K", Oxidised: "K⁺", Reduced: "K", E0: -2.93},
	{Reaction: "Ca²⁺ + 2e⁻ → Ca", Oxidised: "Ca²⁺", Reduced: "Ca", E0: -2.87},
	{Reaction: "Na⁺ + e⁻ → Na", Oxidised: "Na⁺", Reduced: "Na", E0: -2.71},
	{Reaction: "Mg²⁺ + 2e⁻ → Mg", Oxidised: "Mg²⁺", Reduced: "Mg", E0: -2.37},
	{Reaction: "Al³⁺ + 3e⁻ → Al", Oxidised: "Al³⁺", Reduced: "Al", E0: -1.66},
	{Reaction: "Zn²⁺ + 2e⁻ → Zn", Oxidised: "Zn²⁺", Reduced: "Zn", E0: -0.76},
	{Reaction: "Fe²⁺ + 2e⁻ → Fe", Oxidised: "Fe²⁺", Reduced: "Fe", E0: -0.44},
	{Reaction: "Ni²⁺ + 2e⁻ → Ni", Oxidised: "Ni²⁺", Reduced: "Ni", E0: -0.25},
	{Reaction: "Sn²⁺ + 2e⁻ → Sn", Oxidised: "Sn²⁺", Reduced: "Sn", E0: -0.14},
	{Reaction: "Pb²⁺ + 2e⁻ → Pb", Oxidised: "Pb²⁺", Reduced: "Pb", E0: -0.13},
	{Reaction: "2H⁺ + 2e⁻ → H₂", Oxidised: "H⁺", Reduced: "H₂", E0: 0.00},
	{Reaction: "Cu²⁺ + 2e⁻ → Cu", Oxidised: "Cu²⁺", Reduced: "Cu", E0: 0.34},
	{Reaction: "I₂ + 2e⁻ → 2I⁻", Oxidised: "I₂", Reduced: "I⁻", E0: 0.54},
	{Reaction: "Ag⁺ + e⁻ → Ag", Oxidised: "Ag⁺", Reduced: "Ag", E0: 0.80},
	{Reaction: "Br₂ + 2e⁻ → 2Br⁻", Oxidised: "Br₂", Reduced: "Br⁻", E0: 1.07},
	{Reaction: "Cl₂ + 2e⁻ → 2Cl⁻", Oxidised: "Cl₂", Reduced: "Cl⁻", E0: 1.36},
	{Reaction: "F₂ + 2e⁻ → 2F⁻", Oxidised: "F₂", Reduced: "F⁻", E0: 2.87},
}

package services

// Bonding reference data for บทที่ 3 พันธะเคมี.

// bondClasses is the closed vocabulary of the bond-type question. Every
// substance below carries exactly one of these labels, so the three wrong
// options are always the three other real categories rather than noise.
var bondClasses = []string{
	"สารประกอบไอออนิก",
	"โมเลกุลโคเวเลนต์มีขั้ว",
	"โมเลกุลโคเวเลนต์ไม่มีขั้ว",
	"พันธะโลหะ",
}

// bondSubstance classifies a substance by the bonding that holds it together.
//
// The classification is of the whole molecule, not of an individual bond —
// CO₂ and CCl₄ have polar bonds but cancel to non-polar molecules, and asking
// about the molecule is what keeps that distinction teachable.
type bondSubstance struct {
	Formula string
	NameTH  string
	Class   string
}

var bondSubstances = []bondSubstance{
	// Ionic
	{Formula: "NaCl", NameTH: "โซเดียมคลอไรด์", Class: "สารประกอบไอออนิก"},
	{Formula: "KBr", NameTH: "โพแทสเซียมโบรไมด์", Class: "สารประกอบไอออนิก"},
	{Formula: "MgO", NameTH: "แมกนีเซียมออกไซด์", Class: "สารประกอบไอออนิก"},
	{Formula: "CaCl₂", NameTH: "แคลเซียมคลอไรด์", Class: "สารประกอบไอออนิก"},
	{Formula: "LiF", NameTH: "ลิเทียมฟลูออไรด์", Class: "สารประกอบไอออนิก"},
	{Formula: "K₂S", NameTH: "โพแทสเซียมซัลไฟด์", Class: "สารประกอบไอออนิก"},
	{Formula: "Al₂O₃", NameTH: "อะลูมิเนียมออกไซด์", Class: "สารประกอบไอออนิก"},

	// Polar covalent molecules
	{Formula: "HCl", NameTH: "ไฮโดรเจนคลอไรด์", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "H₂O", NameTH: "น้ำ", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "NH₃", NameTH: "แอมโมเนีย", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "HF", NameTH: "ไฮโดรเจนฟลูออไรด์", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "H₂S", NameTH: "ไฮโดรเจนซัลไฟด์", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "SO₂", NameTH: "ซัลเฟอร์ไดออกไซด์", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},
	{Formula: "PCl₃", NameTH: "ฟอสฟอรัสไตรคลอไรด์", Class: "โมเลกุลโคเวเลนต์มีขั้ว"},

	// Non-polar covalent molecules
	{Formula: "O₂", NameTH: "ออกซิเจน", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "N₂", NameTH: "ไนโตรเจน", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "Cl₂", NameTH: "คลอรีน", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "CH₄", NameTH: "มีเทน", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "CCl₄", NameTH: "คาร์บอนเตตระคลอไรด์", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},
	{Formula: "BF₃", NameTH: "โบรอนไตรฟลูออไรด์", Class: "โมเลกุลโคเวเลนต์ไม่มีขั้ว"},

	// Metallic
	{Formula: "Cu(s)", NameTH: "ทองแดง", Class: "พันธะโลหะ"},
	{Formula: "Fe(s)", NameTH: "เหล็ก", Class: "พันธะโลหะ"},
	{Formula: "Al(s)", NameTH: "อะลูมิเนียม", Class: "พันธะโลหะ"},
	{Formula: "Zn(s)", NameTH: "สังกะสี", Class: "พันธะโลหะ"},
	{Formula: "Na(s)", NameTH: "โซเดียม", Class: "พันธะโลหะ"},
	{Formula: "Ag(s)", NameTH: "เงิน", Class: "พันธะโลหะ"},
}

// molecularShapes is the closed vocabulary of the VSEPR question.
var molecularShapes = []string{
	"เส้นตรง",
	"มุมงอ",
	"สามเหลี่ยมแบนราบ",
	"ทรงสี่หน้า",
	"พีระมิดฐานสามเหลี่ยม",
	"พีระมิดคู่ฐานสามเหลี่ยม",
	"ทรงแปดหน้า",
}

// shapeMolecule is one molecule with the shape VSEPR predicts for it.
// Bonded and Lone are the electron-pair counts around the central atom, which
// the question quotes so the answer is derivable rather than memorised.
type shapeMolecule struct {
	Formula string
	NameTH  string
	Central string
	Bonded  int
	Lone    int
	Shape   string
}

var shapeMolecules = []shapeMolecule{
	{Formula: "BeCl₂", NameTH: "เบริลเลียมคลอไรด์", Central: "Be", Bonded: 2, Lone: 0, Shape: "เส้นตรง"},
	{Formula: "CO₂", NameTH: "คาร์บอนไดออกไซด์", Central: "C", Bonded: 2, Lone: 0, Shape: "เส้นตรง"},
	{Formula: "HCN", NameTH: "ไฮโดรเจนไซยาไนด์", Central: "C", Bonded: 2, Lone: 0, Shape: "เส้นตรง"},
	{Formula: "BF₃", NameTH: "โบรอนไตรฟลูออไรด์", Central: "B", Bonded: 3, Lone: 0, Shape: "สามเหลี่ยมแบนราบ"},
	{Formula: "BCl₃", NameTH: "โบรอนไตรคลอไรด์", Central: "B", Bonded: 3, Lone: 0, Shape: "สามเหลี่ยมแบนราบ"},
	{Formula: "SO₃", NameTH: "ซัลเฟอร์ไตรออกไซด์", Central: "S", Bonded: 3, Lone: 0, Shape: "สามเหลี่ยมแบนราบ"},
	{Formula: "CH₄", NameTH: "มีเทน", Central: "C", Bonded: 4, Lone: 0, Shape: "ทรงสี่หน้า"},
	{Formula: "CCl₄", NameTH: "คาร์บอนเตตระคลอไรด์", Central: "C", Bonded: 4, Lone: 0, Shape: "ทรงสี่หน้า"},
	{Formula: "SiH₄", NameTH: "ไซเลน", Central: "Si", Bonded: 4, Lone: 0, Shape: "ทรงสี่หน้า"},
	{Formula: "NH₃", NameTH: "แอมโมเนีย", Central: "N", Bonded: 3, Lone: 1, Shape: "พีระมิดฐานสามเหลี่ยม"},
	{Formula: "PCl₃", NameTH: "ฟอสฟอรัสไตรคลอไรด์", Central: "P", Bonded: 3, Lone: 1, Shape: "พีระมิดฐานสามเหลี่ยม"},
	{Formula: "H₂O", NameTH: "น้ำ", Central: "O", Bonded: 2, Lone: 2, Shape: "มุมงอ"},
	{Formula: "H₂S", NameTH: "ไฮโดรเจนซัลไฟด์", Central: "S", Bonded: 2, Lone: 2, Shape: "มุมงอ"},
	{Formula: "SO₂", NameTH: "ซัลเฟอร์ไดออกไซด์", Central: "S", Bonded: 2, Lone: 1, Shape: "มุมงอ"},
	{Formula: "PCl₅", NameTH: "ฟอสฟอรัสเพนตะคลอไรด์", Central: "P", Bonded: 5, Lone: 0, Shape: "พีระมิดคู่ฐานสามเหลี่ยม"},
	{Formula: "SF₆", NameTH: "ซัลเฟอร์เฮกซะฟลูออไรด์", Central: "S", Bonded: 6, Lone: 0, Shape: "ทรงแปดหน้า"},
}

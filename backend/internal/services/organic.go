package services

// Organic and polymer reference data for บทที่ 12 เคมีอินทรีย์ and
// บทที่ 13 พอลิเมอร์.

// organicClasses and functionalGroups are the closed vocabularies the two
// organic question forms answer from. Every distractor is a real class or a
// real functional group, so no option can be ruled out on shape alone.
var organicClasses = []string{
	"แอลเคน", "แอลคีน", "แอลไคน์", "อะโรมาติกไฮโดรคาร์บอน",
	"แอลกอฮอล์", "อีเทอร์", "แอลดีไฮด์", "คีโตน",
	"กรดคาร์บอกซิลิก", "เอสเทอร์", "เอมีน", "เอไมด์",
}

var functionalGroups = []string{
	"–OH (ไฮดรอกซิล)", "–O– (อีเทอร์)", "–CHO (ฟอร์มิล)", "–CO– (คาร์บอนิล)",
	"–COOH (คาร์บอกซิล)", "–COO– (เอสเทอร์)", "–NH₂ (อะมิโน)", "–CONH₂ (เอไมด์)",
	"C=C (พันธะคู่)", "C≡C (พันธะสาม)", "วงเบนซีน", "พันธะเดี่ยว C–C เท่านั้น",
}

// organicCompound ties a compound to the class it belongs to and the group
// that defines that class.
type organicCompound struct {
	NameTH  string
	Formula string
	Class   string
	Group   string
}

var organicCompounds = []organicCompound{
	{NameTH: "อีเทน", Formula: "CH₃CH₃", Class: "แอลเคน", Group: "พันธะเดี่ยว C–C เท่านั้น"},
	{NameTH: "โพรเพน", Formula: "CH₃CH₂CH₃", Class: "แอลเคน", Group: "พันธะเดี่ยว C–C เท่านั้น"},
	{NameTH: "อีทีน", Formula: "CH₂=CH₂", Class: "แอลคีน", Group: "C=C (พันธะคู่)"},
	{NameTH: "โพรพีน", Formula: "CH₃CH=CH₂", Class: "แอลคีน", Group: "C=C (พันธะคู่)"},
	{NameTH: "อีไทน์", Formula: "CH≡CH", Class: "แอลไคน์", Group: "C≡C (พันธะสาม)"},
	{NameTH: "โพรไพน์", Formula: "CH₃C≡CH", Class: "แอลไคน์", Group: "C≡C (พันธะสาม)"},
	{NameTH: "เบนซีน", Formula: "C₆H₆", Class: "อะโรมาติกไฮโดรคาร์บอน", Group: "วงเบนซีน"},
	{NameTH: "โทลูอีน", Formula: "C₆H₅CH₃", Class: "อะโรมาติกไฮโดรคาร์บอน", Group: "วงเบนซีน"},
	{NameTH: "เมทานอล", Formula: "CH₃OH", Class: "แอลกอฮอล์", Group: "–OH (ไฮดรอกซิล)"},
	{NameTH: "เอทานอล", Formula: "CH₃CH₂OH", Class: "แอลกอฮอล์", Group: "–OH (ไฮดรอกซิล)"},
	{NameTH: "เมทอกซีมีเทน", Formula: "CH₃OCH₃", Class: "อีเทอร์", Group: "–O– (อีเทอร์)"},
	{NameTH: "ไดเอทิลอีเทอร์", Formula: "CH₃CH₂OCH₂CH₃", Class: "อีเทอร์", Group: "–O– (อีเทอร์)"},
	{NameTH: "เมทานาล", Formula: "HCHO", Class: "แอลดีไฮด์", Group: "–CHO (ฟอร์มิล)"},
	{NameTH: "เอทานาล", Formula: "CH₃CHO", Class: "แอลดีไฮด์", Group: "–CHO (ฟอร์มิล)"},
	{NameTH: "โพรพาโนน", Formula: "CH₃COCH₃", Class: "คีโตน", Group: "–CO– (คาร์บอนิล)"},
	{NameTH: "บิวทาโนน", Formula: "CH₃COCH₂CH₃", Class: "คีโตน", Group: "–CO– (คาร์บอนิล)"},
	{NameTH: "กรดเมทาโนอิก", Formula: "HCOOH", Class: "กรดคาร์บอกซิลิก", Group: "–COOH (คาร์บอกซิล)"},
	{NameTH: "กรดแอซีติก", Formula: "CH₃COOH", Class: "กรดคาร์บอกซิลิก", Group: "–COOH (คาร์บอกซิล)"},
	{NameTH: "เมทิลเมทาโนเอต", Formula: "HCOOCH₃", Class: "เอสเทอร์", Group: "–COO– (เอสเทอร์)"},
	{NameTH: "เอทิลเอทาโนเอต", Formula: "CH₃COOCH₂CH₃", Class: "เอสเทอร์", Group: "–COO– (เอสเทอร์)"},
	{NameTH: "เมทิลเอมีน", Formula: "CH₃NH₂", Class: "เอมีน", Group: "–NH₂ (อะมิโน)"},
	{NameTH: "เอทิลเอมีน", Formula: "CH₃CH₂NH₂", Class: "เอมีน", Group: "–NH₂ (อะมิโน)"},
	{NameTH: "เอทานาไมด์", Formula: "CH₃CONH₂", Class: "เอไมด์", Group: "–CONH₂ (เอไมด์)"},
	{NameTH: "เมทานาไมด์", Formula: "HCONH₂", Class: "เอไมด์", Group: "–CONH₂ (เอไมด์)"},
}

// Polymerisation routes, spelled the way the IPST book does.
const (
	additionPolymer     = "พอลิเมอไรเซชันแบบเติม"
	condensationPolymer = "พอลิเมอไรเซชันแบบควบแน่น"
)

// polymer pairs a polymer with the monomer it is built from.
//
// Monomers are unique across the table. A "which polymer comes from this
// monomer?" question needs a single right answer, and starch and cellulose —
// both glucose — would give it two.
type polymer struct {
	NameTH  string
	Monomer string
	Route   string
	UseTH   string
}

var polymers = []polymer{
	{NameTH: "พอลิเอทิลีน (PE)", Monomer: "เอทิลีน (CH₂=CH₂)", Route: additionPolymer, UseTH: "ถุงพลาสติก ขวดน้ำ"},
	{NameTH: "พอลิโพรพิลีน (PP)", Monomer: "โพรพิลีน (CH₂=CHCH₃)", Route: additionPolymer, UseTH: "กล่องอาหาร ถังน้ำ"},
	{NameTH: "พอลิสไตรีน (PS)", Monomer: "สไตรีน (C₆H₅CH=CH₂)", Route: additionPolymer, UseTH: "โฟม กล่องบรรจุภัณฑ์"},
	{NameTH: "พอลิไวนิลคลอไรด์ (PVC)", Monomer: "ไวนิลคลอไรด์ (CH₂=CHCl)", Route: additionPolymer, UseTH: "ท่อน้ำ สายไฟ"},
	{NameTH: "เทฟลอน (PTFE)", Monomer: "เตตระฟลูออโรเอทิลีน (CF₂=CF₂)", Route: additionPolymer, UseTH: "สารเคลือบกระทะ"},
	{NameTH: "ยางธรรมชาติ", Monomer: "ไอโซพรีน (C₅H₈)", Route: additionPolymer, UseTH: "ยางล้อ ถุงมือยาง"},
	{NameTH: "พอลิเอทิลีนเทเรฟแทเลต (PET)", Monomer: "กรดเทเรฟแทลิก + เอทิลีนไกลคอล", Route: condensationPolymer, UseTH: "ขวดน้ำดื่ม เส้นใยโพลีเอสเทอร์"},
	{NameTH: "ไนลอน-6,6", Monomer: "กรดอะดิปิก + เฮกซะเมทิลีนไดเอมีน", Route: condensationPolymer, UseTH: "เส้นใยสังเคราะห์ เชือก"},
	{NameTH: "เบกาไลต์", Monomer: "ฟีนอล + ฟอร์มาลดีไฮด์", Route: condensationPolymer, UseTH: "ด้ามจับเครื่องครัว ฉนวนไฟฟ้า"},
	{NameTH: "แป้ง", Monomer: "กลูโคส (C₆H₁₂O₆)", Route: condensationPolymer, UseTH: "พอลิเมอร์ธรรมชาติสะสมพลังงานในพืช"},
}

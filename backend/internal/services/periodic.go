package services

// Periodic-table reference data for บทที่ 2 อะตอมและสมบัติของธาตุ.
//
// Limited to Z = 1–20, the representative elements the IPST book uses when it
// teaches shell filling. Beyond calcium the 4s/3d overlap makes the simple
// 2-8-8-2 picture wrong, and a question built on it would be teaching an error.
type shellElement struct {
	Symbol   string
	NameTH   string
	Number   int
	Shells   []int  // electrons per principal energy level
	Subshell string // full subshell notation
	Group    string // representative group, IA–VIIIA
	Period   int
}

var shellElements = []shellElement{
	{Symbol: "H", NameTH: "ไฮโดรเจน", Number: 1, Shells: []int{1}, Subshell: "1s¹", Group: "IA", Period: 1},
	{Symbol: "He", NameTH: "ฮีเลียม", Number: 2, Shells: []int{2}, Subshell: "1s²", Group: "VIIIA", Period: 1},
	{Symbol: "Li", NameTH: "ลิเทียม", Number: 3, Shells: []int{2, 1}, Subshell: "1s² 2s¹", Group: "IA", Period: 2},
	{Symbol: "Be", NameTH: "เบริลเลียม", Number: 4, Shells: []int{2, 2}, Subshell: "1s² 2s²", Group: "IIA", Period: 2},
	{Symbol: "B", NameTH: "โบรอน", Number: 5, Shells: []int{2, 3}, Subshell: "1s² 2s² 2p¹", Group: "IIIA", Period: 2},
	{Symbol: "C", NameTH: "คาร์บอน", Number: 6, Shells: []int{2, 4}, Subshell: "1s² 2s² 2p²", Group: "IVA", Period: 2},
	{Symbol: "N", NameTH: "ไนโตรเจน", Number: 7, Shells: []int{2, 5}, Subshell: "1s² 2s² 2p³", Group: "VA", Period: 2},
	{Symbol: "O", NameTH: "ออกซิเจน", Number: 8, Shells: []int{2, 6}, Subshell: "1s² 2s² 2p⁴", Group: "VIA", Period: 2},
	{Symbol: "F", NameTH: "ฟลูออรีน", Number: 9, Shells: []int{2, 7}, Subshell: "1s² 2s² 2p⁵", Group: "VIIA", Period: 2},
	{Symbol: "Ne", NameTH: "นีออน", Number: 10, Shells: []int{2, 8}, Subshell: "1s² 2s² 2p⁶", Group: "VIIIA", Period: 2},
	{Symbol: "Na", NameTH: "โซเดียม", Number: 11, Shells: []int{2, 8, 1}, Subshell: "1s² 2s² 2p⁶ 3s¹", Group: "IA", Period: 3},
	{Symbol: "Mg", NameTH: "แมกนีเซียม", Number: 12, Shells: []int{2, 8, 2}, Subshell: "1s² 2s² 2p⁶ 3s²", Group: "IIA", Period: 3},
	{Symbol: "Al", NameTH: "อะลูมิเนียม", Number: 13, Shells: []int{2, 8, 3}, Subshell: "1s² 2s² 2p⁶ 3s² 3p¹", Group: "IIIA", Period: 3},
	{Symbol: "Si", NameTH: "ซิลิคอน", Number: 14, Shells: []int{2, 8, 4}, Subshell: "1s² 2s² 2p⁶ 3s² 3p²", Group: "IVA", Period: 3},
	{Symbol: "P", NameTH: "ฟอสฟอรัส", Number: 15, Shells: []int{2, 8, 5}, Subshell: "1s² 2s² 2p⁶ 3s² 3p³", Group: "VA", Period: 3},
	{Symbol: "S", NameTH: "กำมะถัน", Number: 16, Shells: []int{2, 8, 6}, Subshell: "1s² 2s² 2p⁶ 3s² 3p⁴", Group: "VIA", Period: 3},
	{Symbol: "Cl", NameTH: "คลอรีน", Number: 17, Shells: []int{2, 8, 7}, Subshell: "1s² 2s² 2p⁶ 3s² 3p⁵", Group: "VIIA", Period: 3},
	{Symbol: "Ar", NameTH: "อาร์กอน", Number: 18, Shells: []int{2, 8, 8}, Subshell: "1s² 2s² 2p⁶ 3s² 3p⁶", Group: "VIIIA", Period: 3},
	{Symbol: "K", NameTH: "โพแทสเซียม", Number: 19, Shells: []int{2, 8, 8, 1}, Subshell: "1s² 2s² 2p⁶ 3s² 3p⁶ 4s¹", Group: "IA", Period: 4},
	{Symbol: "Ca", NameTH: "แคลเซียม", Number: 20, Shells: []int{2, 8, 8, 2}, Subshell: "1s² 2s² 2p⁶ 3s² 3p⁶ 4s²", Group: "IIA", Period: 4},
}

// representativeGroups is the closed vocabulary a group question answers from.
var representativeGroups = []string{"IA", "IIA", "IIIA", "IVA", "VA", "VIA", "VIIA", "VIIIA"}

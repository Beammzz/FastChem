package services

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/takumi/fastchem/internal/models"
)

// RankedQuestionSet holds a deterministically-generated set of questions for a ranked match.
type RankedQuestionSet struct {
	Questions []RankedQuestion
}

// RankedQuestion is a single question with its metadata.
type RankedQuestion struct {
	Question     string   `json:"question"`
	Choices      []string `json:"choices"`
	CorrectIndex int      `json:"correctIndex"`
	TimeLimit    int      `json:"timeLimit"`
	Category     string   `json:"category"`
	Difficulty   string   `json:"difficulty"`
}

// GenerateRankedQuestions generates a deterministic set of 10 questions
// using the given seed. Both players receive the same questions.
// Distribution: 4 easy, 3 medium, 3 hard — shuffled.
func GenerateRankedQuestions(seed int64) *RankedQuestionSet {
	rng := rand.New(rand.NewSource(seed))

	questions := make([]RankedQuestion, 0, models.RankedQuestionsPerMatch)

	// Generate 4 easy questions
	for i := 0; i < models.RankedEasyCount; i++ {
		q := generateSeededEasy(rng)
		questions = append(questions, q)
	}

	// Generate 3 medium questions
	for i := 0; i < models.RankedMediumCount; i++ {
		q := generateSeededMedium(rng)
		questions = append(questions, q)
	}

	// Generate 3 hard questions
	for i := 0; i < models.RankedHardCount; i++ {
		q := generateSeededHard(rng)
		questions = append(questions, q)
	}

	// Shuffle the questions deterministically using the same seed-derived RNG
	rng.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return &RankedQuestionSet{Questions: questions}
}

// ─── Seeded Easy Generators ──────────────────────────────────────

func generateSeededEasy(rng *rand.Rand) RankedQuestion {
	if rng.Intn(2) == 0 {
		return seededAtomicStructure(rng)
	}
	return seededOxidation(rng)
}

func seededAtomicStructure(rng *rand.Rand) RankedQuestion {
	elem := elements[rng.Intn(len(elements))]
	atomicNumber := elem.AtomicNumber
	massNumber := atomicNumber + rng.Intn(11)

	questionType := rng.Intn(4)

	var questionText string
	var correctAnswer int
	var category string

	switch questionType {
	case 0:
		questionText = "ธาตุ " + elem.NameTH + " (" + elem.Symbol + ") มีเลขอะตอมเท่าใด?"
		correctAnswer = atomicNumber
		category = "atomic_structure"
	case 1:
		questionText = "ธาตุ " + elem.NameTH + " (" + elem.Symbol + ") มีจำนวนโปรตอนเท่าใด?"
		correctAnswer = atomicNumber
		category = "atomic_structure"
	case 2:
		questionText = "อะตอมที่เป็นกลางของ " + elem.NameTH + " (" + elem.Symbol + ") มีจำนวนอิเล็กตรอนเท่าใด?"
		correctAnswer = atomicNumber
		category = "atomic_structure"
	case 3:
		neutrons := massNumber - atomicNumber
		questionText = "ไอโซโทปของ " + elem.NameTH + " (" + elem.Symbol + ") ที่มีเลขมวล " + itoa(massNumber) + " มีจำนวนนิวตรอนเท่าใด?"
		correctAnswer = neutrons
		category = "atomic_structure"
	}

	choices, correctIndex := seededDistractors(rng, correctAnswer, 0, 40)

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("easy").TimeLimit,
		Category:     category,
		Difficulty:   "easy",
	}
}

func seededOxidation(rng *rand.Rand) RankedQuestion {
	compound := compounds[rng.Intn(len(compounds))]
	target := compound.Elements[rng.Intn(len(compound.Elements))]

	questionText := "เลขออกซิเดชันของ " + target.Symbol + " ใน " + compound.Formula + " คือเท่าใด?"

	choices, correctIndex := seededOxidationDistractors(rng, target.OxidationNumber)

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("easy").TimeLimit,
		Category:     "oxidation_number",
		Difficulty:   "easy",
	}
}

// ─── Seeded Medium Generators ────────────────────────────────────

func generateSeededMedium(rng *rand.Rand) RankedQuestion {
	switch rng.Intn(3) {
	case 0:
		return seededMassToMoles(rng)
	case 1:
		return seededMolesToMass(rng)
	default:
		return seededMolesToParticles(rng)
	}
}

func seededMassToMoles(rng *rand.Rand) RankedQuestion {
	elem := elements[rng.Intn(len(elements))]
	mass := roundTo(1.0+rng.Float64()*499.0, 1)
	moles := roundTo(mass/elem.MolarMass, 2)

	questionText := "ถ้ามี " + elem.NameTH + " (" + elem.Symbol + ") มวล " + ftoa(mass, 1) + " g จำนวนโมลคือเท่าใด?"
	choices, correctIndex := seededFloatDistractors(rng, moles, "%.2f mol")

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

func seededMolesToMass(rng *rand.Rand) RankedQuestion {
	elem := elements[rng.Intn(len(elements))]
	moles := roundTo(0.1+rng.Float64()*9.9, 2)
	mass := roundTo(moles*elem.MolarMass, 2)

	questionText := elem.NameTH + " (" + elem.Symbol + ") จำนวน " + ftoa(moles, 2) + " mol มีมวลกี่กรัม?"
	choices, correctIndex := seededFloatDistractors(rng, mass, "%.2f g")

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

func seededMolesToParticles(rng *rand.Rand) RankedQuestion {
	elem := elements[rng.Intn(len(elements))]
	moles := roundTo(0.1+rng.Float64()*4.9, 2)
	particles := moles * avogadro

	questionText := elem.NameTH + " (" + elem.Symbol + ") จำนวน " + ftoa(moles, 2) + " mol มีจำนวนอนุภาคเท่าใด? (Nₐ = 6.02×10²³)"
	choices, correctIndex := seededSciDistractors(rng, particles)

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

// ─── Seeded Hard Generators ─────────────────────────────────────

func generateSeededHard(rng *rand.Rand) RankedQuestion {
	switch rng.Intn(3) {
	case 0:
		return seededDilution(rng)
	case 1:
		return seededPreparingFromSolid(rng)
	default:
		return seededFreezingPoint(rng)
	}
}

func seededDilution(rng *rand.Rand) RankedQuestion {
	s := solutes[rng.Intn(len(solutes))]

	m1 := roundTo(0.5+rng.Float64()*5.5, 2)
	v1 := roundTo(10+rng.Float64()*490, 1)
	m2 := roundTo(0.01+rng.Float64()*(m1-0.02), 2)
	if m2 <= 0 {
		m2 = 0.01
	}
	v2 := roundTo((m1*v1)/m2, 1)

	questionText := "ต้องเจือจางสารละลาย " + s.Name + " (" + s.Formula + ") ความเข้มข้น " + ftoa(m1, 2) + " M ปริมาตร " + ftoa(v1, 1) + " mL ให้เหลือ " + ftoa(m2, 2) + " M ปริมาตรสุดท้ายคือเท่าใด?"
	choices, correctIndex := seededFloatDistractors(rng, v2, "%.1f mL")

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "dilution",
		Difficulty:   "hard",
	}
}

func seededPreparingFromSolid(rng *rand.Rand) RankedQuestion {
	s := solutes[rng.Intn(len(solutes))]

	molarity := roundTo(0.1+rng.Float64()*2.9, 2)
	volumeML := roundTo(50+rng.Float64()*950, 1)
	volumeL := volumeML / 1000.0
	mass := roundTo(molarity*volumeL*s.MolarMass, 2)

	questionText := "ต้องใช้ " + s.Name + " (" + s.Formula + ", MW=" + ftoa(s.MolarMass, 2) + " g/mol) กี่กรัม เพื่อเตรียมสารละลาย " + ftoa(molarity, 2) + " M ปริมาตร " + ftoa(volumeML, 1) + " mL?"
	choices, correctIndex := seededFloatDistractors(rng, mass, "%.2f g")

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "preparing_solution",
		Difficulty:   "hard",
	}
}

func seededFreezingPoint(rng *rand.Rand) RankedQuestion {
	s := solutes[rng.Intn(len(solutes))]

	molality := roundTo(0.1+rng.Float64()*4.9, 2)
	deltaTf := roundTo(float64(s.VantHoff)*kfWater*molality, 2)
	newFP := roundTo(0.0-deltaTf, 2)

	questionText := "สารละลาย " + s.Name + " (" + s.Formula + ", i=" + itoa(s.VantHoff) + ") ความเข้มข้น " + ftoa(molality, 2) + " mol/kg ในน้ำ (Kf=1.86 °C·kg/mol) จุดเยือกแข็งใหม่คือเท่าใด?"
	choices, correctIndex := seededFloatDistractors(rng, newFP, "%.2f °C")

	return RankedQuestion{
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "freezing_point",
		Difficulty:   "hard",
	}
}

// ─── Seeded distractor helpers ───────────────────────────────────

func seededDistractors(rng *rand.Rand, correct, minVal, maxVal int) ([]string, int) {
	choices := make(map[int]bool)
	choices[correct] = true

	for len(choices) < 4 {
		offset := rng.Intn(9) - 4
		if offset == 0 {
			offset = rng.Intn(2)*2 - 1
		}
		d := correct + offset
		if d < minVal {
			d = minVal + rng.Intn(5)
		}
		if d > maxVal {
			d = maxVal - rng.Intn(5)
		}
		if d != correct {
			choices[d] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, itoa(c))
	}
	rng.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctStr := itoa(correct)
	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}
	return choiceSlice, correctIndex
}

func seededOxidationDistractors(rng *rand.Rand, correct int) ([]string, int) {
	choices := make(map[int]bool)
	choices[correct] = true
	commonOx := []int{-3, -2, -1, 0, +1, +2, +3, +4, +5, +6, +7}

	for len(choices) < 4 {
		candidate := commonOx[rng.Intn(len(commonOx))]
		if candidate != correct {
			choices[candidate] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		if c > 0 {
			choiceSlice = append(choiceSlice, "+"+itoa(c))
		} else {
			choiceSlice = append(choiceSlice, itoa(c))
		}
	}
	rng.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	var correctStr string
	if correct > 0 {
		correctStr = "+" + itoa(correct)
	} else {
		correctStr = itoa(correct)
	}
	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}
	return choiceSlice, correctIndex
}

func seededFloatDistractors(rng *rand.Rand, correct float64, format string) ([]string, int) {
	choices := make(map[string]bool)
	correctStr := fmt.Sprintf(format, correct)
	choices[correctStr] = true

	for len(choices) < 4 {
		factor := 0.5 + rng.Float64()*1.0
		if factor > 0.9 && factor < 1.1 {
			factor = 1.3
		}
		d := correct * factor
		dStr := fmt.Sprintf(format, d)
		if dStr != correctStr {
			choices[dStr] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, c)
	}
	rng.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}
	return choiceSlice, correctIndex
}

func seededSciDistractors(rng *rand.Rand, correct float64) ([]string, int) {
	correctStr := formatSci(correct)
	choices := make(map[string]bool)
	choices[correctStr] = true

	for len(choices) < 4 {
		factor := 0.4 + rng.Float64()*1.2
		if factor > 0.9 && factor < 1.1 {
			factor = 1.5
		}
		d := correct * factor
		dStr := formatSci(d)
		if dStr != correctStr {
			choices[dStr] = true
		}
	}

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, c)
	}
	rng.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}
	return choiceSlice, correctIndex
}

// ─── Utility functions ──────────────────────────────────────────

func itoa(n int) string {
	return strconv.Itoa(n)
}

func ftoa(f float64, decimals int) string {
	format := "%." + strconv.Itoa(decimals) + "f"
	return fmt.Sprintf(format, f)
}

func roundTo(v float64, decimals int) float64 {
	pow := 1.0
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int(v*pow+0.5)) / pow
}

package services

import (
	"fmt"
	"math"
	"strings"
)

// solute data for hard questions
type solute struct {
	Name      string // Thai name
	Formula   string
	MolarMass float64 // g/mol
	VantHoff  int     // van't Hoff factor i
}

var solutes = []solute{
	{Name: "โซเดียมคลอไรด์", Formula: "NaCl", MolarMass: 58.44, VantHoff: 2},
	{Name: "กลูโคส", Formula: "C₆H₁₂O₆", MolarMass: 180.16, VantHoff: 1},
	{Name: "แคลเซียมคลอไรด์", Formula: "CaCl₂", MolarMass: 110.98, VantHoff: 3},
	{Name: "โพแทสเซียมคลอไรด์", Formula: "KCl", MolarMass: 74.55, VantHoff: 2},
	{Name: "ยูเรีย", Formula: "CO(NH₂)₂", MolarMass: 60.06, VantHoff: 1},
	{Name: "ซูโครส", Formula: "C₁₂H₂₂O₁₁", MolarMass: 342.30, VantHoff: 1},
	{Name: "โซเดียมไฮดรอกไซด์", Formula: "NaOH", MolarMass: 40.00, VantHoff: 2},
	{Name: "แมกนีเซียมซัลเฟต", Formula: "MgSO₄", MolarMass: 120.37, VantHoff: 2},
}

// Cryoscopic and ebullioscopic constants for water, in °C·kg/mol.
// Reaching for Kb where the problem calls for Kf is a standard slip, so it
// earns a distractor.
const (
	kfWater = 1.86
	kbWater = 0.512
)

// ─── Dilution: M₁V₁ = M₂V₂ ──────────────────────────────────────

type dilutionTopic struct{}

func (dilutionTopic) Category() string   { return "dilution" }
func (dilutionTopic) Difficulty() string { return "hard" }

func (t dilutionTopic) Generate(s *Source) Problem {
	const format = "%.1f mL"
	sol := solutes[s.Pick("solutes", len(solutes))]

	m1 := snapTo(0.5+s.Float64()*5.5, 0.01) // 0.50–6.00 M
	v1 := snapTo(10+s.Float64()*490, 0.1)   // 10.0–500.0 mL
	m2 := snapTo(0.01+s.Float64()*(m1-0.02), 0.01)
	if m2 <= 0 {
		m2 = 0.01
	}
	v2 := snapTo(m1*v1/m2, 0.1)

	candidates := floatCandidates(v2, format, []float64{
		v2 - v1,      // gave the volume of water to add, not the final volume
		m2 * v1 / m1, // inverted the concentration ratio
		m1 * v1 * m2, // multiplied through where the rule divides
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, v2), candidates,
		floatTopUp(v2, format))

	return Problem{
		Question: fmt.Sprintf(
			"ต้องเจือจางสารละลาย %s (%s) ความเข้มข้น %.2f M ปริมาตร %.1f mL ให้เหลือ %.2f M ปริมาตรสุดท้ายคือเท่าใด?",
			sol.Name, sol.Formula, m1, v1, m2),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Preparing from solid: mass = M × V(L) × MW ─────────────────

type preparingSolutionTopic struct{}

func (preparingSolutionTopic) Category() string   { return "preparing_solution" }
func (preparingSolutionTopic) Difficulty() string { return "hard" }

func (t preparingSolutionTopic) Generate(s *Source) Problem {
	const format = "%.2f g"
	sol := solutes[s.Pick("solutes", len(solutes))]

	molarity := snapTo(0.1+s.Float64()*2.9, 0.01) // 0.10–3.00 M
	volumeML := snapTo(50+s.Float64()*950, 0.1)   // 50.0–1000.0 mL
	volumeL := volumeML / 1000.0
	mass := snapTo(molarity*volumeL*sol.MolarMass, 0.01)

	candidates := floatCandidates(mass, format, []float64{
		volumeL * sol.MolarMass,             // dropped the concentration
		molarity * sol.MolarMass,            // dropped the volume
		molarity * volumeL / sol.MolarMass,  // divided by the molar mass
		molarity * volumeML * sol.MolarMass, // never converted mL to L
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, mass), candidates,
		floatTopUp(mass, format))

	return Problem{
		Question: fmt.Sprintf(
			"ต้องใช้ %s (%s, MW=%.2f g/mol) กี่กรัม เพื่อเตรียมสารละลาย %.2f M ปริมาตร %.1f mL?",
			sol.Name, sol.Formula, sol.MolarMass, molarity, volumeML),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Freezing point depression: ΔTf = i × Kf × m ────────────────

type freezingPointTopic struct{}

func (freezingPointTopic) Category() string   { return "freezing_point" }
func (freezingPointTopic) Difficulty() string { return "hard" }

func (t freezingPointTopic) Generate(s *Source) Problem {
	const format = "%.2f °C"
	sol := solutes[s.Pick("solutes", len(solutes))]

	molality := snapTo(0.1+s.Float64()*4.9, 0.01) // 0.10–5.00 mol/kg
	i := float64(sol.VantHoff)
	deltaTf := snapTo(i*kfWater*molality, 0.01)
	freezingPoint := snapTo(-deltaTf, 0.01)

	candidates := floatCandidates(freezingPoint, format, []float64{
		deltaTf,                 // reported the depression, not the new freezing point
		-kfWater * molality,     // omitted the van 't Hoff factor
		-i * kbWater * molality, // used Kb for water instead of Kf
		-i * kfWater / molality, // divided by molality where the rule multiplies
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, freezingPoint), candidates,
		floatTopUp(freezingPoint, format))

	return Problem{
		Question: fmt.Sprintf(
			"สารละลาย %s (%s, i=%d) ความเข้มข้น %.2f mol/kg ในน้ำ (Kf=1.86 °C·kg/mol) จุดเยือกแข็งใหม่คือเท่าใด?",
			sol.Name, sol.Formula, sol.VantHoff, molality),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Limiting reagent and percent yield (บทที่ 6) ─────────────────

type limitingReagentTopic struct{}

func (limitingReagentTopic) Category() string   { return "limiting_reagent" }
func (limitingReagentTopic) Difficulty() string { return "hard" }

func (t limitingReagentTopic) Generate(s *Source) Problem {
	r := twoReactantReactions[s.Pick("twoReactantReactions", len(twoReactantReactions))]
	a, b := r.Reactants[0], r.Reactants[1]
	product := r.Products[s.Intn(len(r.Products))]

	// Draw the two amounts independently so neither reactant is systematically
	// the limiting one, and the student has to actually divide by the
	// coefficients to find out which.
	nA := snapTo(0.5+s.Float64()*9.5, 0.01)
	nB := snapTo(0.5+s.Float64()*9.5, 0.01)

	extentA := nA / float64(a.Coef)
	extentB := nB / float64(b.Coef)
	extent := math.Min(extentA, extentB)

	excess := math.Max(extentA, extentB)
	if s.Intn(2) == 0 {
		return t.productAmount(s, r, a, b, product, nA, nB, extent, excess)
	}
	return t.percentYield(s, r, a, b, product, nA, nB, extent, excess)
}

func (t limitingReagentTopic) productAmount(
	s *Source, r reaction, a, b, product reactionSpecies,
	nA, nB, extent, excess float64,
) Problem {
	const format = "%.2f mol"
	answer := snapTo(extent*float64(product.Coef), 0.01)

	candidates := floatCandidates(answer, format, []float64{
		excess * float64(product.Coef),           // worked from the reactant that is in excess
		math.Min(nA, nB) * float64(product.Coef), // called the smaller amount limiting without dividing by the coefficients
		(nA + nB) * float64(product.Coef),        // added the reactants instead of choosing between them
		extent,                                   // stopped at the extent and never applied the product coefficient
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, answer), candidates,
		floatTopUp(answer, format))

	return Problem{
		Question: fmt.Sprintf(
			"จากสมการ %s ถ้าเริ่มด้วย %s %.2f mol และ %s %.2f mol จะเกิด %s ได้กี่โมล?",
			r.Equation, a.Formula, nA, b.Formula, nB, product.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t limitingReagentTopic) percentYield(
	s *Source, r reaction, a, b, product reactionSpecies,
	nA, nB, extent, excess float64,
) Problem {
	const format = "%.2f %%"

	theoretical := extent * float64(product.Coef) * product.MolarMass
	// A yield between 40% and 95% is what a real preparation gives, and it
	// keeps the answer clear of both 0 and 100.
	fraction := 0.40 + s.Float64()*0.55
	actual := snapTo(theoretical*fraction, 0.01)
	percent := snapTo(actual/theoretical*100, 0.01)

	candidates := percentCandidates(percent, format, []float64{
		(theoretical - actual) / theoretical * 100,                          // reported the loss instead of the yield
		actual / (excess * float64(product.Coef) * product.MolarMass) * 100, // sized the theoretical yield off the reactant in excess
		actual / theoretical, // never multiplied by 100
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, percent), candidates,
		floatTopUp(percent, format))

	return Problem{
		Question: fmt.Sprintf(
			"จากสมการ %s เริ่มด้วย %s %.2f mol และ %s %.2f mol ได้ %s (M = %.2f g/mol) จริง %.2f g ร้อยละผลได้เท่าใด?",
			r.Equation, a.Formula, nA, b.Formula, nB, product.Formula, product.MolarMass, actual),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Ideal gas law (บทที่ 7) ──────────────────────────────────────

type idealGasTopic struct{}

func (idealGasTopic) Category() string   { return "ideal_gas" }
func (idealGasTopic) Difficulty() string { return "hard" }

// gasConstantAtm is R in L·atm/(mol·K). gasConstantSI is the same constant in
// J/(mol·K) — reaching for it while working in atm and litres is the standard
// slip, so it earns a distractor rather than being hidden.
const (
	gasConstantAtm = 0.0821
	gasConstantSI  = 8.314
)

func (t idealGasTopic) Generate(s *Source) Problem {
	moles := snapTo(0.2+s.Float64()*4.8, 0.01)
	celsius := snapTo(10+s.Float64()*240, 1)
	kelvin := celsius + kelvinOffset

	if s.Intn(2) == 0 {
		return t.volume(s, moles, celsius, kelvin)
	}
	return t.moles(s, moles, celsius, kelvin)
}

// V = nRT / P
func (t idealGasTopic) volume(s *Source, moles, celsius, kelvin float64) Problem {
	const format = "%.2f L"

	pressure := snapTo(0.5+s.Float64()*3.5, 0.1)
	volume := snapTo(moles*gasConstantAtm*kelvin/pressure, 0.01)

	candidates := floatCandidates(volume, format, []float64{
		moles * gasConstantAtm * celsius / pressure, // stayed in celsius instead of kelvin
		moles * gasConstantSI * kelvin / pressure,   // used R = 8.314 while working in atm and litres
		moles * gasConstantAtm * kelvin * pressure,  // multiplied by the pressure instead of dividing
		pressure / (moles * gasConstantAtm * kelvin),
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, volume), candidates,
		floatTopUp(volume, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สอุดมคติ %.2f mol ที่ความดัน %.1f atm อุณหภูมิ %.0f °C มีปริมาตรเท่าใด? (R = 0.0821 L·atm/(mol·K))",
			moles, pressure, celsius),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// n = PV / RT
func (t idealGasTopic) moles(s *Source, _ float64, celsius, kelvin float64) Problem {
	const format = "%.3f mol"

	pressure := snapTo(0.5+s.Float64()*3.5, 0.1)
	volume := snapTo(1+s.Float64()*49, 0.1)
	n := snapTo(pressure*volume/(gasConstantAtm*kelvin), 0.001)

	candidates := floatCandidates(n, format, []float64{
		pressure * volume / (gasConstantAtm * celsius), // stayed in celsius
		pressure * volume / (gasConstantSI * kelvin),   // used R = 8.314
		gasConstantAtm * kelvin / (pressure * volume),  // inverted the whole expression
		pressure * volume * gasConstantAtm * kelvin,
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, n), candidates,
		floatTopUp(n, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สอุดมคติปริมาตร %.1f L ที่ความดัน %.1f atm อุณหภูมิ %.0f °C มีจำนวนกี่โมล? (R = 0.0821 L·atm/(mol·K))",
			volume, pressure, celsius),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Reaction rate (บทที่ 8) ──────────────────────────────────────

type reactionRateTopic struct{}

func (reactionRateTopic) Category() string   { return "reaction_rate" }
func (reactionRateTopic) Difficulty() string { return "hard" }

func (t reactionRateTopic) Generate(s *Source) Problem {
	if s.Intn(2) == 0 {
		return t.averageRate(s)
	}
	return t.rateOrder(s)
}

// average rate = −Δ[A] / Δt
func (t reactionRateTopic) averageRate(s *Source) Problem {
	const format = "%.4f mol/(L·s)"

	start := snapTo(0.50+s.Float64()*1.50, 0.01)
	end := snapTo(0.05+s.Float64()*(start-0.10), 0.01)
	seconds := snapTo(10+s.Float64()*190, 1)
	rate := snapTo((start-end)/seconds, 0.0001)

	candidates := floatCandidates(rate, format, []float64{
		(end - start) / seconds, // kept the sign of the change instead of reporting a positive rate
		start - end,             // never divided by the elapsed time
		start / seconds,         // used the initial concentration instead of the change
		end / seconds,
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, rate), candidates,
		floatTopUp(rate, format))

	return Problem{
		Question: fmt.Sprintf(
			"ความเข้มข้นของสารตั้งต้นลดลงจาก %.2f mol/L เหลือ %.2f mol/L ในเวลา %.0f วินาที อัตราการเกิดปฏิกิริยาเฉลี่ยเท่าใด?",
			start, end, seconds),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// rate = k[A]ⁿ, so multiplying [A] by f multiplies the rate by fⁿ
func (t reactionRateTopic) rateOrder(s *Source) Problem {
	factor := 2 + s.Intn(2) // 2 or 3
	order := 1 + s.Intn(3)  // first, second or third order
	answer := intPow(factor, order)

	// Every wrong option is a different way of combining the same two numbers,
	// which is exactly how this question is failed.
	candidates := intCandidates(1, 200, []int{
		factor * order,          // multiplied the factor by the order
		intPow(order, factor),   // swapped base and exponent
		factor,                  // assumed the rate scales one-for-one
		intPow(factor, order+1), // counted the exponent one too high
	}, formatInt)

	choices, idx := assembleChoices(s, formatInt(answer), candidates,
		intTopUp(answer, 1, 200, formatInt))

	return Problem{
		Question: fmt.Sprintf(
			"ปฏิกิริยาหนึ่งมีกฎอัตรา rate = k[A]^%d ถ้าเพิ่มความเข้มข้นของ A เป็น %d เท่า อัตราการเกิดปฏิกิริยาจะเป็นกี่เท่าของเดิม?",
			order, factor),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func intPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// ─── Chemical equilibrium (บทที่ 9) ───────────────────────────────

type equilibriumTopic struct{}

func (equilibriumTopic) Category() string   { return "equilibrium" }
func (equilibriumTopic) Difficulty() string { return "hard" }

// exponentEquilibria are the equilibria with at least one coefficient above 1.
// Where every coefficient is 1, "dropped the exponents" and "used the
// coefficients as multipliers" both render the correct expression, leaving too
// few distinct distractors for the expression question.
var exponentEquilibria = filterExponentEquilibria(equilibriumReactions)

func filterExponentEquilibria(all []equilibriumReaction) []equilibriumReaction {
	out := make([]equilibriumReaction, 0, len(all))
	for _, eq := range all {
		for _, sp := range append(append([]equilibriumSpecies{}, eq.Reactants...), eq.Products...) {
			if sp.Coef > 1 {
				out = append(out, eq)
				break
			}
		}
	}
	return out
}

// superscripts renders an exponent the way an equilibrium expression is
// printed; an exponent of 1 is left off, as the book leaves it off.
var superscripts = map[int]string{2: "²", 3: "³", 4: "⁴", 5: "⁵"}

func concentrationTerm(sp equilibriumSpecies, exponent int) string {
	return "[" + sp.Formula + "]" + superscripts[exponent]
}

// termProduct renders one side of the expression, using exponent(sp) for each
// species so the same renderer can produce the correct form and the
// misconception forms.
func termProduct(species []equilibriumSpecies, exponent func(equilibriumSpecies) int) string {
	parts := make([]string, len(species))
	for i, sp := range species {
		parts[i] = concentrationTerm(sp, exponent(sp))
	}
	joined := strings.Join(parts, "")
	if len(parts) > 1 {
		return "(" + joined + ")"
	}
	return joined
}

func (t equilibriumTopic) Generate(s *Source) Problem {
	if s.Intn(2) == 0 {
		return t.expression(s)
	}
	return t.constant(s)
}

func (t equilibriumTopic) expression(s *Source) Problem {
	eq := exponentEquilibria[s.Pick("exponentEquilibria", len(exponentEquilibria))]

	byCoef := func(sp equilibriumSpecies) int { return sp.Coef }
	one := func(equilibriumSpecies) int { return 1 }

	ratio := func(top, bottom []equilibriumSpecies, exponent func(equilibriumSpecies) int) string {
		return termProduct(top, exponent) + " / " + termProduct(bottom, exponent)
	}

	correct := ratio(eq.Products, eq.Reactants, byCoef)

	candidates := []string{
		ratio(eq.Reactants, eq.Products, byCoef), // wrote reactants over products
		ratio(eq.Products, eq.Reactants, one),    // dropped the coefficients instead of using them as exponents
		termProduct(eq.Products, byCoef),         // left out the denominator
		ratio(eq.Reactants, eq.Products, one),
	}

	choices, idx := assembleChoices(s, correct, candidates, noTopUp)

	return Problem{
		Question:     fmt.Sprintf("ค่าคงที่สมดุล K ของปฏิกิริยา %s เขียนได้อย่างไร?", eq.Equation),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t equilibriumTopic) constant(s *Source) Problem {
	eq := equilibriumReactions[s.Pick("equilibriumReactions", len(equilibriumReactions))]

	// Concentrations near 1 M keep K, and therefore 1/K, inside the band where
	// a wrong answer is still worth printing next to the right one.
	conc := make(map[string]float64, len(eq.Reactants)+len(eq.Products))
	lines := make([]string, 0, len(eq.Reactants)+len(eq.Products))
	for _, sp := range append(append([]equilibriumSpecies{}, eq.Reactants...), eq.Products...) {
		c := snapTo(0.50+s.Float64()*1.50, 0.01)
		conc[sp.Formula] = c
		lines = append(lines, fmt.Sprintf("[%s] = %.2f", sp.Formula, c))
	}

	power := func(species []equilibriumSpecies, exponent func(equilibriumSpecies) int) float64 {
		out := 1.0
		for _, sp := range species {
			out *= math.Pow(conc[sp.Formula], float64(exponent(sp)))
		}
		return out
	}
	byCoef := func(sp equilibriumSpecies) int { return sp.Coef }
	one := func(equilibriumSpecies) int { return 1 }

	const format = "%.3f"
	k := snapTo(power(eq.Products, byCoef)/power(eq.Reactants, byCoef), 0.001)

	candidates := floatCandidates(k, format, []float64{
		power(eq.Reactants, byCoef) / power(eq.Products, byCoef), // inverted the expression
		power(eq.Products, one) / power(eq.Reactants, one),       // ignored the coefficients
		power(eq.Products, byCoef),                               // left out the denominator
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, k), candidates,
		floatTopUp(k, format))

	return Problem{
		Question: fmt.Sprintf("ปฏิกิริยา %s ที่สมดุลมีความเข้มข้น (mol/L) ดังนี้ %s ค่าคงที่สมดุล K เท่ากับเท่าใด?",
			eq.Equation, strings.Join(lines, ", ")),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Acids and bases (บทที่ 10) ───────────────────────────────────

type acidBaseTopic struct{}

func (acidBaseTopic) Category() string   { return "acid_base" }
func (acidBaseTopic) Difficulty() string { return "hard" }

// pKw is pH + pOH at 25 °C.
const pKw = 14.0

func (t acidBaseTopic) Generate(s *Source) Problem {
	switch s.Intn(3) {
	case 0:
		return t.strongAcidPH(s)
	case 1:
		return t.weakAcidPH(s)
	default:
		return t.pOHToPH(s)
	}
}

// A strong acid dissociates completely: [H₃O⁺] = protons × concentration.
func (t acidBaseTopic) strongAcidPH(s *Source) Problem {
	const format = "%.2f"
	acid := strongAcids[s.Pick("strongAcids", len(strongAcids))]

	molarity := snapTo(0.001+s.Float64()*0.099, 0.001)
	h := float64(acid.Protons) * molarity
	ph := snapTo(-math.Log10(h), 0.01)

	candidates := floatCandidates(ph, format, []float64{
		-math.Log10(molarity), // ignored that H₂SO₄ supplies two protons
		pKw + math.Log10(h),   // answered with pOH
		math.Log10(h),         // dropped the minus sign
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, ph), candidates,
		floatTopUp(ph, format))

	return Problem{
		Question: fmt.Sprintf("สารละลาย %s (%s) เข้มข้น %.3f mol/L มีค่า pH เท่าใด?",
			acid.NameTH, acid.Formula, molarity),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// A weak acid barely dissociates: [H₃O⁺] ≈ √(Ka·C).
func (t acidBaseTopic) weakAcidPH(s *Source) Problem {
	const format = "%.2f"
	acid := weakAcids[s.Pick("weakAcids", len(weakAcids))]

	molarity := snapTo(0.01+s.Float64()*0.99, 0.01)
	h := math.Sqrt(acid.Ka * molarity)
	ph := snapTo(-math.Log10(h), 0.01)

	candidates := floatCandidates(ph, format, []float64{
		-math.Log10(acid.Ka * molarity), // forgot the square root
		-math.Log10(molarity),           // treated the weak acid as fully dissociated
		-math.Log10(acid.Ka),            // answered with pKa
		pKw + math.Log10(h),             // answered with pOH
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, ph), candidates,
		floatTopUp(ph, format))

	return Problem{
		Question: fmt.Sprintf("สารละลาย %s (%s, Ka = %s) เข้มข้น %.2f mol/L มีค่า pH เท่าใด?",
			acid.NameTH, acid.Formula, formatSci(acid.Ka), molarity),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// [OH⁻] → pOH → pH
func (t acidBaseTopic) pOHToPH(s *Source) Problem {
	const format = "%.2f"

	oh := snapTo(0.001+s.Float64()*0.099, 0.001)
	poh := -math.Log10(oh)
	ph := snapTo(pKw-poh, 0.01)

	candidates := floatCandidates(ph, format, []float64{
		poh,                      // stopped at pOH and called it pH
		pKw + poh,                // added where the relation subtracts
		-math.Log10(kWater / oh), // took the negative log of [H₃O⁺] but of the wrong quantity
		math.Log10(oh),
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, ph), candidates,
		floatTopUp(ph, format))

	return Problem{
		Question: fmt.Sprintf(
			"สารละลายเบสมี [OH⁻] = %.3f mol/L ที่ 25 °C มีค่า pH เท่าใด? (pH + pOH = 14)", oh),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Electrochemistry (บทที่ 11) ──────────────────────────────────

type electrochemistryTopic struct{}

func (electrochemistryTopic) Category() string   { return "electrochemistry" }
func (electrochemistryTopic) Difficulty() string { return "hard" }

func (t electrochemistryTopic) Generate(s *Source) Problem {
	if s.Intn(2) == 0 {
		i := s.Pick("standardPotentials", len(standardPotentials))
		j := s.Intn(len(standardPotentials) - 1)
		if j >= i {
			j++
		}

		// In a galvanic cell the higher standard reduction potential is
		// reduced, so it is the cathode and E°cell comes out positive.
		cathode, anode := standardPotentials[i], standardPotentials[j]
		if cathode.E0 < anode.E0 {
			cathode, anode = anode, cathode
		}
		return t.cellPotential(s, cathode, anode)
	}
	return t.strongestAgent(s)
}

func (t electrochemistryTopic) cellPotential(s *Source, cathode, anode halfCell) Problem {
	const format = "%.2f V"
	e := snapTo(cathode.E0-anode.E0, 0.01)

	candidates := floatCandidates(e, format, []float64{
		anode.E0 - cathode.E0, // subtracted in the wrong direction
		cathode.E0 + anode.E0, // added the two potentials
		math.Abs(cathode.E0) + math.Abs(anode.E0),
		cathode.E0, // quoted the cathode's potential as the cell's
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, e), candidates,
		floatTopUp(e, format))

	return Problem{
		Question: fmt.Sprintf(
			"เซลล์กัลวานิกประกอบด้วยครึ่งเซลล์ %s (E° = %+.2f V) และ %s (E° = %+.2f V) ค่า E°cell เท่ากับเท่าใด?",
			cathode.Reaction, cathode.E0, anode.Reaction, anode.E0),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// strongestAgent lists four half-reactions and asks which species is the
// strongest oxidising or reducing agent among them.
//
// Every option is one of the four half-cells on the table, so none can be
// eliminated for not being in play, and the two directions of the question
// share a table — reading it the wrong way round is the mistake being tested.
func (t electrochemistryTopic) strongestAgent(s *Source) Problem {
	const shown = 4

	// Draw four distinct half-cells by partial Fisher–Yates over an index
	// slice, so no index can repeat and the draw stays seeded.
	order := make([]int, len(standardPotentials))
	for i := range order {
		order[i] = i
	}
	s.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

	cells := make([]halfCell, 0, shown)
	for _, i := range order[:shown] {
		cells = append(cells, standardPotentials[i])
	}

	oxidising := s.Intn(2) == 0
	best := cells[0]
	for _, c := range cells[1:] {
		// The strongest oxidising agent is reduced most readily, at the top of
		// the table; the strongest reducing agent is at the bottom.
		if (oxidising && c.E0 > best.E0) || (!oxidising && c.E0 < best.E0) {
			best = c
		}
	}

	species := func(c halfCell) string {
		if oxidising {
			return c.Oxidised
		}
		return c.Reduced
	}

	lines := make([]string, 0, shown)
	pool := make([]string, 0, shown)
	for _, c := range cells {
		lines = append(lines, fmt.Sprintf("%s (E° = %+.2f V)", c.Reaction, c.E0))
		pool = append(pool, species(c))
	}

	correct := species(best)
	choices, idx := assembleChoices(s, correct, textCandidates(s, correct, pool), noTopUp)

	role := "ตัวออกซิไดซ์"
	if !oxidising {
		role = "ตัวรีดิวซ์"
	}

	return Problem{
		Question: fmt.Sprintf("จากค่า E° ต่อไปนี้ %s สารใดเป็น%sที่แรงที่สุด?",
			strings.Join(lines, ", "), role),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

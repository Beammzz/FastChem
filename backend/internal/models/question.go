package models

// Question represents a generated chemistry question.
//
// CorrectIndex is deliberately not serialized: a question is handed to the
// client before it is answered, so anything on the wire is readable from the
// network tab. The answer lives in GlobalQuestionStore keyed by ID, and comes
// back in AnswerResponse once the player has committed to a choice.
type Question struct {
	ID           string   `json:"id"`
	Question     string   `json:"question"`
	Choices      []string `json:"choices"`
	CorrectIndex int      `json:"-"`
	TimeLimit    int      `json:"timeLimit"`
	Category     string   `json:"category"`
	Difficulty   string   `json:"difficulty"`
}

// ValidateRequest is the payload for answer validation (legacy).
type ValidateRequest struct {
	QuestionID    string `json:"questionId" binding:"required"`
	SelectedIndex int    `json:"selectedIndex"`
	CorrectIndex  int    `json:"correctIndex"`
}

// ValidateResponse is the response after validation (legacy).
type ValidateResponse struct {
	Correct      bool   `json:"correct"`
	CorrectIndex int    `json:"correctIndex"`
	Message      string `json:"message"`
}

// AnswerRequest is the payload for the new per-question answer submission.
// The client only sends questionId and selectedIndex — the server looks up
// the correct answer and calculates time_spent from its own records.
type AnswerRequest struct {
	QuestionID    string `json:"questionId" binding:"required"`
	SelectedIndex int    `json:"selectedIndex"`
}

// AnswerResponse is returned after server-side answer validation and scoring.
type AnswerResponse struct {
	Correct      bool    `json:"correct"`
	CorrectIndex int     `json:"correctIndex"`
	ScoreEarned  int     `json:"scoreEarned"`
	SpeedBonus   int     `json:"speedBonus"`
	TimeSpent    float64 `json:"timeSpent"`
	Difficulty   string  `json:"difficulty"`
	TimedOut     bool    `json:"timedOut"`
}

// Element holds basic element data for question generation
type Element struct {
	Symbol       string
	Name         string
	NameTH       string
	AtomicNumber int
	MolarMass    float64 // g/mol – used for medium & hard questions
	Isotopes     []int   // real mass numbers, so neutron questions stay physical
}

// Compound holds compound data for oxidation number and state-of-matter questions
type Compound struct {
	Formula  string
	Name     string // English name for display
	NameTH   string // Thai name for display
	State    string // standard state at 25 °C, 1 atm: "s", "l", "g", or "aq"
	Elements []CompoundElement
}

// CompoundElement holds element info within a compound
type CompoundElement struct {
	Symbol          string
	OxidationNumber int
	Count           int
}

package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// A question is served before it is answered, so anything in its JSON is
// readable from the browser's network tab. The answer must not be there.
func TestQuestionJSONOmitsTheAnswer(t *testing.T) {
	q := Question{
		ID:           "q1",
		Question:     "ธาตุ คาร์บอน (C) มีจำนวนโปรตอนเท่าใด?",
		Choices:      []string{"4", "6", "8", "5"},
		CorrectIndex: 1,
		TimeLimit:    30,
		Category:     "atomic_structure",
		Difficulty:   "easy",
	}

	encoded, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("marshal question: %v", err)
	}
	if strings.Contains(string(encoded), "correctIndex") {
		t.Fatalf("question JSON leaks the answer: %s", encoded)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal question: %v", err)
	}
	for _, field := range []string{"id", "question", "choices", "timeLimit", "category", "difficulty"} {
		if _, ok := decoded[field]; !ok {
			t.Errorf("question JSON is missing %q, the client needs it to render", field)
		}
	}
}

// The answer is revealed only once the player has committed to a choice.
func TestAnswerResponseCarriesTheAnswer(t *testing.T) {
	encoded, err := json.Marshal(AnswerResponse{CorrectIndex: 2})
	if err != nil {
		t.Fatalf("marshal answer response: %v", err)
	}
	if !strings.Contains(string(encoded), `"correctIndex":2`) {
		t.Fatalf("answer response should reveal the answer, got %s", encoded)
	}
}

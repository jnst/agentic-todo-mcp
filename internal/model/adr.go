package model

import (
	"errors"
	"fmt"
)

// ADR represents an Architecture Decision Record
type ADR struct {
	Number       int    `json:"number"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Context      string `json:"context"`
	Decision     string `json:"decision"`
	Rationale    string `json:"rationale"`
	Consequences string `json:"consequences,omitempty"`
}

// NewADR creates a new ADR with the given parameters
func NewADR(number int, title, context, decision, rationale string) ADR {
	return ADR{
		Number:    number,
		Title:     title,
		Status:    "Proposed",
		Context:   context,
		Decision:  decision,
		Rationale: rationale,
	}
}

// Validate validates the ADR fields
func (a ADR) Validate() error {
	if a.Number <= 0 {
		return errors.New("ADR number must be positive")
	}
	if a.Title == "" {
		return errors.New("ADR title cannot be empty")
	}
	if a.Context == "" {
		return errors.New("ADR context cannot be empty")
	}
	if a.Decision == "" {
		return errors.New("ADR decision cannot be empty")
	}
	if a.Rationale == "" {
		return errors.New("ADR rationale cannot be empty")
	}

	validStatuses := map[string]bool{
		"Proposed":   true,
		"Accepted":   true,
		"Deprecated": true,
	}

	if !validStatuses[a.Status] {
		return fmt.Errorf("invalid status: %s", a.Status)
	}

	return nil
}

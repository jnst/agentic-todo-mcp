package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestADRCreation(t *testing.T) {
	// Red: 失敗するテストを書く
	want := ADR{
		Number:    1,
		Title:     "Test ADR",
		Status:    "Proposed",
		Context:   "Test context",
		Decision:  "Test decision",
		Rationale: "Test rationale",
	}

	got := NewADR(1, "Test ADR", "Test context", "Test decision", "Test rationale")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewADR() mismatch (-want +got):\n%s", diff)
	}
}

func TestADRStatusValidation(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		wantErr  bool
	}{
		{"valid Proposed", "Proposed", false},
		{"valid Accepted", "Accepted", false},
		{"valid Deprecated", "Deprecated", false},
		{"invalid status", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr := ADR{
				Number:    1,
				Title:     "Test ADR",
				Status:    tt.status,
				Context:   "Test context",
				Decision:  "Test decision",
				Rationale: "Test rationale",
			}
			err := adr.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ADR.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
package auth

import "testing"

func TestNormalizePatientStatus_DefaultFallback(t *testing.T) {
	got, err := normalizePatientStatus(nil, patientStatusActive)
	if err != nil {
		t.Fatalf("normalizePatientStatus error: %v", err)
	}
	if got == nil || *got != patientStatusActive {
		t.Fatalf("expected default %q, got %#v", patientStatusActive, got)
	}
}

func TestNormalizePatientStatus_AllowsKnownValues(t *testing.T) {
	cases := []string{
		patientStatusProvisional,
		patientStatusActive,
		patientStatusInactive,
		patientStatusDeceased,
		"ACTIVE",
	}

	for _, c := range cases {
		input := c
		got, err := normalizePatientStatus(&input, "")
		if err != nil {
			t.Fatalf("normalizePatientStatus(%q) error: %v", c, err)
		}
		if got == nil {
			t.Fatalf("normalizePatientStatus(%q) returned nil value", c)
		}
	}
}

func TestNormalizePatientStatus_RejectsInvalidValue(t *testing.T) {
	input := "archived"
	_, err := normalizePatientStatus(&input, "")
	if err != ErrInvalidPatientStatus {
		t.Fatalf("expected ErrInvalidPatientStatus, got %v", err)
	}
}

func TestNormalizePatientStatus_NilWhenEmptyInputAndFallback(t *testing.T) {
	got, err := normalizePatientStatus(nil, "")
	if err != nil {
		t.Fatalf("normalizePatientStatus error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

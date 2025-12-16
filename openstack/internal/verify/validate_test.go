package verify

import (
	"testing"
)

func TestUnit_ValidateDate_valid(t *testing.T) {
	validDates := []string{
		"2025-12-12T12:34:56Z",
		"2025-12-12T12:34:56.123Z",
	}

	for _, date := range validDates {
		_, errors := ValidateDate(date, "date")
		if len(errors) > 0 {
			t.Fatalf("expected the date %q to be in valid format, got error %q", date, errors)
		}
	}
}

func TestUnit_ValidateDate_invalid(t *testing.T) {
	invalidDates := []string{
		"a",
		"1234",
		"0000-00-00",
		"",
	}

	for _, date := range invalidDates {
		_, errors := ValidateDate(date, "date")
		if len(errors) == 0 {
			t.Fatalf("expected the date %q to fail validation", date)
		}
	}
}

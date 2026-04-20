package utils

import (
	"testing"
	"time"
)

func TestParseDataType1(t *testing.T) {
	tests := []struct {
		str      string
		expected time.Time
	}{
		{
			"Jan 29, 2026 10:49:04 PM YEKT",
			time.Date(2026, time.January, 29, 22, 49, 04, 0, time.Local),
		},
	}

	for i, test := range tests {
		result, err := ParseDataType1(test.str)
		if err != nil {
			t.Error(err)
		}

		if result != test.expected {
			t.Errorf("%d: Expected %v, got %v", i, test.expected, result)
		}
	}
}

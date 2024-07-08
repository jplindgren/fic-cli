package stock

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAddParameters(t *testing.T) {
	var tests = []struct {
		name   string
		ticker string
		target string
		want   bool
	}{
		{"empty ticker", "", "39.00", false},
		{"small ticker", "A", "39.00", false},
		{"empty target", "BVMF:EGIE3", "", false},
		{"invalid target", "BVMF:EGIE3", "AA", false},
		{"negative target", "BVMF:EGIE3", "-1", false},
		{"target too big", "BVMF:EGIE3", "11000", false},
		{"valid ticker and target", "BVMF:EGIE3", "39.00", true},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok, _ := isValid(tt.ticker, tt.target)
			if ok != tt.want {
				t.Errorf("got %t, want %t", ok, tt.want)
			}
		})
	}
}

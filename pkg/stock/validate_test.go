package stock

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	var tests = []struct {
		email string
		want  bool
	}{
		{"user@example.com", true},
		{"another.user@domain.com", true},
		{"invalid-email", false},
		{"@missing-username.com", false},
		{"username@.missing-domain", false},
		{"user@domain.c", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.want {
				t.Errorf("got %t, want %t", result, tt.want)
			}
		})
	}
}

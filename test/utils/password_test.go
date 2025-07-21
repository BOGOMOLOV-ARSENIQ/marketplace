package utils

import (
	"testing"

	"github.com/BOGOMOLOV-ARSENIQ/marketplace/pkg/utils"
)

func TestValidatePasswordComplexity(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Should return true for letters and digits",
			password: "password123",
			expected: true,
		},
		{
			name:     "Should return true for letters and special chars",
			password: "password!@#",
			expected: true,
		},
		{
			name:     "Should return true for digits and special chars",
			password: "12345!@#",
			expected: true,
		},
		{
			name:     "Should return true for all three types",
			password: "Password123!@#",
			expected: true,
		},
		{
			name:     "Should return false for only letters",
			password: "onlyletters",
			expected: false,
		},
		{
			name:     "Should return false for only digits",
			password: "123456789",
			expected: false,
		},
		{
			name:     "Should return false for only special chars",
			password: "!@#$%^&*",
			expected: false,
		},
		{
			name:     "Should return false for empty string",
			password: "",
			expected: false,
		},
		{
			name:     "Should return false for letters and spaces only",
			password: "hello world", // Spaces are not considered special characters by our regex
			expected: false,
		},
		{
			name:     "Should return false for digits and spaces only",
			password: "123 456",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidatePasswordComplexity(tt.password)
			if result != tt.expected {
				t.Errorf("ValidatePasswordComplexity(%q) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestHasDigit(t *testing.T) {
	if !utils.HasDigit("abc1def") {
		t.Error("Expected 'abc1def' to have a digit")
	}
	if utils.HasDigit("abcdef") {
		t.Error("Expected 'abcdef' not to have a digit")
	}
}

func TestHasLetter(t *testing.T) {
	if !utils.HasLetter("123a456") {
		t.Error("Expected '123a456' to have a letter")
	}
	if utils.HasLetter("123456") {
		t.Error("Expected '123456' not to have a letter")
	}
}

func TestHasSpecialChar(t *testing.T) {
	if !utils.HasSpecialChar("abc!def") {
		t.Error("Expected 'abc!def' to have a special character")
	}
	if utils.HasSpecialChar("abcdef123") {
		t.Error("Expected 'abcdef123' not to have a special character")
	}
	if utils.HasSpecialChar("abc def") {
		t.Error("Expected 'abc def' not to have a special character (space is excluded)")
	}
}

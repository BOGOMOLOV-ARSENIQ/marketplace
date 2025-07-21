package utils

func HasDigit(s string) bool {
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			return true
		}
	}
	return false
}

func HasLetter(s string) bool {
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			return true
		}
	}
	return false
}

func HasSpecialChar(s string) bool {
	for _, ch := range s {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
			return true
		}
	}
	return false
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ValidatePasswordComplexity(password string) bool {
	var flag int
	flag += btoi(HasDigit(password))
	flag += btoi(HasLetter(password))
	flag += btoi(HasSpecialChar(password))

	return flag >= 2
}

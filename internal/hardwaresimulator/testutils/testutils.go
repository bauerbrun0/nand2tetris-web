package testutils

func RepeatBool(value bool, count int) []bool {
	result := make([]bool, count)
	for i := range count {
		result[i] = value
	}
	return result
}

func StringToBoolArray(s string) []bool {
	result := make([]bool, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		char := s[i]
		if char == '1' {
			result[len(s)-1-i] = true
		} else {
			result[len(s)-1-i] = false
		}
	}
	return result
}

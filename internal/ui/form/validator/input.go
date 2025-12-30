package validator

import "unicode"

func IpValidator(textToCheck string, lastChar rune) bool {
	if textToCheck == "" {
		return true
	}

	if unicode.IsDigit(lastChar) {
		return true
	}

	if lastChar == '.' {
		return true
	}

	return false
}

func PortValidator(textToCheck string, lastChar rune) bool {
	if unicode.IsDigit(lastChar) {
		return true
	}

	return false
}

func DeviceValidator(textToCheck string, lastChar rune) bool {
	if unicode.IsLetter(lastChar) {
		return true
	}
	if unicode.IsDigit(lastChar) {
		return true
	}

	return false
}

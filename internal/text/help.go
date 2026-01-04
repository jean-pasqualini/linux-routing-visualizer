package text

import "regexp"

func TagSelectionRegions(text, query string) string {
	if query == "" {
		return text
	}

	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))

	return re.ReplaceAllStringFunc(text, func(m string) string {
		return `["selection"]` + m + `[""]`
	})
}

func TagColorRegions(color, defaultColor, text, query string) string {
	if query == "" {
		return text
	}

	re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(query))

	return re.ReplaceAllStringFunc(text, func(m string) string {
		return `[` + color + `]` + m + `[` + defaultColor + `]`
	})
}

var bracketRe = regexp.MustCompile(`\[[^\]]*\]`)

func RemoveSquareBrackets(s string) string {
	return bracketRe.ReplaceAllString(s, "")
}

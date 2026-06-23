package phone

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	localRegex  = regexp.MustCompile(`^09\d{9}$`)
	e164Regex   = regexp.MustCompile(`^\+989\d{9}$`)
	digitsRegex = regexp.MustCompile(`^\d{10,11}$`)
)

// Normalize converts Iranian phone numbers to E.164 (+989...).
func Normalize(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")

	switch {
	case e164Regex.MatchString(s):
		return s, nil
	case localRegex.MatchString(s):
		return "+98" + s[1:], nil
	case strings.HasPrefix(s, "989") && len(s) == 12:
		return "+" + s, nil
	case strings.HasPrefix(s, "9") && len(s) == 10:
		return "+98" + s, nil
	default:
		return "", fmt.Errorf("invalid phone number")
	}
}

// IsValid reports whether the phone can be normalized.
func IsValid(raw string) bool {
	_, err := Normalize(raw)
	return err == nil
}

// ToLocal converts E.164 to 09XXXXXXXXX for legacy validators.
func ToLocal(e164 string) string {
	if strings.HasPrefix(e164, "+989") && len(e164) == 13 {
		return "0" + e164[3:]
	}
	return e164
}

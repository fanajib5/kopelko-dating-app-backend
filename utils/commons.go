package utils

import "strings"

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		// return the original email if it doesn't have exactly one '@' character
		return email
	}
	local := parts[0]
	if len(local) <= 3 {
		return local + "@*****"
	}

	maskedLocal := local[:3] + strings.Repeat("*", 5)
	return maskedLocal + "@*****"
}

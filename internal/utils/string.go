package utils

import (
	"fmt"
	"regexp"
)

func IsAlphaNumeric(str string) bool {
	startsWithLetter := "[a-z]"
	allowedCharacters := "[a-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-z]"
	characterLimit := "{1,}"

	matchers := fmt.Sprintf(`^(%s%s*%s)%s$`, // ^[a-z][a-z0-9_-]*[0-9a-z]{1,}$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(str)
}

package patterns

import (
	"regexp"
)

func IsObjectReference(str string) bool {
	exp := regexp.MustCompile(`rec[0-9a-zA-Z]{14}`)
	return exp.Match([]byte(str))
}

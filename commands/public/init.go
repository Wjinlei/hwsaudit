package public

import "regexp"

var (
	regularizer *regexp.Regexp
)

func init() {
	regularizer = regexp.MustCompile(`(?i)\["[^\[]*"\]`)
}

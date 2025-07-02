//go:build !windows

package sanitize

import "strings"

func SanitizeOutputPath(s string) string {
	return strings.Replace(s, " ", "\\ ", -1)
}

//go:build windows

package sanitize

func SanitizeOutputPath(s string) string {
	return s
}

package utils

// StringOrDefault value
func StringOrDefault(s, def string) string {
	if len(s) > 0 {
		return s
	}
	return def
}

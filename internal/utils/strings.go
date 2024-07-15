package utils

// StringOrDefault value
func StringOrDefault(s, def string) string {
	if s != "" {
		return s
	}
	return def
}

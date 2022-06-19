// util
package util

// find out whether a string is in a slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Remove a string from a slice
func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			s = append(s[:i], s[i+1:]...)
			return s
		}
	}
	return s
}

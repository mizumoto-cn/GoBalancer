// util
package util

import "hash/fnv"

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

// Hash key with salt to 2 uint32 numbers modulo len(s)
func Hash(key string, salt string, len int) (uint32, uint32) {
	h := fnv.New32a()
	h.Write([]byte(key + salt))
	k1 := h.Sum32()
	h.Reset()
	h.Write([]byte(key + salt + salt))
	k2 := h.Sum32()
	return k1 % uint32(len), k2 % uint32(len)
}

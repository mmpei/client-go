package utils

import "os"

// FileExists check if the specified exists.
func FileExists(file string) bool {
	if len(file) > 0 {
		_, err := os.Stat(file)
		if err == nil {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return false
}

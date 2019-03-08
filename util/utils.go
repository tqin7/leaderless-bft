package util

func StringInArray(target string, arr []string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}
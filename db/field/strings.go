package field

func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func Offset(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}

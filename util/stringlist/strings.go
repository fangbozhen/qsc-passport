package stringlist

// do nothing if not found
func Remove(sl []string, s string) []string {
	maxIdx := len(sl) - 1
	for i := maxIdx; i >= 0; i-- {
		if s == sl[i] {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	return sl
}

// -1 if not found
func Find(sl []string, s string) int {
	for i, si := range sl {
		if si == s {
			return i
		}
	}
	return -1
}

func IsIn(sl []string, s string) bool {
	for _, si := range sl {
		if si == s {
			return true
		}
	}
	return false
}

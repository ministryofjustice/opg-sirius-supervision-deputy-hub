package util

func IsFeatureFlagged(features []string) func(string) bool {
	return func(f string) bool {
		for _, v := range features {
			if v == f {
				return true
			}
		}
		return false
	}
}

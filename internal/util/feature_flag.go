package util

var features = map[string][]string{
	"dev":        {"assurance_visits"},
	"production": {},
}

func FeatureFlag(env string) func(string) string {
	return func(f string) string {
		for _, v := range features[env] {
			if v == f {
				return ""
			}
		}
		return "hide"
	}
}

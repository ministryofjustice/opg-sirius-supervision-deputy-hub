package util

import "testing"

func TestTranslate(t *testing.T) {
	type args struct {
		prefix string
		s      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Translates value when prefix and s match",
			args{"FIELD", "email"},
			"Email address",
		},
		{
			"Returns original value when no translation exists",
			args{"", "Favourite colour"},
			"Favourite colour",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Translate(tt.args.prefix, tt.args.s); got != tt.want {
				t.Errorf("Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

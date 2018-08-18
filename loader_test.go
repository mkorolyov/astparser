package astparser

import (
	"regexp"
	"testing"
)

func Test_validFile(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		include *regexp.Regexp
		exclude *regexp.Regexp
		want    bool
	}{
		{
			name:    "include ok",
			s:       "event.go",
			include: regexp.MustCompile("event"),
			want:    true,
		},
		{
			name:    "include dont match",
			s:       "event.go",
			include: regexp.MustCompile("type"),
			want:    false,
		},
		{
			name: "valid go file",
			s:    "event.go",
			want: true,
		},
		{
			name: "test file",
			s:    "event_test.go",
			want: false,
		},
		{
			name:    "exclude ok",
			s:       "event.go",
			exclude: regexp.MustCompile("event"),
			want:    false,
		},
		{
			name:    "exclude dont match",
			s:       "event.go",
			exclude: regexp.MustCompile("type"),
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validFile(tt.s, tt.include, tt.exclude); got != tt.want {
				t.Errorf("validFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

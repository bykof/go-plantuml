package astParser

import (
	"regexp"
	"testing"
)

func Test_isExcluded(t *testing.T) {
	type args struct {
		fileName string
		regex    *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"positive_non_go_file_returns_true", args{"parser.js", nil}, true},
		{"positive_go_test_file_returns_true", args{"parser_test.go", nil}, true},
		{"positive_match_regex_returns_true", args{"parser_mock.go", regexp.MustCompile(`^.+_mock.go$`)}, true},
		{"positive_match_regex_due_to_dot_returns_true", args{"parser_mock|tmp.go", regexp.MustCompile(`^.+_mock.tmp.go$`)}, true},
		{"negative_go_file_no_regex_returns_false", args{"parser.go", nil}, false},
		{"negative_go_file_did_not_match_regex_returns_false", args{"parser_mick.go", regexp.MustCompile(`^.+_mock.go$`)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExcluded(tt.args.fileName, tt.args.regex); got != tt.want {
				t.Errorf("isExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}

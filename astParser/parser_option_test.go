package astParser

import (
	"reflect"
	"regexp"
	"testing"
)

func TestWithRecursive(t *testing.T) {
	tests := []struct {
		name string
		want *parserOptions
	}{
		{"positive", &parserOptions{recursive: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			op := &parserOptions{}
			opt := WithRecursive()
			opt(op)
			if !reflect.DeepEqual(op, tt.want) {
				t.Errorf("WithRecursive() = %v, want %v", op, tt.want)
			}
		})
	}
}

func TestWithFileExclusion(t *testing.T) {
	type args struct {
		exclusionRegex string
	}
	tests := []struct {
		name string
		args args
		want *parserOptions
	}{
		{"positive", args{"^somefile.go$"}, &parserOptions{excludedFilesRegex: regexp.MustCompile("^somefile.go$")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			op := &parserOptions{}
			opt := WithFileExclusion(tt.args.exclusionRegex)
			opt(op)
			if !reflect.DeepEqual(op, tt.want) {
				t.Errorf("WithRecursive() = %v, want %v", op, tt.want)
			}
		})
	}
}

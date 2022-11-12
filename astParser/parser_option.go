package astParser

import (
	"regexp"
)

type parserOptions struct {
	recursive          bool
	excludedFilesRegex *regexp.Regexp
}

type ParserOptionFunc func(*parserOptions)

func WithRecursive() func(*parserOptions) {
	return func(opt *parserOptions) {
		opt.recursive = true
	}
}

func WithFileExclusion(exclusionRegex string) func(*parserOptions) {
	return func(opt *parserOptions) {
		var re = regexp.MustCompile(exclusionRegex)
		opt.excludedFilesRegex = re
	}
}

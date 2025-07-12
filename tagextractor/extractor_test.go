package tagextractor

import (
	"article-processing-microservice/utils"
	"testing"
)

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		n        int
		expected []string
		err      error
	}{
		{
			name:     "Empty string",
			body:     "",
			n:        0,
			expected: []string{},
			err: nil,
		},
		{
			name:     "Single word",
			body:     "Golang",
			n:        1,
			expected: []string{"golang"},
			err: nil,
		},
		{
			name:     "Multiple words sorted",
			body:     "apple banana cherry",
			n:        3,
			expected: []string{"apple", "banana", "cherry"},
			err: nil,
		},
		{
			name:     "Multiple words sorted with punctuation and case",
			body:     "aPple, ban.ana,! cheRRy",
			n:        3,
			expected: []string{"ana", "apple", "ban"},
			err: nil,
		},
		{
			name:     "Multiple words sorted with stopwords",
			body:     "apple banana cherry the and of in to is",
			n:        3,
			expected: []string{"apple", "banana", "cherry"},
			err: nil,
		},
		{
			name:     "Wrong input",
			body:     "A! B.",
			n:        3,
			expected: nil,
			err: ErrWrongHighTagCount,
		},
	}

	for _, tt := range tests {
		realSlide, realErr := ExtractTags(tt.body, tt.n)
		if realErr != tt.err {
			t.Errorf("The error is wrong in %s", tt.name)
		}
		if !utils.IsSlideEqual(realSlide, tt.expected) {
			t.Errorf("ExtractTags(%q, %d) = %v, want %v", tt.body, tt.n, realSlide, tt.expected)
		}
	}
}
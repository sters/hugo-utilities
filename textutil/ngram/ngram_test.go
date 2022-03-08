package ngram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, []NGramGroup{
		{'こ', 'ん'},
		{'ん', 'に'},
		{'に', 'ち'},
		{'ち', 'は'},
	}, Parse(2, "こんにちは").Groups)

	assert.Equal(t, []NGramGroup{
		{'こ', 'ん', 'に'},
		{'ん', 'に', 'ち'},
		{'に', 'ち', 'は'},
	}, Parse(3, "こんにちは").Groups)
}

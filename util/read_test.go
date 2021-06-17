package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dirwalk(t *testing.T) {
	results, err := dirwalk("../tools")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Contains(t, results[0], "tools/tools.go")
}

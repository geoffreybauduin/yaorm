package yaorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSelectColumns(t *testing.T) {
	bs := buildSelectColumns{}
	assert.Equal(t, []string{"a", "b", "c"}, bs.reduce([]string{"a", "b", "c"}),
		"no filter on columns")

	bs = buildSelectColumns{dontLoadColumns: []string{"c", "d"}}
	assert.Equal(t, []string{"a", "b"}, bs.reduce([]string{"a", "b", "c"}),
		"do not use c, d columns")

	bs = buildSelectColumns{loadColumns: []string{"b", "c", "d"}}
	assert.Equal(t, []string{"b", "c"}, bs.reduce([]string{"a", "b", "c"}),
		"only use b, c, d columns")

	bs = buildSelectColumns{
		loadColumns:     []string{"b", "c", "d", "g"},
		dontLoadColumns: []string{"b", "e", "f"},
	}
	assert.Equal(t, []string{"c", "d"}, bs.reduce([]string{"a", "b", "c", "d", "e"}),
		"only use b, c, d, g columns, BUT do not use b, e, f columns")
}

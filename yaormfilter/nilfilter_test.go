package yaormfilter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

func TestNewNilFilter(t *testing.T) {
	filter := yaormfilter.NewNilFilter()
	assert.IsType(t, &yaormfilter.NilFilter{}, filter)
}

package yaormfilter_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestNewNilFilter(t *testing.T) {
	filter := yaormfilter.NewNilFilter()
	assert.IsType(t, &yaormfilter.NilFilter{}, filter)
}

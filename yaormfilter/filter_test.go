package yaormfilter_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

type fakeFilter struct {
	yaormfilter.ModelFilter
}

func (ff *fakeFilter) AddOption(opt yaormfilter.RequestOption) yaormfilter.Filter {
	ff.AddOption_(opt)
	return ff
}

func TestFilter_GetSelectOptions(t *testing.T) {
	filter := &fakeFilter{}
	assert.Len(t, filter.GetSelectOptions(), 0)

	filter.AddOption(yaormfilter.RequestOptions.SelectForUpdate)
	assert.Len(t, filter.GetSelectOptions(), 1)
	assert.Equal(t, filter.GetSelectOptions()[0], yaormfilter.RequestOptions.SelectForUpdate)
}
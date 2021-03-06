package container

import (
	"testing"

	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestBuildContainerListOptions(t *testing.T) {
	filters := opts.NewFilterOpt()
	assert.NilError(t, filters.Set("foo=bar"))
	assert.NilError(t, filters.Set("baz=foo"))

	contexts := []struct {
		psOpts          *psOptions
		expectedAll     bool
		expectedSize    bool
		expectedLimit   int
		expectedFilters map[string]string
	}{
		{
			psOpts: &psOptions{
				all:    true,
				size:   true,
				last:   5,
				filter: filters,
			},
			expectedAll:   true,
			expectedSize:  true,
			expectedLimit: 5,
			expectedFilters: map[string]string{
				"foo": "bar",
				"baz": "foo",
			},
		},
		{
			psOpts: &psOptions{
				all:     true,
				size:    true,
				last:    -1,
				nLatest: true,
			},
			expectedAll:     true,
			expectedSize:    true,
			expectedLimit:   1,
			expectedFilters: make(map[string]string),
		},
	}

	for _, c := range contexts {
		options, err := buildContainerListOptions(c.psOpts)
		assert.NilError(t, err)

		assert.Equal(t, c.expectedAll, options.All)
		assert.Equal(t, c.expectedSize, options.Size)
		assert.Equal(t, c.expectedLimit, options.Limit)
		assert.Equal(t, options.Filter.Len(), len(c.expectedFilters))

		for k, v := range c.expectedFilters {
			f := options.Filter
			if !f.ExactMatch(k, v) {
				t.Fatalf("Expected filter with key %s to be %s but got %s", k, v, f.Get(k))
			}
		}
	}
}

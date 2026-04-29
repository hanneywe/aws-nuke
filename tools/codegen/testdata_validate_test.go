package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSampleDSLFiles_ParseAndValidate(t *testing.T) {
	files := []string{
		"testdata/flat-simple.yaml",
		"testdata/flat-with-describe.yaml",
		"testdata/flat-with-tags.yaml",
		"testdata/nested-two-level.yaml",
		"testdata/nested-three-level.yaml",
		"testdata/singleton.yaml",
		"testdata/with-filters.yaml",
		"testdata/with-settings.yaml",
		"testdata/with-pre-deletion.yaml",
		"testdata/with-overrides.yaml",
		"testdata/with-integration-test.yaml",
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			def, err := Parse(f)
			require.NoError(t, err, "Parse failed for %s", f)
			require.NotNil(t, def)

			errs := Validate(def)
			assert.Empty(t, errs, "Validation errors for %s: %v", f, errs)
		})
	}
}

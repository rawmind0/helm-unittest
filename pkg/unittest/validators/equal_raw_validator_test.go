package validators_test

import (
	"testing"

	"github.com/lrills/helm-unittest/internal/common"
	. "github.com/lrills/helm-unittest/pkg/unittest/validators"
	"github.com/stretchr/testify/assert"
)

var docToTestEqualRaw = `
raw: This is a NOTES.txt document.
`

func TestEqualRawValidatorWhenOk(t *testing.T) {
	manifest := makeManifest(docToTestEqualRaw)
	validator := EqualRawValidator{"This is a NOTES.txt document."}

	pass, diff := validator.Validate(&ValidateContext{
		Docs: []common.K8sManifest{manifest},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestEqualRawValidatorWhenNegativeAndOk(t *testing.T) {
	manifest := makeManifest(docToTestEqualRaw)

	validator := EqualRawValidator{"Invalid text."}
	pass, diff := validator.Validate(&ValidateContext{
		Docs:     []common.K8sManifest{manifest},
		Negative: true,
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestEqualRawValidatorWhenFail(t *testing.T) {
	manifest := makeManifest(docToTestEqualRaw)

	validator := EqualRawValidator{"Invalid text."}
	pass, diff := validator.Validate(&ValidateContext{
		Docs: []common.K8sManifest{manifest},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"Expected to equal:",
		"	Invalid text.",
		"Actual:",
		"	This is a NOTES.txt document.",
		"Diff:",
		"	--- Expected",
		"	+++ Actual",
		"	@@ -1,2 +1,2 @@",
		"	-Invalid text.",
		"	+This is a NOTES.txt document.",
	}, diff)
}

func TestEqualRawValidatorWhenNegativeAndFail(t *testing.T) {
	manifest := makeManifest(docToTestEqualRaw)

	v := EqualRawValidator{"This is a NOTES.txt document."}
	pass, diff := v.Validate(&ValidateContext{
		Docs:     []common.K8sManifest{manifest},
		Negative: true,
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"Expected NOT to equal:",
		"	This is a NOTES.txt document.",
	}, diff)
}

func TestEqualRawValidatorWhenInvalidIndex(t *testing.T) {
	manifest := makeManifest(docToTestEqualRaw)
	validator := EqualRawValidator{"Invalid text."}
	pass, diff := validator.Validate(&ValidateContext{
		Docs:  []common.K8sManifest{manifest},
		Index: 2,
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"Error:",
		"	documentIndex 2 out of range",
	}, diff)
}

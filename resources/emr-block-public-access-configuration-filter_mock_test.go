package resources

import (
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_Mock_EMRBlockPublicAccessConfiguration_Filter_AlreadyFalse(t *testing.T) {
	a := assert.New(t)
	r := &EMRBlockPublicAccessConfiguration{
		BlockPublicSecurityGroupRules: ptr.Bool(false),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "default configuration")
}

func Test_Mock_EMRBlockPublicAccessConfiguration_Filter_True(t *testing.T) {
	a := assert.New(t)
	r := &EMRBlockPublicAccessConfiguration{
		BlockPublicSecurityGroupRules: ptr.Bool(true),
	}
	a.NoError(r.Filter())
}

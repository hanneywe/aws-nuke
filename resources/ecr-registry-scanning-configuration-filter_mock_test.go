package resources

import (
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"

	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func Test_Mock_ECRRegistryScanningConfiguration_Filter_DefaultBasic(t *testing.T) {
	a := assert.New(t)
	r := &ECRRegistryScanningConfiguration{
		ScanType: ptr.String("BASIC"),
		Rules:    []ecrtypes.RegistryScanningRule{},
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "default configuration")
}

func Test_Mock_ECRRegistryScanningConfiguration_Filter_Enhanced(t *testing.T) {
	a := assert.New(t)
	r := &ECRRegistryScanningConfiguration{
		ScanType: ptr.String("ENHANCED"),
		Rules:    []ecrtypes.RegistryScanningRule{},
	}
	a.NoError(r.Filter())
}

func Test_Mock_ECRRegistryScanningConfiguration_Filter_BasicWithRules(t *testing.T) {
	a := assert.New(t)
	r := &ECRRegistryScanningConfiguration{
		ScanType: ptr.String("BASIC"),
		Rules: []ecrtypes.RegistryScanningRule{
			{ScanFrequency: ecrtypes.ScanFrequencyScanOnPush},
		},
	}
	a.NoError(r.Filter())
}

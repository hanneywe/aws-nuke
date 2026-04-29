package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func Test_Mock_ECRRegistryScanningConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("GetRegistryScanningConfiguration", mock.Anything, mock.Anything).
		Return(&ecr.GetRegistryScanningConfigurationOutput{
			ScanningConfiguration: &ecrtypes.RegistryScanningConfiguration{
				ScanType: ecrtypes.ScanTypeEnhanced,
			},
		}, nil)
	lister := &ECRRegistryScanningConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("ENHANCED", *resources[0].(*ECRRegistryScanningConfiguration).ScanType)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRegistryScanningConfiguration_List_NoConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("GetRegistryScanningConfiguration", mock.Anything, mock.Anything).
		Return(&ecr.GetRegistryScanningConfigurationOutput{}, nil)
	lister := &ECRRegistryScanningConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRegistryScanningConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRRegistryScanningConfiguration{
		svc:      mockClient,
		ScanType: ptr.String("ENHANCED"),
	}
	mockClient.On("PutRegistryScanningConfiguration", mock.Anything, mock.Anything).
		Return(&ecr.PutRegistryScanningConfigurationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRegistryScanningConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRRegistryScanningConfiguration{
		ScanType: ptr.String("ENHANCED"),
	}
	props := r.Properties()
	a.Equal("ENHANCED", props.Get("ScanType"))
}

func Test_Mock_ECRRegistryScanningConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &ECRRegistryScanningConfiguration{ScanType: ptr.String("ENHANCED")}
	a.Equal("Registry Scanning Configuration", r.String())
}

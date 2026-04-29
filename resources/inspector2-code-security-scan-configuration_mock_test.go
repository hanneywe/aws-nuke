package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	inspector2types "github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

const TestScanConfigArn = "arn:aws:inspector2:us-east-1:123456789012:" +
	"code-security/scan-configuration/test"

func Test_Mock_Inspector2CodeSecurityScanConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockInspector2V2Client)
	mockClient.On("ListCodeSecurityScanConfigurations", mock.Anything, mock.Anything).
		Return(&inspector2.ListCodeSecurityScanConfigurationsOutput{
			Configurations: []inspector2types.CodeSecurityScanConfigurationSummary{
				{ScanConfigurationArn: ptr.String(TestScanConfigArn)},
			},
		}, nil)
	lister := &Inspector2CodeSecurityScanConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testInspector2V2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal(TestScanConfigArn,
		resources[0].(*Inspector2CodeSecurityScanConfiguration).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_Inspector2CodeSecurityScanConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockInspector2V2Client)
	mockClient.On("ListCodeSecurityScanConfigurations", mock.Anything, mock.Anything).
		Return(&inspector2.ListCodeSecurityScanConfigurationsOutput{
			Configurations: []inspector2types.CodeSecurityScanConfigurationSummary{},
		}, nil)
	lister := &Inspector2CodeSecurityScanConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testInspector2V2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Inspector2CodeSecurityScanConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockInspector2V2Client)
	r := &Inspector2CodeSecurityScanConfiguration{
		svc:                  mockClient,
		ScanConfigurationArn: ptr.String(TestScanConfigArn),
	}
	mockClient.On("DeleteCodeSecurityScanConfiguration", mock.Anything,
		&inspector2.DeleteCodeSecurityScanConfigurationInput{
			ScanConfigurationArn: r.ScanConfigurationArn,
		}).Return(&inspector2.DeleteCodeSecurityScanConfigurationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Inspector2CodeSecurityScanConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := Inspector2CodeSecurityScanConfiguration{
		ScanConfigurationArn: ptr.String(TestScanConfigArn),
	}
	a.Equal(TestScanConfigArn,
		r.Properties().Get("ScanConfigurationArn"))
}

func Test_Mock_Inspector2CodeSecurityScanConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &Inspector2CodeSecurityScanConfiguration{
		ScanConfigurationArn: ptr.String(TestScanConfigArn),
	}
	a.Equal(TestScanConfigArn, r.String())
}

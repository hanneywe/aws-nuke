package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"
)

func Test_Mock_RedshiftHsmConfiguration_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeHsmConfigurations", mock.Anything, mock.Anything).
		Return(&redshift.DescribeHsmConfigurationsOutput{
			HsmConfigurations: []redshifttypes.HsmConfiguration{
				{
					HsmConfigurationIdentifier: ptr.String("my-hsm-config"),
				},
			},
		}, nil)

	lister := &RedshiftHsmConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	hsmConfig := resources[0].(*RedshiftHsmConfiguration)
	assertions.Equal("my-hsm-config", *hsmConfig.HsmConfigurationIdentifier)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmConfiguration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeHsmConfigurations", mock.Anything, mock.Anything).
		Return(&redshift.DescribeHsmConfigurationsOutput{
			HsmConfigurations: []redshifttypes.HsmConfiguration{},
		}, nil)

	lister := &RedshiftHsmConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmConfiguration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	hsmConfig := &RedshiftHsmConfiguration{
		svc:                        mockClient,
		HsmConfigurationIdentifier: ptr.String("my-hsm-config"),
	}

	mockClient.On("DeleteHsmConfiguration", mock.Anything, &redshift.DeleteHsmConfigurationInput{
		HsmConfigurationIdentifier: hsmConfig.HsmConfigurationIdentifier,
	}).Return(&redshift.DeleteHsmConfigurationOutput{}, nil)

	err := hsmConfig.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftHsmConfiguration_Properties(t *testing.T) {
	assertions := assert.New(t)

	hsmConfig := RedshiftHsmConfiguration{
		HsmConfigurationIdentifier: ptr.String("my-hsm-config"),
	}

	properties := hsmConfig.Properties()
	assertions.Equal("my-hsm-config", properties.Get("HsmConfigurationIdentifier"))
}

func Test_Mock_RedshiftHsmConfiguration_String(t *testing.T) {
	assertions := assert.New(t)
	hsmConfig := RedshiftHsmConfiguration{HsmConfigurationIdentifier: ptr.String("my-hsm-config")}
	assertions.Equal("my-hsm-config", hsmConfig.String())
}

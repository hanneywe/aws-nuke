package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

func Test_Mock_ConfigServiceRetentionConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeRetentionConfigurations", mock.Anything, mock.Anything).
		Return(&configservice.DescribeRetentionConfigurationsOutput{
			RetentionConfigurations: []configtypes.RetentionConfiguration{
				{Name: ptr.String("default")},
			},
		}, nil)
	lister := &ConfigServiceRetentionConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("default", *resources[0].(*ConfigServiceRetentionConfiguration).Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceRetentionConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeRetentionConfigurations", mock.Anything, mock.Anything).
		Return(&configservice.DescribeRetentionConfigurationsOutput{RetentionConfigurations: []configtypes.RetentionConfiguration{}}, nil)
	lister := &ConfigServiceRetentionConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceRetentionConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	r := &ConfigServiceRetentionConfiguration{svc: mockClient, Name: ptr.String("default")}
	mockClient.On("DeleteRetentionConfiguration", mock.Anything,
		&configservice.DeleteRetentionConfigurationInput{RetentionConfigurationName: r.Name}).
		Return(&configservice.DeleteRetentionConfigurationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceRetentionConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := ConfigServiceRetentionConfiguration{Name: ptr.String("default")}
	a.Equal("default", r.Properties().Get("Name"))
}

func Test_Mock_ConfigServiceRetentionConfiguration_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("default", (&ConfigServiceRetentionConfiguration{Name: ptr.String("default")}).String())
}

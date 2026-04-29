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

func Test_Mock_ConfigServiceConfigurationAggregator_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeConfigurationAggregators", mock.Anything, mock.Anything).
		Return(&configservice.DescribeConfigurationAggregatorsOutput{
			ConfigurationAggregators: []configtypes.ConfigurationAggregator{
				{
					ConfigurationAggregatorName: ptr.String("my-aggregator"),
					ConfigurationAggregatorArn:  ptr.String("arn:aws:config:us-east-1:123456789012:config-aggregator/config-aggregator-abcdef"),
				},
			},
		}, nil)
	mockClient.On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(&configservice.ListTagsForResourceOutput{
			Tags: []configtypes.Tag{
				{Key: ptr.String("env"), Value: ptr.String("test")},
			},
		}, nil)
	lister := &ConfigServiceConfigurationAggregatorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*ConfigServiceConfigurationAggregator)
	a.Equal("my-aggregator", r.String())
	a.Equal("test", r.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceConfigurationAggregator_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeConfigurationAggregators", mock.Anything, mock.Anything).
		Return(&configservice.DescribeConfigurationAggregatorsOutput{
			ConfigurationAggregators: []configtypes.ConfigurationAggregator{},
		}, nil)
	lister := &ConfigServiceConfigurationAggregatorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceConfigurationAggregator_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	r := &ConfigServiceConfigurationAggregator{
		svc:  mockClient,
		Name: ptr.String("my-aggregator"),
	}
	mockClient.On("DeleteConfigurationAggregator", mock.Anything, &configservice.DeleteConfigurationAggregatorInput{
		ConfigurationAggregatorName: r.Name,
	}).Return(&configservice.DeleteConfigurationAggregatorOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceConfigurationAggregator_Properties(t *testing.T) {
	a := assert.New(t)
	r := ConfigServiceConfigurationAggregator{
		Name: ptr.String("my-aggregator"),
		Tags: map[string]string{"env": "test"},
	}
	a.Equal("my-aggregator", r.Properties().Get("Name"))
	a.Equal("test", r.Properties().Get("tag:env"))
}

func Test_Mock_ConfigServiceConfigurationAggregator_String(t *testing.T) {
	a := assert.New(t)
	r := &ConfigServiceConfigurationAggregator{
		Name: ptr.String("my-aggregator"),
	}
	a.Equal("my-aggregator", r.String())
}

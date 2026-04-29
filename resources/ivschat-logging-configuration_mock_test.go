package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"
	ivschattypes "github.com/aws/aws-sdk-go-v2/service/ivschat/types"
)

func Test_Mock_IVSChatLoggingConfiguration_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSChatClient)

	mockClient.On("ListLoggingConfigurations", mock.Anything, mock.Anything).
		Return(&ivschat.ListLoggingConfigurationsOutput{
			LoggingConfigurations: []ivschattypes.LoggingConfigurationSummary{
				{
					Arn:  ptr.String("arn:aws:ivschat:us-east-1:123456789012:logging-configuration/abc123"),
					Name: ptr.String("my-logging-config"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSChatLoggingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	loggingConfig := resources[0].(*IVSChatLoggingConfiguration)
	assertions.Equal("my-logging-config", *loggingConfig.Name)
	assertions.Equal("test", loggingConfig.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatLoggingConfiguration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSChatClient)

	mockClient.On("ListLoggingConfigurations", mock.Anything, mock.Anything).
		Return(&ivschat.ListLoggingConfigurationsOutput{
			LoggingConfigurations: []ivschattypes.LoggingConfigurationSummary{},
		}, nil)

	lister := &IVSChatLoggingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatLoggingConfiguration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSChatClient)

	loggingConfig := &IVSChatLoggingConfiguration{
		svc: mockClient,
		Arn: ptr.String("arn:aws:ivschat:us-east-1:123456789012:logging-configuration/abc123"),
	}

	mockClient.On("DeleteLoggingConfiguration", mock.Anything, &ivschat.DeleteLoggingConfigurationInput{
		Identifier: loggingConfig.Arn,
	}).Return(&ivschat.DeleteLoggingConfigurationOutput{}, nil)

	assertions.NoError(loggingConfig.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatLoggingConfiguration_Properties(t *testing.T) {
	assertions := assert.New(t)

	loggingConfig := IVSChatLoggingConfiguration{
		Arn:  ptr.String("arn:aws:ivschat:us-east-1:123456789012:logging-configuration/abc123"),
		Name: ptr.String("my-logging-config"),
		Tags: map[string]string{"env": "test"},
	}

	properties := loggingConfig.Properties()
	assertions.Equal("my-logging-config", properties.Get("Name"))
	assertions.Equal("test", properties.Get("tag:env"))
}

func Test_Mock_IVSChatLoggingConfiguration_String(t *testing.T) {
	assertions := assert.New(t)

	loggingConfig := IVSChatLoggingConfiguration{Name: ptr.String("my-logging-config")}
	assertions.Equal("my-logging-config", loggingConfig.String())
}

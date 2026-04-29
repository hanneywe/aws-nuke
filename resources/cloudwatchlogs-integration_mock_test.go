package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func Test_Mock_CloudWatchLogsIntegration_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("ListIntegrations", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.ListIntegrationsOutput{
			IntegrationSummaries: []cloudwatchlogstypes.IntegrationSummary{
				{
					IntegrationName: ptr.String("my-integration"),
					IntegrationType: cloudwatchlogstypes.IntegrationTypeOpensearch,
				},
			},
		}, nil)

	lister := &CloudWatchLogsIntegrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	integration := resources[0].(*CloudWatchLogsIntegration)
	assertions.Equal("my-integration", *integration.IntegrationName)
	assertions.Equal(cloudwatchlogstypes.IntegrationTypeOpensearch, integration.IntegrationType)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsIntegration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("ListIntegrations", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.ListIntegrationsOutput{
			IntegrationSummaries: []cloudwatchlogstypes.IntegrationSummary{},
		}, nil)

	lister := &CloudWatchLogsIntegrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsIntegration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	integration := &CloudWatchLogsIntegration{
		svc:             mockClient,
		IntegrationName: ptr.String("my-integration"),
	}

	mockClient.On("DeleteIntegration", mock.Anything, &cloudwatchlogs.DeleteIntegrationInput{
		IntegrationName: integration.IntegrationName,
		Force:           true,
	}).Return(&cloudwatchlogs.DeleteIntegrationOutput{}, nil)

	err := integration.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsIntegration_Properties(t *testing.T) {
	assertions := assert.New(t)

	integration := CloudWatchLogsIntegration{
		IntegrationName: ptr.String("my-integration"),
		IntegrationType: cloudwatchlogstypes.IntegrationTypeOpensearch,
	}

	properties := integration.Properties()
	assertions.Equal("my-integration", properties.Get("IntegrationName"))
	assertions.Equal(string(cloudwatchlogstypes.IntegrationTypeOpensearch), properties.Get("IntegrationType"))
}

func Test_Mock_CloudWatchLogsIntegration_String(t *testing.T) {
	assertions := assert.New(t)
	integration := CloudWatchLogsIntegration{IntegrationName: ptr.String("my-integration")}
	assertions.Equal("my-integration", integration.String())
}

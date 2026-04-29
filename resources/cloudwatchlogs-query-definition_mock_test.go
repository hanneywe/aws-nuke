package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testCloudWatchLogsV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_CloudWatchLogsQueryDefinition_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeQueryDefinitions", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeQueryDefinitionsOutput{
			QueryDefinitions: []cloudwatchlogstypes.QueryDefinition{
				{
					QueryDefinitionId: ptr.String("qd-12345"),
					Name:              ptr.String("my-query"),
				},
			},
		}, nil)

	lister := &CloudWatchLogsQueryDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	queryDef := resources[0].(*CloudWatchLogsQueryDefinition)
	assertions.Equal("qd-12345", *queryDef.QueryDefinitionID)
	assertions.Equal("my-query", *queryDef.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsQueryDefinition_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeQueryDefinitions", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeQueryDefinitionsOutput{
			QueryDefinitions: []cloudwatchlogstypes.QueryDefinition{},
		}, nil)

	lister := &CloudWatchLogsQueryDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsQueryDefinition_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	queryDef := &CloudWatchLogsQueryDefinition{
		svc:               mockClient,
		QueryDefinitionID: ptr.String("qd-12345"),
	}

	mockClient.On("DeleteQueryDefinition", mock.Anything, &cloudwatchlogs.DeleteQueryDefinitionInput{
		QueryDefinitionId: queryDef.QueryDefinitionID,
	}).Return(&cloudwatchlogs.DeleteQueryDefinitionOutput{}, nil)

	err := queryDef.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsQueryDefinition_Properties(t *testing.T) {
	assertions := assert.New(t)

	queryDef := CloudWatchLogsQueryDefinition{
		QueryDefinitionID: ptr.String("qd-12345"),
		Name:              ptr.String("my-query"),
	}

	properties := queryDef.Properties()
	assertions.Equal("qd-12345", properties.Get("QueryDefinitionId"))
	assertions.Equal("my-query", properties.Get("Name"))
}

func Test_Mock_CloudWatchLogsQueryDefinition_String(t *testing.T) {
	assertions := assert.New(t)
	queryDef := CloudWatchLogsQueryDefinition{Name: ptr.String("my-query")}
	assertions.Equal("my-query", queryDef.String())
}

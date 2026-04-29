package resources

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Mock_CloudWatchLogsBearerTokenAuthentication_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{
				{
					LogGroupName:                     ptr.String("/aws/test-group"),
					BearerTokenAuthenticationEnabled: ptr.Bool(true),
				},
				{
					LogGroupName:                     ptr.String("/aws/no-bearer"),
					BearerTokenAuthenticationEnabled: ptr.Bool(false),
				},
			},
		}, nil)

	lister := &CloudWatchLogsBearerTokenAuthenticationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CloudWatchLogsBearerTokenAuthentication)
	a.Equal("/aws/test-group", *r.LogGroupName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsBearerTokenAuthentication_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{},
		}, nil)

	lister := &CloudWatchLogsBearerTokenAuthenticationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsBearerTokenAuthentication_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	r := &CloudWatchLogsBearerTokenAuthentication{
		svc:          mockClient,
		LogGroupName: ptr.String("/aws/test-group"),
	}

	mockClient.On("PutBearerTokenAuthentication", mock.Anything, &cloudwatchlogs.PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               r.LogGroupName,
		BearerTokenAuthenticationEnabled: aws.Bool(false),
	}).Return(&cloudwatchlogs.PutBearerTokenAuthenticationOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsBearerTokenAuthentication_Properties(t *testing.T) {
	a := assert.New(t)

	r := CloudWatchLogsBearerTokenAuthentication{
		LogGroupName: ptr.String("/aws/test-group"),
	}

	props := r.Properties()
	a.Equal("/aws/test-group", props.Get("LogGroupName"))
}

func Test_Mock_CloudWatchLogsBearerTokenAuthentication_String(t *testing.T) {
	a := assert.New(t)
	r := CloudWatchLogsBearerTokenAuthentication{
		LogGroupName: ptr.String("/aws/test-group"),
	}
	a.Equal("/aws/test-group", r.String())
}

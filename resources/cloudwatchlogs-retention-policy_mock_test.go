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

func Test_Mock_CloudWatchLogsRetentionPolicy_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{
				{
					LogGroupName:    ptr.String("/aws/test-group"),
					RetentionInDays: ptr.Int32(30),
				},
				{
					LogGroupName:    ptr.String("/aws/no-retention"),
					RetentionInDays: nil,
				},
			},
		}, nil)

	lister := &CloudWatchLogsRetentionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	policy := resources[0].(*CloudWatchLogsRetentionPolicy)
	a.Equal("/aws/test-group", *policy.LogGroupName)
	a.Equal(int32(30), *policy.RetentionInDays)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsRetentionPolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{},
		}, nil)

	lister := &CloudWatchLogsRetentionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsRetentionPolicy_List_NoRetention(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{
				{
					LogGroupName:    ptr.String("/aws/group-1"),
					RetentionInDays: nil,
				},
				{
					LogGroupName:    ptr.String("/aws/group-2"),
					RetentionInDays: nil,
				},
			},
		}, nil)

	lister := &CloudWatchLogsRetentionPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsRetentionPolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	policy := &CloudWatchLogsRetentionPolicy{
		svc:          mockClient,
		LogGroupName: ptr.String("/aws/test-group"),
	}

	mockClient.On("DeleteRetentionPolicy", mock.Anything, &cloudwatchlogs.DeleteRetentionPolicyInput{
		LogGroupName: policy.LogGroupName,
	}).Return(&cloudwatchlogs.DeleteRetentionPolicyOutput{}, nil)

	err := policy.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsRetentionPolicy_Properties(t *testing.T) {
	a := assert.New(t)

	policy := CloudWatchLogsRetentionPolicy{
		LogGroupName:    ptr.String("/aws/test-group"),
		RetentionInDays: ptr.Int32(30),
	}

	props := policy.Properties()
	a.Equal("/aws/test-group", props.Get("LogGroupName"))
	a.Equal("30", props.Get("RetentionInDays"))
}

func Test_Mock_CloudWatchLogsRetentionPolicy_String(t *testing.T) {
	a := assert.New(t)
	policy := CloudWatchLogsRetentionPolicy{
		LogGroupName: ptr.String("/aws/test-group"),
	}
	a.Equal("/aws/test-group", policy.String())
}

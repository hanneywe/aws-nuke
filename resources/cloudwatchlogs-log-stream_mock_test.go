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

func Test_Mock_CloudWatchLogsLogStream_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{
				{
					LogGroupName: ptr.String("/aws/test-group"),
				},
			},
		}, nil)

	mockClient.On("DescribeLogStreams", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogStreamsOutput{
			LogStreams: []cloudwatchlogstypes.LogStream{
				{
					LogStreamName:      ptr.String("stream-1"),
					LastEventTimestamp: ptr.Int64(1700000000000),
				},
				{
					LogStreamName: ptr.String("stream-2"),
				},
			},
		}, nil)

	lister := &CloudWatchLogsLogStreamLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 2)

	stream := resources[0].(*CloudWatchLogsLogStream)
	a.Equal("/aws/test-group", *stream.LogGroupName)
	a.Equal("stream-1", *stream.LogStreamName)
	a.NotNil(stream.LastEventTimestamp)

	stream2 := resources[1].(*CloudWatchLogsLogStream)
	a.Equal("stream-2", *stream2.LogStreamName)
	a.Nil(stream2.LastEventTimestamp)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsLogStream_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{},
		}, nil)

	lister := &CloudWatchLogsLogStreamLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsLogStream_List_MultipleLogGroups(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	mockClient.On("DescribeLogGroups", mock.Anything, mock.Anything).
		Return(&cloudwatchlogs.DescribeLogGroupsOutput{
			LogGroups: []cloudwatchlogstypes.LogGroup{
				{LogGroupName: ptr.String("group-1")},
				{LogGroupName: ptr.String("group-2")},
			},
		}, nil)

	mockClient.On("DescribeLogStreams", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.DescribeLogStreamsInput) bool {
		return *input.LogGroupName == "group-1"
	})).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []cloudwatchlogstypes.LogStream{
			{LogStreamName: ptr.String("stream-a")},
		},
	}, nil)

	mockClient.On("DescribeLogStreams", mock.Anything, mock.MatchedBy(func(input *cloudwatchlogs.DescribeLogStreamsInput) bool {
		return *input.LogGroupName == "group-2"
	})).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []cloudwatchlogstypes.LogStream{
			{LogStreamName: ptr.String("stream-b")},
			{LogStreamName: ptr.String("stream-c")},
		},
	}, nil)

	lister := &CloudWatchLogsLogStreamLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 3)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsLogStream_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchLogsV2Client)

	stream := &CloudWatchLogsLogStream{
		svc:           mockClient,
		LogGroupName:  ptr.String("/aws/test-group"),
		LogStreamName: ptr.String("stream-1"),
	}

	mockClient.On("DeleteLogStream", mock.Anything, &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  stream.LogGroupName,
		LogStreamName: stream.LogStreamName,
	}).Return(&cloudwatchlogs.DeleteLogStreamOutput{}, nil)

	err := stream.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchLogsLogStream_Properties(t *testing.T) {
	a := assert.New(t)

	stream := CloudWatchLogsLogStream{
		LogGroupName:  ptr.String("/aws/test-group"),
		LogStreamName: ptr.String("stream-1"),
	}

	props := stream.Properties()
	a.Equal("/aws/test-group", props.Get("LogGroupName"))
	a.Equal("stream-1", props.Get("LogStreamName"))
}

func Test_Mock_CloudWatchLogsLogStream_String(t *testing.T) {
	a := assert.New(t)
	stream := CloudWatchLogsLogStream{
		LogGroupName:  ptr.String("/aws/test-group"),
		LogStreamName: ptr.String("stream-1"),
	}
	a.Equal("/aws/test-group/stream-1", stream.String())
}

package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2FlowLog_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeFlowLogs", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeFlowLogsOutput{
				FlowLogs: []ec2types.FlowLog{
					{
						FlowLogId:          ptr.String("fl-11111111111111111"),
						ResourceId:         ptr.String("vpc-22222222222222222"),
						TrafficType:        ec2types.TrafficTypeAll,
						LogDestinationType: ec2types.LogDestinationTypeCloudWatchLogs,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-flow-log")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2FlowLogLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	flowLog := resources[0].(*EC2FlowLog)
	assertions.Equal("fl-11111111111111111", *flowLog.FlowLogID)
	assertions.Equal("vpc-22222222222222222", *flowLog.ResourceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2FlowLog_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeFlowLogs", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeFlowLogsOutput{
				FlowLogs: []ec2types.FlowLog{},
			}, nil,
		)

	lister := &EC2FlowLogLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2FlowLog_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.FlowLog, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.FlowLog{
			FlowLogId:          ptr.String(fmt.Sprintf("fl-%d", i)),
			ResourceId:         ptr.String(fmt.Sprintf("vpc-%d", i)),
			TrafficType:        ec2types.TrafficTypeAll,
			LogDestinationType: ec2types.LogDestinationTypeCloudWatchLogs,
		}
	}

	mockClient.
		On(
			"DescribeFlowLogs",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeFlowLogsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeFlowLogsOutput{
				FlowLogs:  firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeFlowLogs",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeFlowLogsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeFlowLogsOutput{
				FlowLogs: []ec2types.FlowLog{
					{
						FlowLogId:          ptr.String("fl-100"),
						ResourceId:         ptr.String("vpc-100"),
						TrafficType:        ec2types.TrafficTypeAccept,
						LogDestinationType: ec2types.LogDestinationTypeS3,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2FlowLogLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2FlowLog_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	flowLog := &EC2FlowLog{
		svc:       mockClient,
		FlowLogID: ptr.String("fl-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteFlowLogs",
			mock.Anything,
			&ec2.DeleteFlowLogsInput{
				FlowLogIds: []string{"fl-11111111111111111"},
			},
		).
		Return(&ec2.DeleteFlowLogsOutput{}, nil)

	err := flowLog.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2FlowLog_Properties(t *testing.T) {
	assertions := assert.New(t)

	flowLog := EC2FlowLog{
		FlowLogID:          ptr.String("fl-11111111111111111"),
		ResourceID:         ptr.String("vpc-22222222222222222"),
		TrafficType:        "ALL",
		LogDestinationType: "cloud-watch-logs",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := flowLog.Properties()

	assertions.Equal("fl-11111111111111111", properties.Get("FlowLogId"))
	assertions.Equal("vpc-22222222222222222", properties.Get("ResourceId"))
	assertions.Equal("ALL", properties.Get("TrafficType"))
	assertions.Equal("cloud-watch-logs", properties.Get("LogDestinationType"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2FlowLog_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	flowLog := EC2FlowLog{
		FlowLogID:          ptr.String("fl-99999999999999999"),
		ResourceID:         ptr.String("vpc-88888888888888888"),
		TrafficType:        "REJECT",
		LogDestinationType: "s3",
		Tags:               []ec2types.Tag{},
	}

	properties := flowLog.Properties()

	assertions.Equal("fl-99999999999999999", properties.Get("FlowLogId"))
	assertions.Equal("REJECT", properties.Get("TrafficType"))
}

func Test_Mock_EC2FlowLog_String(t *testing.T) {
	assertions := assert.New(t)

	flowLog := EC2FlowLog{
		FlowLogID: ptr.String("fl-11111111111111111"),
	}

	assertions.Equal("fl-11111111111111111", flowLog.String())
}

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

func Test_Mock_EC2TrafficMirrorSession_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorSessions", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorSessionsOutput{
				TrafficMirrorSessions: []ec2types.TrafficMirrorSession{
					{
						TrafficMirrorSessionId: ptr.String("tms-11111111111111111"),
						TrafficMirrorTargetId:  ptr.String("tmt-22222222222222222"),
						TrafficMirrorFilterId:  ptr.String("tmf-33333333333333333"),
						NetworkInterfaceId:     ptr.String("eni-44444444444444444"),
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-session")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TrafficMirrorSessionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	session := resources[0].(*EC2TrafficMirrorSession)
	assertions.Equal("tms-11111111111111111", *session.TrafficMirrorSessionID)
	assertions.Equal("tmt-22222222222222222", *session.TrafficMirrorTargetID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorSession_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorSessions", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorSessionsOutput{
				TrafficMirrorSessions: []ec2types.TrafficMirrorSession{},
			}, nil,
		)

	lister := &EC2TrafficMirrorSessionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorSession_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TrafficMirrorSession, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TrafficMirrorSession{
			TrafficMirrorSessionId: ptr.String(fmt.Sprintf("tms-%d", i)),
			TrafficMirrorTargetId:  ptr.String(fmt.Sprintf("tmt-%d", i)),
			TrafficMirrorFilterId:  ptr.String(fmt.Sprintf("tmf-%d", i)),
			NetworkInterfaceId:     ptr.String(fmt.Sprintf("eni-%d", i)),
		}
	}

	mockClient.
		On(
			"DescribeTrafficMirrorSessions",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorSessionsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorSessionsOutput{
				TrafficMirrorSessions: firstPageItems,
				NextToken:             ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTrafficMirrorSessions",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorSessionsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorSessionsOutput{
				TrafficMirrorSessions: []ec2types.TrafficMirrorSession{
					{
						TrafficMirrorSessionId: ptr.String("tms-100"),
						TrafficMirrorTargetId:  ptr.String("tmt-100"),
						TrafficMirrorFilterId:  ptr.String("tmf-100"),
						NetworkInterfaceId:     ptr.String("eni-100"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TrafficMirrorSessionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorSession_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	session := &EC2TrafficMirrorSession{
		svc:                    mockClient,
		TrafficMirrorSessionID: ptr.String("tms-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTrafficMirrorSession",
			mock.Anything,
			&ec2.DeleteTrafficMirrorSessionInput{
				TrafficMirrorSessionId: session.TrafficMirrorSessionID,
			},
		).
		Return(&ec2.DeleteTrafficMirrorSessionOutput{}, nil)

	err := session.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorSession_Properties(t *testing.T) {
	assertions := assert.New(t)

	session := EC2TrafficMirrorSession{
		TrafficMirrorSessionID: ptr.String("tms-11111111111111111"),
		TrafficMirrorTargetID:  ptr.String("tmt-22222222222222222"),
		TrafficMirrorFilterID:  ptr.String("tmf-33333333333333333"),
		NetworkInterfaceID:     ptr.String("eni-44444444444444444"),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := session.Properties()

	assertions.Equal("tms-11111111111111111", properties.Get("TrafficMirrorSessionId"))
	assertions.Equal("tmt-22222222222222222", properties.Get("TrafficMirrorTargetId"))
	assertions.Equal("tmf-33333333333333333", properties.Get("TrafficMirrorFilterId"))
	assertions.Equal("eni-44444444444444444", properties.Get("NetworkInterfaceId"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TrafficMirrorSession_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	session := EC2TrafficMirrorSession{
		TrafficMirrorSessionID: ptr.String("tms-99999999999999999"),
		TrafficMirrorTargetID:  ptr.String("tmt-88888888888888888"),
		TrafficMirrorFilterID:  ptr.String("tmf-77777777777777777"),
		NetworkInterfaceID:     ptr.String("eni-66666666666666666"),
		Tags:                   []ec2types.Tag{},
	}

	properties := session.Properties()

	assertions.Equal("tms-99999999999999999", properties.Get("TrafficMirrorSessionId"))
	assertions.Equal("tmt-88888888888888888", properties.Get("TrafficMirrorTargetId"))
}

func Test_Mock_EC2TrafficMirrorSession_String(t *testing.T) {
	assertions := assert.New(t)

	session := EC2TrafficMirrorSession{
		TrafficMirrorSessionID: ptr.String("tms-11111111111111111"),
	}

	assertions.Equal("tms-11111111111111111", session.String())
}

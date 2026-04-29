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

func Test_Mock_EC2TrafficMirrorTarget_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorTargets", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorTargetsOutput{
				TrafficMirrorTargets: []ec2types.TrafficMirrorTarget{
					{
						TrafficMirrorTargetId:  ptr.String("tmt-11111111111111111"),
						Type:                   ec2types.TrafficMirrorTargetTypeNetworkInterface,
						NetworkInterfaceId:     ptr.String("eni-22222222222222222"),
						NetworkLoadBalancerArn: nil,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-target")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TrafficMirrorTargetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	target := resources[0].(*EC2TrafficMirrorTarget)
	assertions.Equal("tmt-11111111111111111", *target.TrafficMirrorTargetID)
	assertions.Equal("eni-22222222222222222", *target.NetworkInterfaceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorTarget_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorTargets", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTrafficMirrorTargetsOutput{
				TrafficMirrorTargets: []ec2types.TrafficMirrorTarget{},
			}, nil,
		)

	lister := &EC2TrafficMirrorTargetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorTarget_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TrafficMirrorTarget, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TrafficMirrorTarget{
			TrafficMirrorTargetId: ptr.String(fmt.Sprintf("tmt-%d", i)),
			Type:                  ec2types.TrafficMirrorTargetTypeNetworkInterface,
			NetworkInterfaceId:    ptr.String(fmt.Sprintf("eni-%d", i)),
		}
	}

	mockClient.
		On(
			"DescribeTrafficMirrorTargets",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorTargetsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorTargetsOutput{
				TrafficMirrorTargets: firstPageItems,
				NextToken:            ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTrafficMirrorTargets",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorTargetsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTrafficMirrorTargetsOutput{
				TrafficMirrorTargets: []ec2types.TrafficMirrorTarget{
					{
						TrafficMirrorTargetId:  ptr.String("tmt-100"),
						Type:                   ec2types.TrafficMirrorTargetTypeNetworkLoadBalancer,
						NetworkLoadBalancerArn: ptr.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/net/my-nlb/1234567890"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TrafficMirrorTargetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorTarget_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	target := &EC2TrafficMirrorTarget{
		svc:                   mockClient,
		TrafficMirrorTargetID: ptr.String("tmt-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTrafficMirrorTarget",
			mock.Anything,
			&ec2.DeleteTrafficMirrorTargetInput{
				TrafficMirrorTargetId: target.TrafficMirrorTargetID,
			},
		).
		Return(&ec2.DeleteTrafficMirrorTargetOutput{}, nil)

	err := target.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorTarget_Properties(t *testing.T) {
	assertions := assert.New(t)

	target := EC2TrafficMirrorTarget{
		TrafficMirrorTargetID:  ptr.String("tmt-11111111111111111"),
		Type:                   "network-interface",
		NetworkInterfaceID:     ptr.String("eni-22222222222222222"),
		NetworkLoadBalancerArn: ptr.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/net/my-nlb/1234567890"),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := target.Properties()

	assertions.Equal("tmt-11111111111111111", properties.Get("TrafficMirrorTargetId"))
	assertions.Equal("network-interface", properties.Get("Type"))
	assertions.Equal("eni-22222222222222222", properties.Get("NetworkInterfaceId"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TrafficMirrorTarget_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	target := EC2TrafficMirrorTarget{
		TrafficMirrorTargetID: ptr.String("tmt-99999999999999999"),
		Type:                  "network-load-balancer",
		Tags:                  []ec2types.Tag{},
	}

	properties := target.Properties()

	assertions.Equal("tmt-99999999999999999", properties.Get("TrafficMirrorTargetId"))
	assertions.Equal("network-load-balancer", properties.Get("Type"))
}

func Test_Mock_EC2TrafficMirrorTarget_String(t *testing.T) {
	assertions := assert.New(t)

	target := EC2TrafficMirrorTarget{
		TrafficMirrorTargetID: ptr.String("tmt-11111111111111111"),
	}

	assertions.Equal("tmt-11111111111111111", target.String())
}

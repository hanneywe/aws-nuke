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

func Test_Mock_EC2NetworkInsightsPath_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsPaths", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsPathsOutput{
				NetworkInsightsPaths: []ec2types.NetworkInsightsPath{
					{
						NetworkInsightsPathId: ptr.String("nip-11111111111111111"),
						Source:                ptr.String("eni-aaaaaaaaaaaaaaaa"),
						Destination:           ptr.String("eni-bbbbbbbbbbbbbbbb"),
						Protocol:              ec2types.ProtocolTcp,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-path")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2NetworkInsightsPathLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	insightsPath := resources[0].(*EC2NetworkInsightsPath)
	assertions.Equal("nip-11111111111111111", *insightsPath.NetworkInsightsPathID)
	assertions.Equal("eni-aaaaaaaaaaaaaaaa", *insightsPath.Source)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsPath_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsPaths", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsPathsOutput{
				NetworkInsightsPaths: []ec2types.NetworkInsightsPath{},
			}, nil,
		)

	lister := &EC2NetworkInsightsPathLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsPath_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.NetworkInsightsPath, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.NetworkInsightsPath{
			NetworkInsightsPathId: ptr.String(fmt.Sprintf("nip-%d", i)),
			Source:                ptr.String(fmt.Sprintf("eni-src-%d", i)),
			Destination:           ptr.String(fmt.Sprintf("eni-dst-%d", i)),
			Protocol:              ec2types.ProtocolTcp,
		}
	}

	mockClient.
		On(
			"DescribeNetworkInsightsPaths",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsPathsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsPathsOutput{
				NetworkInsightsPaths: firstPageItems,
				NextToken:            ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeNetworkInsightsPaths",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsPathsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsPathsOutput{
				NetworkInsightsPaths: []ec2types.NetworkInsightsPath{
					{
						NetworkInsightsPathId: ptr.String("nip-100"),
						Source:                ptr.String("eni-src-100"),
						Destination:           ptr.String("eni-dst-100"),
						Protocol:              ec2types.ProtocolUdp,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2NetworkInsightsPathLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsPath_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	insightsPath := &EC2NetworkInsightsPath{
		svc:                   mockClient,
		NetworkInsightsPathID: ptr.String("nip-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteNetworkInsightsPath",
			mock.Anything,
			&ec2.DeleteNetworkInsightsPathInput{
				NetworkInsightsPathId: insightsPath.NetworkInsightsPathID,
			},
		).
		Return(&ec2.DeleteNetworkInsightsPathOutput{}, nil)

	err := insightsPath.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsPath_Properties(t *testing.T) {
	assertions := assert.New(t)

	insightsPath := EC2NetworkInsightsPath{
		NetworkInsightsPathID: ptr.String("nip-11111111111111111"),
		Source:                ptr.String("eni-aaaaaaaaaaaaaaaa"),
		Destination:           ptr.String("eni-bbbbbbbbbbbbbbbb"),
		Protocol:              "tcp",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("staging")},
		},
	}

	properties := insightsPath.Properties()

	assertions.Equal("nip-11111111111111111", properties.Get("NetworkInsightsPathId"))
	assertions.Equal("eni-aaaaaaaaaaaaaaaa", properties.Get("Source"))
	assertions.Equal("eni-bbbbbbbbbbbbbbbb", properties.Get("Destination"))
	assertions.Equal("tcp", properties.Get("Protocol"))
	assertions.Equal("staging", properties.Get("tag:Environment"))
}

func Test_Mock_EC2NetworkInsightsPath_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	insightsPath := EC2NetworkInsightsPath{
		NetworkInsightsPathID: ptr.String("nip-99999999999999999"),
		Source:                ptr.String("eni-src-empty"),
		Destination:           ptr.String("eni-dst-empty"),
		Protocol:              "udp",
		Tags:                  []ec2types.Tag{},
	}

	properties := insightsPath.Properties()

	assertions.Equal("nip-99999999999999999", properties.Get("NetworkInsightsPathId"))
	assertions.Equal("udp", properties.Get("Protocol"))
}

func Test_Mock_EC2NetworkInsightsPath_String(t *testing.T) {
	assertions := assert.New(t)

	insightsPath := EC2NetworkInsightsPath{
		NetworkInsightsPathID: ptr.String("nip-11111111111111111"),
	}

	assertions.Equal("nip-11111111111111111", insightsPath.String())
}

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

func Test_Mock_EC2NetworkInsightsAnalysis_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsAnalyses", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsAnalysesOutput{
				NetworkInsightsAnalyses: []ec2types.NetworkInsightsAnalysis{
					{
						NetworkInsightsAnalysisId: ptr.String("nia-11111111111111111"),
						NetworkInsightsPathId:     ptr.String("nip-22222222222222222"),
						Status:                    ec2types.AnalysisStatusSucceeded,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-analysis")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2NetworkInsightsAnalysisLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	analysis := resources[0].(*EC2NetworkInsightsAnalysis)
	assertions.Equal("nia-11111111111111111", *analysis.NetworkInsightsAnalysisID)
	assertions.Equal("nip-22222222222222222", *analysis.NetworkInsightsPathID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAnalysis_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsAnalyses", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsAnalysesOutput{
				NetworkInsightsAnalyses: []ec2types.NetworkInsightsAnalysis{},
			}, nil,
		)

	lister := &EC2NetworkInsightsAnalysisLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAnalysis_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.NetworkInsightsAnalysis, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.NetworkInsightsAnalysis{
			NetworkInsightsAnalysisId: ptr.String(fmt.Sprintf("nia-%d", i)),
			NetworkInsightsPathId:     ptr.String(fmt.Sprintf("nip-%d", i)),
			Status:                    ec2types.AnalysisStatusSucceeded,
		}
	}

	mockClient.
		On(
			"DescribeNetworkInsightsAnalyses",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsAnalysesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsAnalysesOutput{
				NetworkInsightsAnalyses: firstPageItems,
				NextToken:               ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeNetworkInsightsAnalyses",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsAnalysesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsAnalysesOutput{
				NetworkInsightsAnalyses: []ec2types.NetworkInsightsAnalysis{
					{
						NetworkInsightsAnalysisId: ptr.String("nia-100"),
						NetworkInsightsPathId:     ptr.String("nip-100"),
						Status:                    ec2types.AnalysisStatusRunning,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2NetworkInsightsAnalysisLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAnalysis_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	analysis := &EC2NetworkInsightsAnalysis{
		svc:                       mockClient,
		NetworkInsightsAnalysisID: ptr.String("nia-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteNetworkInsightsAnalysis",
			mock.Anything,
			&ec2.DeleteNetworkInsightsAnalysisInput{
				NetworkInsightsAnalysisId: analysis.NetworkInsightsAnalysisID,
			},
		).
		Return(&ec2.DeleteNetworkInsightsAnalysisOutput{}, nil)

	err := analysis.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAnalysis_Properties(t *testing.T) {
	assertions := assert.New(t)

	analysis := EC2NetworkInsightsAnalysis{
		NetworkInsightsAnalysisID: ptr.String("nia-11111111111111111"),
		NetworkInsightsPathID:     ptr.String("nip-22222222222222222"),
		Status:                    "succeeded",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := analysis.Properties()

	assertions.Equal("nia-11111111111111111", properties.Get("NetworkInsightsAnalysisId"))
	assertions.Equal("nip-22222222222222222", properties.Get("NetworkInsightsPathId"))
	assertions.Equal("succeeded", properties.Get("Status"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2NetworkInsightsAnalysis_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	analysis := EC2NetworkInsightsAnalysis{
		NetworkInsightsAnalysisID: ptr.String("nia-99999999999999999"),
		NetworkInsightsPathID:     ptr.String("nip-88888888888888888"),
		Status:                    "running",
		Tags:                      []ec2types.Tag{},
	}

	properties := analysis.Properties()

	assertions.Equal("nia-99999999999999999", properties.Get("NetworkInsightsAnalysisId"))
	assertions.Equal("running", properties.Get("Status"))
}

func Test_Mock_EC2NetworkInsightsAnalysis_String(t *testing.T) {
	assertions := assert.New(t)

	analysis := EC2NetworkInsightsAnalysis{
		NetworkInsightsAnalysisID: ptr.String("nia-11111111111111111"),
	}

	assertions.Equal("nia-11111111111111111", analysis.String())
}

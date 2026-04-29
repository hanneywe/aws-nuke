package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Mock_EC2NetworkInsightsAccessScope_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsAccessScopes", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsAccessScopesOutput{
				NetworkInsightsAccessScopes: []ec2types.NetworkInsightsAccessScope{
					{
						NetworkInsightsAccessScopeId:  ptr.String("nis-11111111111111111"),
						NetworkInsightsAccessScopeArn: ptr.String("arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-11111111111111111"),
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-scope")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2NetworkInsightsAccessScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	scope := resources[0].(*EC2NetworkInsightsAccessScope)
	assertions.Equal("nis-11111111111111111", *scope.NetworkInsightsAccessScopeID)
	expectedArn := "arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-11111111111111111"
	assertions.Equal(expectedArn, *scope.NetworkInsightsAccessScopeArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAccessScope_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeNetworkInsightsAccessScopes", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeNetworkInsightsAccessScopesOutput{
				NetworkInsightsAccessScopes: []ec2types.NetworkInsightsAccessScope{},
			}, nil,
		)

	lister := &EC2NetworkInsightsAccessScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAccessScope_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.NetworkInsightsAccessScope, 100)
	for index := range firstPageItems {
		firstPageItems[index] = ec2types.NetworkInsightsAccessScope{
			NetworkInsightsAccessScopeId:  ptr.String(fmt.Sprintf("nis-%d", index)),
			NetworkInsightsAccessScopeArn: ptr.String(fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-%d", index)),
		}
	}

	mockClient.
		On(
			"DescribeNetworkInsightsAccessScopes",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsAccessScopesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsAccessScopesOutput{
				NetworkInsightsAccessScopes: firstPageItems,
				NextToken:                   ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeNetworkInsightsAccessScopes",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeNetworkInsightsAccessScopesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeNetworkInsightsAccessScopesOutput{
				NetworkInsightsAccessScopes: []ec2types.NetworkInsightsAccessScope{
					{
						NetworkInsightsAccessScopeId:  ptr.String("nis-100"),
						NetworkInsightsAccessScopeArn: ptr.String("arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-100"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2NetworkInsightsAccessScopeLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAccessScope_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	scope := &EC2NetworkInsightsAccessScope{
		svc:                          mockClient,
		NetworkInsightsAccessScopeID: ptr.String("nis-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteNetworkInsightsAccessScope",
			mock.Anything,
			&ec2.DeleteNetworkInsightsAccessScopeInput{
				NetworkInsightsAccessScopeId: scope.NetworkInsightsAccessScopeID,
			},
		).
		Return(&ec2.DeleteNetworkInsightsAccessScopeOutput{}, nil)

	err := scope.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2NetworkInsightsAccessScope_Properties(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2NetworkInsightsAccessScope{
		NetworkInsightsAccessScopeID:  ptr.String("nis-11111111111111111"),
		NetworkInsightsAccessScopeArn: ptr.String("arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-11111111111111111"),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := scope.Properties()

	assertions.Equal("nis-11111111111111111", properties.Get("NetworkInsightsAccessScopeId"))
	expectedArn := "arn:aws:ec2:us-east-1:123456789012:network-insights-access-scope/nis-11111111111111111"
	assertions.Equal(expectedArn, properties.Get("NetworkInsightsAccessScopeArn"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2NetworkInsightsAccessScope_String(t *testing.T) {
	assertions := assert.New(t)

	scope := EC2NetworkInsightsAccessScope{
		NetworkInsightsAccessScopeID: ptr.String("nis-11111111111111111"),
	}

	assertions.Equal("nis-11111111111111111", scope.String())
}

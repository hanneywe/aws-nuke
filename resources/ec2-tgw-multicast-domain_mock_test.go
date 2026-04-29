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

// --- Listing ---

func Test_Mock_EC2TGWMulticastDomain_ListWithOneDomain(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayMulticastDomains", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayMulticastDomainsOutput{
				TransitGatewayMulticastDomains: []ec2types.TransitGatewayMulticastDomain{
					{
						TransitGatewayMulticastDomainId: ptr.String("tgw-mcast-domain-11111111111111111"),
						TransitGatewayId:                ptr.String("tgw-11111111111111111"),
						State:                           ec2types.TransitGatewayMulticastDomainStateAvailable,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-domain")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TGWMulticastDomainLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	domain := resources[0].(*EC2TGWMulticastDomain)
	assertions.Equal("tgw-mcast-domain-11111111111111111", *domain.TransitGatewayMulticastDomainID)
	assertions.Equal("tgw-11111111111111111", *domain.TransitGatewayID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWMulticastDomain_ListWithNoDomains(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayMulticastDomains", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayMulticastDomainsOutput{
				TransitGatewayMulticastDomains: []ec2types.TransitGatewayMulticastDomain{},
			}, nil,
		)

	lister := &EC2TGWMulticastDomainLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWMulticastDomain_ListWithMultiplePages(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TransitGatewayMulticastDomain, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TransitGatewayMulticastDomain{
			TransitGatewayMulticastDomainId: ptr.String(fmt.Sprintf("tgw-mcast-domain-%017d", i)),
			TransitGatewayId:                ptr.String("tgw-11111111111111111"),
			State:                           ec2types.TransitGatewayMulticastDomainStateAvailable,
		}
	}

	mockClient.
		On(
			"DescribeTransitGatewayMulticastDomains",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayMulticastDomainsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayMulticastDomainsOutput{
				TransitGatewayMulticastDomains: firstPageItems,
				NextToken:                      ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTransitGatewayMulticastDomains",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayMulticastDomainsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayMulticastDomainsOutput{
				TransitGatewayMulticastDomains: []ec2types.TransitGatewayMulticastDomain{
					{
						TransitGatewayMulticastDomainId: ptr.String("tgw-mcast-domain-00000000000000100"),
						TransitGatewayId:                ptr.String("tgw-11111111111111111"),
						State:                           ec2types.TransitGatewayMulticastDomainStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TGWMulticastDomainLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_EC2TGWMulticastDomain_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	domain := &EC2TGWMulticastDomain{
		svc:                             mockClient,
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTransitGatewayMulticastDomain",
			mock.Anything,
			&ec2.DeleteTransitGatewayMulticastDomainInput{
				TransitGatewayMulticastDomainId: domain.TransitGatewayMulticastDomainID,
			},
		).
		Return(&ec2.DeleteTransitGatewayMulticastDomainOutput{}, nil)

	err := domain.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_EC2TGWMulticastDomain_Properties(t *testing.T) {
	assertions := assert.New(t)

	domain := EC2TGWMulticastDomain{
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-11111111111111111"),
		TransitGatewayID:                ptr.String("tgw-11111111111111111"),
		State:                           "available",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("staging")},
		},
	}

	properties := domain.Properties()

	assertions.Equal("tgw-mcast-domain-11111111111111111", properties.Get("TransitGatewayMulticastDomainId"))
	assertions.Equal("tgw-11111111111111111", properties.Get("TransitGatewayId"))
	assertions.Equal("available", properties.Get("State"))
	assertions.Equal("staging", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TGWMulticastDomain_PropertiesWithEmptyTags(t *testing.T) {
	assertions := assert.New(t)

	domain := EC2TGWMulticastDomain{
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-22222222222222222"),
		TransitGatewayID:                ptr.String("tgw-22222222222222222"),
		State:                           "available",
		Tags:                            []ec2types.Tag{},
	}

	properties := domain.Properties()

	assertions.Equal("tgw-mcast-domain-22222222222222222", properties.Get("TransitGatewayMulticastDomainId"))
}

// --- Display ---

func Test_Mock_EC2TGWMulticastDomain_String(t *testing.T) {
	assertions := assert.New(t)

	domain := EC2TGWMulticastDomain{
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-11111111111111111"),
	}

	assertions.Equal("tgw-mcast-domain-11111111111111111", domain.String())
}

// --- Filter ---

func Test_Mock_EC2TGWMulticastDomain_FilterExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	domain := EC2TGWMulticastDomain{
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-11111111111111111"),
		State:                           string(ec2types.TransitGatewayMulticastDomainStateDeleted),
	}

	err := domain.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2TGWMulticastDomain_FilterPassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	domain := EC2TGWMulticastDomain{
		TransitGatewayMulticastDomainID: ptr.String("tgw-mcast-domain-11111111111111111"),
		State:                           string(ec2types.TransitGatewayMulticastDomainStateAvailable),
	}

	err := domain.Filter()
	assertions.NoError(err)
}

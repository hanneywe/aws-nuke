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

func Test_Mock_EC2TGWRouteTable_ListWithOneRouteTable(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayRouteTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: []ec2types.TransitGatewayRouteTable{
					{
						TransitGatewayRouteTableId:   ptr.String("tgw-rtb-11111111111111111"),
						TransitGatewayId:             ptr.String("tgw-11111111111111111"),
						State:                        ec2types.TransitGatewayRouteTableStateAvailable,
						DefaultAssociationRouteTable: ptr.Bool(false),
						DefaultPropagationRouteTable: ptr.Bool(false),
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-route-table")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TGWRouteTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	routeTable := resources[0].(*EC2TGWRouteTable)
	assertions.Equal("tgw-rtb-11111111111111111", *routeTable.TransitGatewayRouteTableID)
	assertions.Equal("tgw-11111111111111111", *routeTable.TransitGatewayID)
	assertions.Equal("available", routeTable.State)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWRouteTable_ListWithNoRouteTables(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayRouteTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: []ec2types.TransitGatewayRouteTable{},
			}, nil,
		)

	lister := &EC2TGWRouteTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWRouteTable_ListWithMultiplePages(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TransitGatewayRouteTable, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TransitGatewayRouteTable{
			TransitGatewayRouteTableId: ptr.String(fmt.Sprintf("tgw-rtb-%017d", i)),
			TransitGatewayId:           ptr.String("tgw-11111111111111111"),
			State:                      ec2types.TransitGatewayRouteTableStateAvailable,
		}
	}

	mockClient.
		On(
			"DescribeTransitGatewayRouteTables",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayRouteTablesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: firstPageItems,
				NextToken:                 ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTransitGatewayRouteTables",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayRouteTablesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: []ec2types.TransitGatewayRouteTable{
					{
						TransitGatewayRouteTableId: ptr.String("tgw-rtb-00000000000000100"),
						TransitGatewayId:           ptr.String("tgw-11111111111111111"),
						State:                      ec2types.TransitGatewayRouteTableStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TGWRouteTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_EC2TGWRouteTable_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	routeTable := &EC2TGWRouteTable{
		svc:                        mockClient,
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTransitGatewayRouteTable",
			mock.Anything,
			&ec2.DeleteTransitGatewayRouteTableInput{
				TransitGatewayRouteTableId: routeTable.TransitGatewayRouteTableID,
			},
		).
		Return(&ec2.DeleteTransitGatewayRouteTableOutput{}, nil)

	err := routeTable.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_EC2TGWRouteTable_Properties(t *testing.T) {
	assertions := assert.New(t)

	routeTable := EC2TGWRouteTable{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		TransitGatewayID:           ptr.String("tgw-11111111111111111"),
		State:                      "available",
		DefaultAssociation:         ptr.Bool(true),
		DefaultPropagation:         ptr.Bool(false),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := routeTable.Properties()

	assertions.Equal("tgw-rtb-11111111111111111", properties.Get("TransitGatewayRouteTableId"))
	assertions.Equal("tgw-11111111111111111", properties.Get("TransitGatewayId"))
	assertions.Equal("available", properties.Get("State"))
	assertions.Equal("true", properties.Get("DefaultAssociation"))
	assertions.Equal("false", properties.Get("DefaultPropagation"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TGWRouteTable_PropertiesWithEmptyTags(t *testing.T) {
	assertions := assert.New(t)

	routeTable := EC2TGWRouteTable{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-22222222222222222"),
		TransitGatewayID:           ptr.String("tgw-22222222222222222"),
		State:                      "available",
		Tags:                       []ec2types.Tag{},
	}

	properties := routeTable.Properties()

	assertions.Equal("tgw-rtb-22222222222222222", properties.Get("TransitGatewayRouteTableId"))
	assertions.Equal("tgw-22222222222222222", properties.Get("TransitGatewayId"))
}

// --- Display ---

func Test_Mock_EC2TGWRouteTable_String(t *testing.T) {
	assertions := assert.New(t)

	routeTable := EC2TGWRouteTable{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
	}

	assertions.Equal("tgw-rtb-11111111111111111", routeTable.String())
}

// --- Filter ---

func Test_Mock_EC2TGWRouteTable_FilterExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	routeTable := EC2TGWRouteTable{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		State:                      string(ec2types.TransitGatewayRouteTableStateDeleted),
	}

	err := routeTable.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2TGWRouteTable_FilterPassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	routeTable := EC2TGWRouteTable{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		State:                      string(ec2types.TransitGatewayRouteTableStateAvailable),
	}

	err := routeTable.Filter()
	assertions.NoError(err)
}

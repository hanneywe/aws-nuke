package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// --- Listing ---

func Test_Mock_EC2TGWPrefixListReference_ListWithOneReference(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayRouteTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: []ec2types.TransitGatewayRouteTable{
					{
						TransitGatewayRouteTableId: ptr.String("tgw-rtb-11111111111111111"),
					},
				},
			}, nil,
		)

	mockClient.
		On("GetTransitGatewayPrefixListReferences", mock.Anything, mock.Anything).
		Return(
			&ec2.GetTransitGatewayPrefixListReferencesOutput{
				TransitGatewayPrefixListReferences: []ec2types.TransitGatewayPrefixListReference{
					{
						PrefixListId:               ptr.String("pl-11111111111111111"),
						TransitGatewayRouteTableId: ptr.String("tgw-rtb-11111111111111111"),
						State:                      ec2types.TransitGatewayPrefixListReferenceStateAvailable,
					},
				},
			}, nil,
		)

	lister := &EC2TGWPrefixListReferenceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	reference := resources[0].(*EC2TGWPrefixListReference)
	assertions.Equal("pl-11111111111111111", *reference.PrefixListID)
	assertions.Equal("tgw-rtb-11111111111111111", *reference.TransitGatewayRouteTableID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWPrefixListReference_ListWithNoReferences(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayRouteTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayRouteTablesOutput{
				TransitGatewayRouteTables: []ec2types.TransitGatewayRouteTable{},
			}, nil,
		)

	lister := &EC2TGWPrefixListReferenceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_EC2TGWPrefixListReference_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	reference := &EC2TGWPrefixListReference{
		svc:                        mockClient,
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		PrefixListID:               ptr.String("pl-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTransitGatewayPrefixListReference",
			mock.Anything,
			&ec2.DeleteTransitGatewayPrefixListReferenceInput{
				TransitGatewayRouteTableId: reference.TransitGatewayRouteTableID,
				PrefixListId:               reference.PrefixListID,
			},
		).
		Return(&ec2.DeleteTransitGatewayPrefixListReferenceOutput{}, nil)

	err := reference.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_EC2TGWPrefixListReference_Properties(t *testing.T) {
	assertions := assert.New(t)

	reference := EC2TGWPrefixListReference{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		PrefixListID:               ptr.String("pl-11111111111111111"),
		State:                      "available",
	}

	properties := reference.Properties()

	assertions.Equal("tgw-rtb-11111111111111111", properties.Get("TransitGatewayRouteTableId"))
	assertions.Equal("pl-11111111111111111", properties.Get("PrefixListId"))
	assertions.Equal("available", properties.Get("State"))
}

// --- Display ---

func Test_Mock_EC2TGWPrefixListReference_String(t *testing.T) {
	assertions := assert.New(t)

	reference := EC2TGWPrefixListReference{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		PrefixListID:               ptr.String("pl-11111111111111111"),
	}

	assertions.Equal("tgw-rtb-11111111111111111 -> pl-11111111111111111", reference.String())
}

// --- Filter ---

func Test_Mock_EC2TGWPrefixListReference_FilterExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	reference := EC2TGWPrefixListReference{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		PrefixListID:               ptr.String("pl-11111111111111111"),
		State:                      string(ec2types.TransitGatewayPrefixListReferenceStateDeleting),
	}

	err := reference.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already being deleted")
}

func Test_Mock_EC2TGWPrefixListReference_FilterPassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	reference := EC2TGWPrefixListReference{
		TransitGatewayRouteTableID: ptr.String("tgw-rtb-11111111111111111"),
		PrefixListID:               ptr.String("pl-11111111111111111"),
		State:                      string(ec2types.TransitGatewayPrefixListReferenceStateAvailable),
	}

	err := reference.Filter()
	assertions.NoError(err)
}

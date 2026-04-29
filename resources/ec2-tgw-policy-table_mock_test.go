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

func Test_Mock_EC2TGWPolicyTable_ListWithOnePolicyTable(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayPolicyTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayPolicyTablesOutput{
				TransitGatewayPolicyTables: []ec2types.TransitGatewayPolicyTable{
					{
						TransitGatewayPolicyTableId: ptr.String("tgw-ptb-11111111111111111"),
						TransitGatewayId:            ptr.String("tgw-11111111111111111"),
						State:                       ec2types.TransitGatewayPolicyTableStateAvailable,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-policy-table")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2TGWPolicyTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	policyTable := resources[0].(*EC2TGWPolicyTable)
	assertions.Equal("tgw-ptb-11111111111111111", *policyTable.TransitGatewayPolicyTableID)
	assertions.Equal("tgw-11111111111111111", *policyTable.TransitGatewayID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWPolicyTable_ListWithNoPolicyTables(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTransitGatewayPolicyTables", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeTransitGatewayPolicyTablesOutput{
				TransitGatewayPolicyTables: []ec2types.TransitGatewayPolicyTable{},
			}, nil,
		)

	lister := &EC2TGWPolicyTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TGWPolicyTable_ListWithMultiplePages(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.TransitGatewayPolicyTable, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.TransitGatewayPolicyTable{
			TransitGatewayPolicyTableId: ptr.String(fmt.Sprintf("tgw-ptb-%017d", i)),
			TransitGatewayId:            ptr.String("tgw-11111111111111111"),
			State:                       ec2types.TransitGatewayPolicyTableStateAvailable,
		}
	}

	mockClient.
		On(
			"DescribeTransitGatewayPolicyTables",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayPolicyTablesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayPolicyTablesOutput{
				TransitGatewayPolicyTables: firstPageItems,
				NextToken:                  ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeTransitGatewayPolicyTables",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTransitGatewayPolicyTablesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeTransitGatewayPolicyTablesOutput{
				TransitGatewayPolicyTables: []ec2types.TransitGatewayPolicyTable{
					{
						TransitGatewayPolicyTableId: ptr.String("tgw-ptb-00000000000000100"),
						TransitGatewayId:            ptr.String("tgw-11111111111111111"),
						State:                       ec2types.TransitGatewayPolicyTableStateAvailable,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2TGWPolicyTableLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

// --- Removal ---

func Test_Mock_EC2TGWPolicyTable_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	policyTable := &EC2TGWPolicyTable{
		svc:                         mockClient,
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteTransitGatewayPolicyTable",
			mock.Anything,
			&ec2.DeleteTransitGatewayPolicyTableInput{
				TransitGatewayPolicyTableId: policyTable.TransitGatewayPolicyTableID,
			},
		).
		Return(&ec2.DeleteTransitGatewayPolicyTableOutput{}, nil)

	err := policyTable.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

// --- Properties ---

func Test_Mock_EC2TGWPolicyTable_Properties(t *testing.T) {
	assertions := assert.New(t)

	policyTable := EC2TGWPolicyTable{
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-11111111111111111"),
		TransitGatewayID:            ptr.String("tgw-11111111111111111"),
		State:                       "available",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := policyTable.Properties()

	assertions.Equal("tgw-ptb-11111111111111111", properties.Get("TransitGatewayPolicyTableId"))
	assertions.Equal("tgw-11111111111111111", properties.Get("TransitGatewayId"))
	assertions.Equal("available", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2TGWPolicyTable_PropertiesWithEmptyTags(t *testing.T) {
	assertions := assert.New(t)

	policyTable := EC2TGWPolicyTable{
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-22222222222222222"),
		TransitGatewayID:            ptr.String("tgw-22222222222222222"),
		State:                       "available",
		Tags:                        []ec2types.Tag{},
	}

	properties := policyTable.Properties()

	assertions.Equal("tgw-ptb-22222222222222222", properties.Get("TransitGatewayPolicyTableId"))
}

// --- Display ---

func Test_Mock_EC2TGWPolicyTable_String(t *testing.T) {
	assertions := assert.New(t)

	policyTable := EC2TGWPolicyTable{
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-11111111111111111"),
	}

	assertions.Equal("tgw-ptb-11111111111111111", policyTable.String())
}

// --- Filter ---

func Test_Mock_EC2TGWPolicyTable_FilterExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	policyTable := EC2TGWPolicyTable{
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-11111111111111111"),
		State:                       string(ec2types.TransitGatewayPolicyTableStateDeleted),
	}

	err := policyTable.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2TGWPolicyTable_FilterPassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	policyTable := EC2TGWPolicyTable{
		TransitGatewayPolicyTableID: ptr.String("tgw-ptb-11111111111111111"),
		State:                       string(ec2types.TransitGatewayPolicyTableStateAvailable),
	}

	err := policyTable.Filter()
	assertions.NoError(err)
}

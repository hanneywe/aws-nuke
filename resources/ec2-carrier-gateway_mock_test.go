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

func Test_Mock_EC2CarrierGateway_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeCarrierGateways", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeCarrierGatewaysOutput{
				CarrierGateways: []ec2types.CarrierGateway{
					{
						CarrierGatewayId: ptr.String("cagw-11111111111111111"),
						VpcId:            ptr.String("vpc-22222222222222222"),
						State:            ec2types.CarrierGatewayStateAvailable,
						OwnerId:          ptr.String("123456789012"),
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-cagw")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2CarrierGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	gateway := resources[0].(*EC2CarrierGateway)
	assertions.Equal("cagw-11111111111111111", *gateway.CarrierGatewayID)
	assertions.Equal("vpc-22222222222222222", *gateway.VpcID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2CarrierGateway_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeCarrierGateways", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeCarrierGatewaysOutput{
				CarrierGateways: []ec2types.CarrierGateway{},
			}, nil,
		)

	lister := &EC2CarrierGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2CarrierGateway_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.CarrierGateway, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.CarrierGateway{
			CarrierGatewayId: ptr.String(fmt.Sprintf("cagw-%d", i)),
			VpcId:            ptr.String(fmt.Sprintf("vpc-%d", i)),
			State:            ec2types.CarrierGatewayStateAvailable,
			OwnerId:          ptr.String("123456789012"),
		}
	}

	mockClient.
		On(
			"DescribeCarrierGateways",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeCarrierGatewaysInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeCarrierGatewaysOutput{
				CarrierGateways: firstPageItems,
				NextToken:       ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeCarrierGateways",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeCarrierGatewaysInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeCarrierGatewaysOutput{
				CarrierGateways: []ec2types.CarrierGateway{
					{
						CarrierGatewayId: ptr.String("cagw-100"),
						VpcId:            ptr.String("vpc-100"),
						State:            ec2types.CarrierGatewayStateAvailable,
						OwnerId:          ptr.String("123456789012"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2CarrierGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2CarrierGateway_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	gateway := &EC2CarrierGateway{
		svc:              mockClient,
		CarrierGatewayID: ptr.String("cagw-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteCarrierGateway",
			mock.Anything,
			&ec2.DeleteCarrierGatewayInput{
				CarrierGatewayId: gateway.CarrierGatewayID,
			},
		).
		Return(&ec2.DeleteCarrierGatewayOutput{}, nil)

	err := gateway.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2CarrierGateway_Properties(t *testing.T) {
	assertions := assert.New(t)

	gateway := EC2CarrierGateway{
		CarrierGatewayID: ptr.String("cagw-11111111111111111"),
		VpcID:            ptr.String("vpc-22222222222222222"),
		State:            "available",
		OwnerID:          ptr.String("123456789012"),
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := gateway.Properties()

	assertions.Equal("cagw-11111111111111111", properties.Get("CarrierGatewayId"))
	assertions.Equal("vpc-22222222222222222", properties.Get("VpcId"))
	assertions.Equal("available", properties.Get("State"))
	assertions.Equal("123456789012", properties.Get("OwnerId"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2CarrierGateway_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	gateway := EC2CarrierGateway{
		CarrierGatewayID: ptr.String("cagw-99999999999999999"),
		VpcID:            ptr.String("vpc-88888888888888888"),
		State:            "available",
		OwnerID:          ptr.String("123456789012"),
		Tags:             []ec2types.Tag{},
	}

	properties := gateway.Properties()

	assertions.Equal("cagw-99999999999999999", properties.Get("CarrierGatewayId"))
	assertions.Equal("available", properties.Get("State"))
}

func Test_Mock_EC2CarrierGateway_String(t *testing.T) {
	assertions := assert.New(t)

	gateway := EC2CarrierGateway{
		CarrierGatewayID: ptr.String("cagw-11111111111111111"),
	}

	assertions.Equal("cagw-11111111111111111", gateway.String())
}

func Test_Mock_EC2CarrierGateway_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	gateway := EC2CarrierGateway{
		CarrierGatewayID: ptr.String("cagw-deleted"),
		State:            string(ec2types.CarrierGatewayStateDeleted),
	}

	err := gateway.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2CarrierGateway_Filter_PassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	gateway := EC2CarrierGateway{
		CarrierGatewayID: ptr.String("cagw-active"),
		State:            string(ec2types.CarrierGatewayStateAvailable),
	}

	err := gateway.Filter()
	assertions.NoError(err)
}

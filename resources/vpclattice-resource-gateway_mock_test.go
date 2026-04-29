package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	lattice "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
)

func Test_Mock_VPCLatticeResourceGateway_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListResourceGateways", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListResourceGatewaysOutput{
				Items: []lattice.ResourceGatewaySummary{
					{
						Id:   ptr.String("rgw-1"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourcegateway/rgw-1"),
						Name: ptr.String("gateway-1"),
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{"env": "dev"},
			}, nil,
		)

	lister := &VPCLatticeResourceGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	resourceGateway := resources[0].(*VPCLatticeResourceGateway)
	assertions.Equal("gateway-1", *resourceGateway.Name)
	assertions.Equal("dev", resourceGateway.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceGateway_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListResourceGateways", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListResourceGatewaysOutput{
				Items: []lattice.ResourceGatewaySummary{},
			}, nil,
		)

	lister := &VPCLatticeResourceGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceGateway_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	firstPageItems := make([]lattice.ResourceGatewaySummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ResourceGatewaySummary{
			Id:   ptr.String(fmt.Sprintf("rgw-%d", i)),
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:resourcegateway/rgw-%d", i)),
			Name: ptr.String(fmt.Sprintf("gateway-%d", i)),
		}
	}

	mockClient.
		On(
			"ListResourceGateways",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListResourceGatewaysInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListResourceGatewaysOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListResourceGateways",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListResourceGatewaysInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListResourceGatewaysOutput{
				Items: []lattice.ResourceGatewaySummary{
					{
						Id:   ptr.String("rgw-100"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourcegateway/rgw-100"),
						Name: ptr.String("gateway-100"),
					},
				},
			}, nil,
		).
		Once()

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{},
			}, nil,
		)

	lister := &VPCLatticeResourceGatewayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceGateway_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	resourceGateway := &VPCLatticeResourceGateway{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourcegateway/rgw-1"),
	}

	mockClient.
		On(
			"DeleteResourceGateway",
			mock.Anything,
			&vpclattice.DeleteResourceGatewayInput{
				ResourceGatewayIdentifier: resourceGateway.ARN,
			},
		).
		Return(&vpclattice.DeleteResourceGatewayOutput{}, nil)

	err := resourceGateway.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeResourceGateway_Properties(t *testing.T) {
	assertions := assert.New(t)

	resourceGateway := VPCLatticeResourceGateway{
		ID:   ptr.String("rgw-12345"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:resourcegateway/rgw-12345"),
		Name: ptr.String("my-resource-gateway"),
		Tags: map[string]string{"Environment": "staging"},
	}

	properties := resourceGateway.Properties()

	assertions.Equal("rgw-12345", properties.Get("ID"))
	assertions.Equal("my-resource-gateway", properties.Get("Name"))
	assertions.Equal("staging", properties.Get("tag:Environment"))
}

func Test_Mock_VPCLatticeResourceGateway_String(t *testing.T) {
	assertions := assert.New(t)

	resourceGateway := VPCLatticeResourceGateway{
		Name: ptr.String("my-resource-gateway"),
	}

	assertions.Equal("my-resource-gateway", resourceGateway.String())
}

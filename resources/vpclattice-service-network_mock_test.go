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

func Test_Mock_VPCLatticeServiceNetwork_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{
						Id:   ptr.String("sn-1"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
						Name: ptr.String("network-1"),
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{"env": "prod"},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	serviceNetwork := resources[0].(*VPCLatticeServiceNetwork)
	assertions.Equal("network-1", *serviceNetwork.Name)
	assertions.Equal("prod", serviceNetwork.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetwork_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetwork_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	firstPageItems := make([]lattice.ServiceNetworkSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ServiceNetworkSummary{
			Id:   ptr.String(fmt.Sprintf("sn-%d", i)),
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-%d", i)),
			Name: ptr.String(fmt.Sprintf("network-%d", i)),
		}
	}

	mockClient.
		On(
			"ListServiceNetworks",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworksInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListServiceNetworks",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworksInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{
						Id:   ptr.String("sn-100"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-100"),
						Name: ptr.String("network-100"),
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

	lister := &VPCLatticeServiceNetworkLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetwork_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	serviceNetwork := &VPCLatticeServiceNetwork{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
	}

	mockClient.
		On(
			"DeleteServiceNetwork",
			mock.Anything,
			&vpclattice.DeleteServiceNetworkInput{
				ServiceNetworkIdentifier: serviceNetwork.ARN,
			},
		).
		Return(&vpclattice.DeleteServiceNetworkOutput{}, nil)

	err := serviceNetwork.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetwork_Properties(t *testing.T) {
	assertions := assert.New(t)

	serviceNetwork := VPCLatticeServiceNetwork{
		ID:   ptr.String("sn-12345"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-12345"),
		Name: ptr.String("my-network"),
		Tags: map[string]string{"Environment": "dev"},
	}

	properties := serviceNetwork.Properties()

	assertions.Equal("sn-12345", properties.Get("ID"))
	assertions.Equal("my-network", properties.Get("Name"))
	assertions.Equal("dev", properties.Get("tag:Environment"))
}

func Test_Mock_VPCLatticeServiceNetwork_String(t *testing.T) {
	assertions := assert.New(t)

	serviceNetwork := VPCLatticeServiceNetwork{
		Name: ptr.String("my-network"),
	}

	assertions.Equal("my-network", serviceNetwork.String())
}

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

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{Id: ptr.String("sn-1")},
				},
			}, nil,
		)

	mockClient.
		On("ListServiceNetworkServiceAssociations", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworkServiceAssociationsOutput{
				Items: []lattice.ServiceNetworkServiceAssociationSummary{
					{
						Id:                 ptr.String("snsa-1"),
						Arn:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snsa/snsa-1"),
						ServiceNetworkName: ptr.String("network-1"),
						ServiceName:        ptr.String("service-1"),
						Status:             lattice.ServiceNetworkServiceAssociationStatusActive,
					},
				},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkServiceAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 1)

	association := resources[0].(*VPCLatticeServiceNetworkServiceAssociation)
	assertions.Equal("network-1", *association.ServiceNetworkName)
	assertions.Equal("service-1", *association.ServiceName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkServiceAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{Id: ptr.String("sn-1")},
				},
			}, nil,
		)

	firstPageItems := make([]lattice.ServiceNetworkServiceAssociationSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ServiceNetworkServiceAssociationSummary{
			Id:                 ptr.String(fmt.Sprintf("snsa-%d", i)),
			Arn:                ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:snsa/snsa-%d", i)),
			ServiceNetworkName: ptr.String("network"),
			ServiceName:        ptr.String(fmt.Sprintf("service-%d", i)),
			Status:             lattice.ServiceNetworkServiceAssociationStatusActive,
		}
	}

	mockClient.
		On(
			"ListServiceNetworkServiceAssociations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworkServiceAssociationsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworkServiceAssociationsOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListServiceNetworkServiceAssociations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworkServiceAssociationsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworkServiceAssociationsOutput{
				Items: []lattice.ServiceNetworkServiceAssociationSummary{
					{
						Id:                 ptr.String("snsa-100"),
						Arn:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snsa/snsa-100"),
						ServiceNetworkName: ptr.String("network"),
						ServiceName:        ptr.String("service-100"),
						Status:             lattice.ServiceNetworkServiceAssociationStatusActive,
					},
				},
			}, nil,
		).
		Once()

	lister := &VPCLatticeServiceNetworkServiceAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	association := &VPCLatticeServiceNetworkServiceAssociation{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snsa/snsa-1"),
	}

	mockClient.
		On(
			"DeleteServiceNetworkServiceAssociation",
			mock.Anything,
			&vpclattice.DeleteServiceNetworkServiceAssociationInput{
				ServiceNetworkServiceAssociationIdentifier: association.ARN,
			},
		).
		Return(&vpclattice.DeleteServiceNetworkServiceAssociationOutput{}, nil)

	err := association.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_Properties(t *testing.T) {
	assertions := assert.New(t)

	association := VPCLatticeServiceNetworkServiceAssociation{
		ID:                 ptr.String("snsa-12345"),
		ARN:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snsa/snsa-12345"),
		ServiceNetworkName: ptr.String("my-network"),
		ServiceName:        ptr.String("my-service"),
		Status:             ptr.String("ACTIVE"),
	}

	properties := association.Properties()

	assertions.Equal("snsa-12345", properties.Get("ID"))
	assertions.Equal("my-network", properties.Get("ServiceNetworkName"))
	assertions.Equal("my-service", properties.Get("ServiceName"))
	assertions.Equal("ACTIVE", properties.Get("Status"))
}

func Test_Mock_VPCLatticeServiceNetworkServiceAssociation_String(t *testing.T) {
	assertions := assert.New(t)

	association := VPCLatticeServiceNetworkServiceAssociation{
		ServiceNetworkName: ptr.String("my-network"),
		ServiceName:        ptr.String("my-service"),
	}

	assertions.Equal("my-network -> my-service", association.String())
}

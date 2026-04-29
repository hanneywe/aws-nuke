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

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_List_One(t *testing.T) {
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
		On("ListServiceNetworkVpcAssociations", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworkVpcAssociationsOutput{
				Items: []lattice.ServiceNetworkVpcAssociationSummary{
					{
						Id:                 ptr.String("snva-1"),
						Arn:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snva/snva-1"),
						ServiceNetworkName: ptr.String("network-1"),
						VpcId:              ptr.String("vpc-abc123"),
						Status:             lattice.ServiceNetworkVpcAssociationStatusActive,
					},
				},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkVPCAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 1)

	association := resources[0].(*VPCLatticeServiceNetworkVPCAssociation)
	assertions.Equal("network-1", *association.ServiceNetworkName)
	assertions.Equal("vpc-abc123", *association.VPCID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{},
			}, nil,
		)

	lister := &VPCLatticeServiceNetworkVPCAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_List_MultiPage(t *testing.T) {
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

	firstPageItems := make([]lattice.ServiceNetworkVpcAssociationSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ServiceNetworkVpcAssociationSummary{
			Id:                 ptr.String(fmt.Sprintf("snva-%d", i)),
			Arn:                ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:snva/snva-%d", i)),
			ServiceNetworkName: ptr.String("network"),
			VpcId:              ptr.String(fmt.Sprintf("vpc-%d", i)),
			Status:             lattice.ServiceNetworkVpcAssociationStatusActive,
		}
	}

	mockClient.
		On(
			"ListServiceNetworkVpcAssociations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworkVpcAssociationsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworkVpcAssociationsOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListServiceNetworkVpcAssociations",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServiceNetworkVpcAssociationsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListServiceNetworkVpcAssociationsOutput{
				Items: []lattice.ServiceNetworkVpcAssociationSummary{
					{
						Id:                 ptr.String("snva-100"),
						Arn:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snva/snva-100"),
						ServiceNetworkName: ptr.String("network"),
						VpcId:              ptr.String("vpc-100"),
						Status:             lattice.ServiceNetworkVpcAssociationStatusActive,
					},
				},
			}, nil,
		).
		Once()

	lister := &VPCLatticeServiceNetworkVPCAssociationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	association := &VPCLatticeServiceNetworkVPCAssociation{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snva/snva-1"),
	}

	mockClient.
		On(
			"DeleteServiceNetworkVpcAssociation",
			mock.Anything,
			&vpclattice.DeleteServiceNetworkVpcAssociationInput{
				ServiceNetworkVpcAssociationIdentifier: association.ARN,
			},
		).
		Return(&vpclattice.DeleteServiceNetworkVpcAssociationOutput{}, nil)

	err := association.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_Properties(t *testing.T) {
	assertions := assert.New(t)

	association := VPCLatticeServiceNetworkVPCAssociation{
		ID:                 ptr.String("snva-12345"),
		ARN:                ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:snva/snva-12345"),
		ServiceNetworkName: ptr.String("my-network"),
		VPCID:              ptr.String("vpc-abc123"),
		Status:             ptr.String("ACTIVE"),
	}

	properties := association.Properties()

	assertions.Equal("snva-12345", properties.Get("ID"))
	assertions.Equal("my-network", properties.Get("ServiceNetworkName"))
	assertions.Equal("vpc-abc123", properties.Get("VPCID"))
	assertions.Equal("ACTIVE", properties.Get("Status"))
}

func Test_Mock_VPCLatticeServiceNetworkVPCAssociation_String(t *testing.T) {
	assertions := assert.New(t)

	association := VPCLatticeServiceNetworkVPCAssociation{
		ServiceNetworkName: ptr.String("my-network"),
		VPCID:              ptr.String("vpc-abc123"),
	}

	assertions.Equal("my-network -> vpc-abc123", association.String())
}

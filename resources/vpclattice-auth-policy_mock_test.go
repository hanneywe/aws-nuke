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

func Test_Mock_VPCLatticeAuthPolicy_List_One_ServiceNetwork(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
						Id:   ptr.String("sn-1"),
						Name: ptr.String("network-1"),
					},
				},
			}, nil,
		)

	mockClient.
		On(
			"GetAuthPolicy",
			mock.Anything,
			&vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
			},
		).
		Return(
			&vpclattice.GetAuthPolicyOutput{
				Policy: ptr.String(`{"Version":"2012-10-17"}`),
			}, nil,
		)

	// No services with policies
	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeAuthPolicyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 1)

	authPolicy := resources[0].(*VPCLatticeAuthPolicy)
	assertions.Equal("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1", *authPolicy.ResourceARN)
	assertions.Equal("SERVICE_NETWORK", *authPolicy.ResourceType)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_List_One_Service(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// No service networks with policies
	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{},
			}, nil,
		)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-1"),
						Id:   ptr.String("svc-1"),
						Name: ptr.String("service-one"),
					},
				},
			}, nil,
		)

	mockClient.
		On(
			"GetAuthPolicy",
			mock.Anything,
			&vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-1"),
			},
		).
		Return(
			&vpclattice.GetAuthPolicyOutput{
				Policy: ptr.String(`{"Version":"2012-10-17"}`),
			}, nil,
		)

	lister := &VPCLatticeAuthPolicyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 1)

	authPolicy := resources[0].(*VPCLatticeAuthPolicy)
	assertions.Equal("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-1", *authPolicy.ResourceARN)
	assertions.Equal("SERVICE", *authPolicy.ResourceType)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_List_Empty_NoPolicies(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// Service network exists but GetAuthPolicy returns error (no policy attached)
	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
						Id:   ptr.String("sn-1"),
						Name: ptr.String("network-1"),
					},
				},
			}, nil,
		)

	mockClient.
		On(
			"GetAuthPolicy",
			mock.Anything,
			&vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
			},
		).
		Return(
			&vpclattice.GetAuthPolicyOutput{},
			fmt.Errorf("ResourceNotFoundException"),
		)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeAuthPolicyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_List_Empty_EmptyPolicy(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// Policy exists but is an empty string — should be skipped
	mockClient.
		On("ListServiceNetworks", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServiceNetworksOutput{
				Items: []lattice.ServiceNetworkSummary{
					{
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
						Id:   ptr.String("sn-1"),
						Name: ptr.String("network-1"),
					},
				},
			}, nil,
		)

	mockClient.
		On(
			"GetAuthPolicy",
			mock.Anything,
			&vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
			},
		).
		Return(
			&vpclattice.GetAuthPolicyOutput{
				Policy: ptr.String(""),
			}, nil,
		)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeAuthPolicyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// 101 service networks across 2 pages, each with an auth policy
	firstPageItems := make([]lattice.ServiceNetworkSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ServiceNetworkSummary{
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-%d", i)),
			Id:   ptr.String(fmt.Sprintf("sn-%d", i)),
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
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-100"),
						Id:   ptr.String("sn-100"),
						Name: ptr.String("network-100"),
					},
				},
			}, nil,
		).
		Once()

	mockClient.
		On("GetAuthPolicy", mock.Anything, mock.Anything).
		Return(
			&vpclattice.GetAuthPolicyOutput{
				Policy: ptr.String(`{"Version":"2012-10-17"}`),
			}, nil,
		)

	// No services
	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeAuthPolicyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	authPolicy := &VPCLatticeAuthPolicy{
		svc:         mockClient,
		ResourceARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-1"),
	}

	mockClient.
		On(
			"DeleteAuthPolicy",
			mock.Anything,
			&vpclattice.DeleteAuthPolicyInput{
				ResourceIdentifier: authPolicy.ResourceARN,
			},
		).
		Return(&vpclattice.DeleteAuthPolicyOutput{}, nil)

	err := authPolicy.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeAuthPolicy_Properties(t *testing.T) {
	assertions := assert.New(t)

	authPolicy := VPCLatticeAuthPolicy{
		ResourceARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-12345"),
		ResourceType: ptr.String("SERVICE_NETWORK"),
	}

	properties := authPolicy.Properties()

	assertions.Equal(
		"arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-12345",
		properties.Get("ResourceARN"),
	)
	assertions.Equal("SERVICE_NETWORK", properties.Get("ResourceType"))
}

func Test_Mock_VPCLatticeAuthPolicy_String(t *testing.T) {
	assertions := assert.New(t)

	authPolicy := VPCLatticeAuthPolicy{
		ResourceARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-12345"),
	}

	assertions.Equal(
		"arn:aws:vpc-lattice:us-east-1:123456789012:servicenetwork/sn-12345",
		authPolicy.String(),
	)
}

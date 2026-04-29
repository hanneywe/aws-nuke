package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
)

// VPCLatticeClient is an interface for the VPC Lattice SDK client methods used by all VPC Lattice resources.
// It enables mock testing of List and Remove operations.
type VPCLatticeClient interface {
	// Listing
	ListServiceNetworks(ctx context.Context, params *vpclattice.ListServiceNetworksInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListServiceNetworksOutput, error)
	ListServices(ctx context.Context, params *vpclattice.ListServicesInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListServicesOutput, error)
	ListTargetGroups(ctx context.Context, params *vpclattice.ListTargetGroupsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListTargetGroupsOutput, error)
	ListListeners(ctx context.Context, params *vpclattice.ListListenersInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListListenersOutput, error)
	ListRules(ctx context.Context, params *vpclattice.ListRulesInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListRulesOutput, error)
	ListServiceNetworkServiceAssociations(ctx context.Context, params *vpclattice.ListServiceNetworkServiceAssociationsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListServiceNetworkServiceAssociationsOutput, error)
	ListServiceNetworkVpcAssociations(ctx context.Context, params *vpclattice.ListServiceNetworkVpcAssociationsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListServiceNetworkVpcAssociationsOutput, error)
	ListResourceConfigurations(ctx context.Context, params *vpclattice.ListResourceConfigurationsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListResourceConfigurationsOutput, error)
	ListResourceGateways(ctx context.Context, params *vpclattice.ListResourceGatewaysInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListResourceGatewaysOutput, error)
	ListTagsForResource(ctx context.Context, params *vpclattice.ListTagsForResourceInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListTagsForResourceOutput, error)

	// Auth policy
	GetAuthPolicy(ctx context.Context, params *vpclattice.GetAuthPolicyInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.GetAuthPolicyOutput, error)

	// Deletion
	DeleteServiceNetwork(ctx context.Context, params *vpclattice.DeleteServiceNetworkInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceNetworkOutput, error)
	DeleteService(ctx context.Context, params *vpclattice.DeleteServiceInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceOutput, error)
	DeleteTargetGroup(ctx context.Context, params *vpclattice.DeleteTargetGroupInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteTargetGroupOutput, error)
	DeleteListener(ctx context.Context, params *vpclattice.DeleteListenerInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteListenerOutput, error)
	DeleteRule(ctx context.Context, params *vpclattice.DeleteRuleInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteRuleOutput, error)
	DeleteAuthPolicy(ctx context.Context, params *vpclattice.DeleteAuthPolicyInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteAuthPolicyOutput, error)
	DeleteServiceNetworkServiceAssociation(ctx context.Context, params *vpclattice.DeleteServiceNetworkServiceAssociationInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceNetworkServiceAssociationOutput, error)
	DeleteServiceNetworkVpcAssociation(ctx context.Context, params *vpclattice.DeleteServiceNetworkVpcAssociationInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteServiceNetworkVpcAssociationOutput, error)
	DeleteResourceConfiguration(ctx context.Context, params *vpclattice.DeleteResourceConfigurationInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteResourceConfigurationOutput, error)
	DeleteResourceGateway(ctx context.Context, params *vpclattice.DeleteResourceGatewayInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeleteResourceGatewayOutput, error)
	ListTargets(ctx context.Context, params *vpclattice.ListTargetsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.ListTargetsOutput, error)
	DeregisterTargets(ctx context.Context, params *vpclattice.DeregisterTargetsInput,
		optFns ...func(*vpclattice.Options)) (*vpclattice.DeregisterTargetsOutput, error)
}

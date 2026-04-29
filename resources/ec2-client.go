package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2Client is a shared interface for EC2 SDK v2 client methods used by multiple resources.
// It enables mock testing of List and Remove operations.
type EC2Client interface {
	// Transit Gateway Route Tables
	DescribeTransitGatewayRouteTables(ctx context.Context, params *ec2.DescribeTransitGatewayRouteTablesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTransitGatewayRouteTablesOutput, error)
	DeleteTransitGatewayRouteTable(ctx context.Context, params *ec2.DeleteTransitGatewayRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTransitGatewayRouteTableOutput, error)
	ModifyTransitGateway(ctx context.Context, params *ec2.ModifyTransitGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.ModifyTransitGatewayOutput, error)

	// Transit Gateway Multicast Domains
	DescribeTransitGatewayMulticastDomains(ctx context.Context, params *ec2.DescribeTransitGatewayMulticastDomainsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTransitGatewayMulticastDomainsOutput, error)
	DeleteTransitGatewayMulticastDomain(ctx context.Context, params *ec2.DeleteTransitGatewayMulticastDomainInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTransitGatewayMulticastDomainOutput, error)

	// Transit Gateway Policy Tables
	DescribeTransitGatewayPolicyTables(ctx context.Context, params *ec2.DescribeTransitGatewayPolicyTablesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTransitGatewayPolicyTablesOutput, error)
	DeleteTransitGatewayPolicyTable(ctx context.Context, params *ec2.DeleteTransitGatewayPolicyTableInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTransitGatewayPolicyTableOutput, error)

	// Transit Gateway Prefix List References
	GetTransitGatewayPrefixListReferences(ctx context.Context, params *ec2.GetTransitGatewayPrefixListReferencesInput,
		optFns ...func(*ec2.Options)) (*ec2.GetTransitGatewayPrefixListReferencesOutput, error)
	DeleteTransitGatewayPrefixListReference(ctx context.Context, params *ec2.DeleteTransitGatewayPrefixListReferenceInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTransitGatewayPrefixListReferenceOutput, error)

	// Flow Logs
	DescribeFlowLogs(ctx context.Context, params *ec2.DescribeFlowLogsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeFlowLogsOutput, error)
	DeleteFlowLogs(ctx context.Context, params *ec2.DeleteFlowLogsInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteFlowLogsOutput, error)

	// Managed Prefix Lists
	DescribeManagedPrefixLists(ctx context.Context, params *ec2.DescribeManagedPrefixListsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeManagedPrefixListsOutput, error)
	DeleteManagedPrefixList(ctx context.Context, params *ec2.DeleteManagedPrefixListInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteManagedPrefixListOutput, error)

	// Network Insights Paths
	DescribeNetworkInsightsPaths(ctx context.Context, params *ec2.DescribeNetworkInsightsPathsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInsightsPathsOutput, error)
	DeleteNetworkInsightsPath(ctx context.Context, params *ec2.DeleteNetworkInsightsPathInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInsightsPathOutput, error)

	// Network Insights Analyses
	DescribeNetworkInsightsAnalyses(ctx context.Context, params *ec2.DescribeNetworkInsightsAnalysesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInsightsAnalysesOutput, error)
	DeleteNetworkInsightsAnalysis(ctx context.Context, params *ec2.DeleteNetworkInsightsAnalysisInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInsightsAnalysisOutput, error)

	// Traffic Mirror Sessions
	DescribeTrafficMirrorSessions(ctx context.Context, params *ec2.DescribeTrafficMirrorSessionsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTrafficMirrorSessionsOutput, error)
	DeleteTrafficMirrorSession(ctx context.Context, params *ec2.DeleteTrafficMirrorSessionInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTrafficMirrorSessionOutput, error)

	// Traffic Mirror Targets
	DescribeTrafficMirrorTargets(ctx context.Context, params *ec2.DescribeTrafficMirrorTargetsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTrafficMirrorTargetsOutput, error)
	DeleteTrafficMirrorTarget(ctx context.Context, params *ec2.DeleteTrafficMirrorTargetInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTrafficMirrorTargetOutput, error)

	// Traffic Mirror Filters
	DescribeTrafficMirrorFilters(ctx context.Context, params *ec2.DescribeTrafficMirrorFiltersInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTrafficMirrorFiltersOutput, error)
	DeleteTrafficMirrorFilter(ctx context.Context, params *ec2.DeleteTrafficMirrorFilterInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTrafficMirrorFilterOutput, error)

	// Carrier Gateways
	DescribeCarrierGateways(ctx context.Context, params *ec2.DescribeCarrierGatewaysInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeCarrierGatewaysOutput, error)
	DeleteCarrierGateway(ctx context.Context, params *ec2.DeleteCarrierGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteCarrierGatewayOutput, error)

	// IPAM
	DescribeIpams(ctx context.Context, params *ec2.DescribeIpamsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeIpamsOutput, error)
	DeleteIpam(ctx context.Context, params *ec2.DeleteIpamInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteIpamOutput, error)

	// IPAM Pools
	DescribeIpamPools(ctx context.Context, params *ec2.DescribeIpamPoolsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeIpamPoolsOutput, error)
	DeleteIpamPool(ctx context.Context, params *ec2.DeleteIpamPoolInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteIpamPoolOutput, error)

	// IPAM Scopes
	DescribeIpamScopes(ctx context.Context, params *ec2.DescribeIpamScopesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeIpamScopesOutput, error)
	DeleteIpamScope(ctx context.Context, params *ec2.DeleteIpamScopeInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteIpamScopeOutput, error)

	// VPC CIDR Blocks
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DisassociateVpcCidrBlock(ctx context.Context, params *ec2.DisassociateVpcCidrBlockInput,
		optFns ...func(*ec2.Options)) (*ec2.DisassociateVpcCidrBlockOutput, error)

	// Route Servers
	DescribeRouteServers(ctx context.Context, params *ec2.DescribeRouteServersInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeRouteServersOutput, error)
	DeleteRouteServer(ctx context.Context, params *ec2.DeleteRouteServerInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteRouteServerOutput, error)

	// Fleets
	DescribeFleets(ctx context.Context, params *ec2.DescribeFleetsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeFleetsOutput, error)
	DeleteFleets(ctx context.Context, params *ec2.DeleteFleetsInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteFleetsOutput, error)

	// Network Insights Access Scopes
	DescribeNetworkInsightsAccessScopes(ctx context.Context, params *ec2.DescribeNetworkInsightsAccessScopesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInsightsAccessScopesOutput, error)
	DeleteNetworkInsightsAccessScope(ctx context.Context, params *ec2.DeleteNetworkInsightsAccessScopeInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteNetworkInsightsAccessScopeOutput, error)

	// VPC Encryption Controls
	DescribeVpcEncryptionControls(ctx context.Context, params *ec2.DescribeVpcEncryptionControlsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcEncryptionControlsOutput, error)
	DeleteVpcEncryptionControl(ctx context.Context, params *ec2.DeleteVpcEncryptionControlInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteVpcEncryptionControlOutput, error)

	// Public IPv4 Pools
	DescribePublicIpv4Pools(ctx context.Context, params *ec2.DescribePublicIpv4PoolsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribePublicIpv4PoolsOutput, error)
	DeletePublicIpv4Pool(ctx context.Context, params *ec2.DeletePublicIpv4PoolInput,
		optFns ...func(*ec2.Options)) (*ec2.DeletePublicIpv4PoolOutput, error)

	// VPN Concentrators
	DescribeVpnConcentrators(ctx context.Context, params *ec2.DescribeVpnConcentratorsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpnConcentratorsOutput, error)
	DeleteVpnConcentrator(ctx context.Context, params *ec2.DeleteVpnConcentratorInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteVpnConcentratorOutput, error)

	// IPAM Prefix List Resolvers
	DescribeIpamPrefixListResolvers(ctx context.Context, params *ec2.DescribeIpamPrefixListResolversInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeIpamPrefixListResolversOutput, error)
	DeleteIpamPrefixListResolver(ctx context.Context, params *ec2.DeleteIpamPrefixListResolverInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteIpamPrefixListResolverOutput, error)

	// Traffic Mirror Filter Rules
	DescribeTrafficMirrorFilterRules(ctx context.Context, params *ec2.DescribeTrafficMirrorFilterRulesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTrafficMirrorFilterRulesOutput, error)
	DeleteTrafficMirrorFilterRule(ctx context.Context, params *ec2.DeleteTrafficMirrorFilterRuleInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteTrafficMirrorFilterRuleOutput, error)
}

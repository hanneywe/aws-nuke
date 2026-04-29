package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testEC2ClientListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockEC2Client struct {
	mock.Mock
}

// Transit Gateway Route Tables

func (m *mockEC2Client) DescribeTransitGatewayRouteTables(
	ctx context.Context, params *ec2.DescribeTransitGatewayRouteTablesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTransitGatewayRouteTablesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTransitGatewayRouteTablesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTransitGatewayRouteTable(
	ctx context.Context, params *ec2.DeleteTransitGatewayRouteTableInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTransitGatewayRouteTableOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTransitGatewayRouteTableOutput), args.Error(1)
}

func (m *mockEC2Client) ModifyTransitGateway(
	ctx context.Context, params *ec2.ModifyTransitGatewayInput, _ ...func(*ec2.Options),
) (*ec2.ModifyTransitGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.ModifyTransitGatewayOutput), args.Error(1)
}

// Transit Gateway Multicast Domains

func (m *mockEC2Client) DescribeTransitGatewayMulticastDomains(
	ctx context.Context, params *ec2.DescribeTransitGatewayMulticastDomainsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTransitGatewayMulticastDomainsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTransitGatewayMulticastDomainsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTransitGatewayMulticastDomain(
	ctx context.Context, params *ec2.DeleteTransitGatewayMulticastDomainInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTransitGatewayMulticastDomainOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTransitGatewayMulticastDomainOutput), args.Error(1)
}

// Transit Gateway Policy Tables

func (m *mockEC2Client) DescribeTransitGatewayPolicyTables(
	ctx context.Context, params *ec2.DescribeTransitGatewayPolicyTablesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTransitGatewayPolicyTablesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTransitGatewayPolicyTablesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTransitGatewayPolicyTable(
	ctx context.Context, params *ec2.DeleteTransitGatewayPolicyTableInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTransitGatewayPolicyTableOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTransitGatewayPolicyTableOutput), args.Error(1)
}

// Transit Gateway Prefix List References

func (m *mockEC2Client) GetTransitGatewayPrefixListReferences(
	ctx context.Context, params *ec2.GetTransitGatewayPrefixListReferencesInput, _ ...func(*ec2.Options),
) (*ec2.GetTransitGatewayPrefixListReferencesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.GetTransitGatewayPrefixListReferencesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTransitGatewayPrefixListReference(
	ctx context.Context, params *ec2.DeleteTransitGatewayPrefixListReferenceInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTransitGatewayPrefixListReferenceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTransitGatewayPrefixListReferenceOutput), args.Error(1)
}

// Flow Logs

func (m *mockEC2Client) DescribeFlowLogs(
	ctx context.Context, params *ec2.DescribeFlowLogsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeFlowLogsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeFlowLogsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteFlowLogs(
	ctx context.Context, params *ec2.DeleteFlowLogsInput, _ ...func(*ec2.Options),
) (*ec2.DeleteFlowLogsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteFlowLogsOutput), args.Error(1)
}

// Managed Prefix Lists

func (m *mockEC2Client) DescribeManagedPrefixLists(
	ctx context.Context, params *ec2.DescribeManagedPrefixListsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeManagedPrefixListsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeManagedPrefixListsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteManagedPrefixList(
	ctx context.Context, params *ec2.DeleteManagedPrefixListInput, _ ...func(*ec2.Options),
) (*ec2.DeleteManagedPrefixListOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteManagedPrefixListOutput), args.Error(1)
}

// Network Insights Paths

func (m *mockEC2Client) DescribeNetworkInsightsPaths(
	ctx context.Context, params *ec2.DescribeNetworkInsightsPathsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeNetworkInsightsPathsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeNetworkInsightsPathsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteNetworkInsightsPath(
	ctx context.Context, params *ec2.DeleteNetworkInsightsPathInput, _ ...func(*ec2.Options),
) (*ec2.DeleteNetworkInsightsPathOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteNetworkInsightsPathOutput), args.Error(1)
}

// Network Insights Analyses

func (m *mockEC2Client) DescribeNetworkInsightsAnalyses(
	ctx context.Context, params *ec2.DescribeNetworkInsightsAnalysesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeNetworkInsightsAnalysesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeNetworkInsightsAnalysesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteNetworkInsightsAnalysis(
	ctx context.Context, params *ec2.DeleteNetworkInsightsAnalysisInput, _ ...func(*ec2.Options),
) (*ec2.DeleteNetworkInsightsAnalysisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteNetworkInsightsAnalysisOutput), args.Error(1)
}

// Traffic Mirror Sessions

func (m *mockEC2Client) DescribeTrafficMirrorSessions(
	ctx context.Context, params *ec2.DescribeTrafficMirrorSessionsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTrafficMirrorSessionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTrafficMirrorSessionsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTrafficMirrorSession(
	ctx context.Context, params *ec2.DeleteTrafficMirrorSessionInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTrafficMirrorSessionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTrafficMirrorSessionOutput), args.Error(1)
}

// Traffic Mirror Targets

func (m *mockEC2Client) DescribeTrafficMirrorTargets(
	ctx context.Context, params *ec2.DescribeTrafficMirrorTargetsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTrafficMirrorTargetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTrafficMirrorTargetsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTrafficMirrorTarget(
	ctx context.Context, params *ec2.DeleteTrafficMirrorTargetInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTrafficMirrorTargetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTrafficMirrorTargetOutput), args.Error(1)
}

// Traffic Mirror Filters

func (m *mockEC2Client) DescribeTrafficMirrorFilters(
	ctx context.Context, params *ec2.DescribeTrafficMirrorFiltersInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTrafficMirrorFiltersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTrafficMirrorFiltersOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTrafficMirrorFilter(
	ctx context.Context, params *ec2.DeleteTrafficMirrorFilterInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTrafficMirrorFilterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTrafficMirrorFilterOutput), args.Error(1)
}

// Carrier Gateways

func (m *mockEC2Client) DescribeCarrierGateways(
	ctx context.Context, params *ec2.DescribeCarrierGatewaysInput, _ ...func(*ec2.Options),
) (*ec2.DescribeCarrierGatewaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeCarrierGatewaysOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteCarrierGateway(
	ctx context.Context, params *ec2.DeleteCarrierGatewayInput, _ ...func(*ec2.Options),
) (*ec2.DeleteCarrierGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteCarrierGatewayOutput), args.Error(1)
}

// IPAM

func (m *mockEC2Client) DescribeIpams(
	ctx context.Context, params *ec2.DescribeIpamsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeIpamsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeIpamsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteIpam(
	ctx context.Context, params *ec2.DeleteIpamInput, _ ...func(*ec2.Options),
) (*ec2.DeleteIpamOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteIpamOutput), args.Error(1)
}

// IPAM Pools

func (m *mockEC2Client) DescribeIpamPools(
	ctx context.Context, params *ec2.DescribeIpamPoolsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeIpamPoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeIpamPoolsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteIpamPool(
	ctx context.Context, params *ec2.DeleteIpamPoolInput, _ ...func(*ec2.Options),
) (*ec2.DeleteIpamPoolOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteIpamPoolOutput), args.Error(1)
}

// IPAM Scopes

func (m *mockEC2Client) DescribeIpamScopes(
	ctx context.Context, params *ec2.DescribeIpamScopesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeIpamScopesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeIpamScopesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteIpamScope(
	ctx context.Context, params *ec2.DeleteIpamScopeInput, _ ...func(*ec2.Options),
) (*ec2.DeleteIpamScopeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteIpamScopeOutput), args.Error(1)
}

// VPC CIDR Blocks

func (m *mockEC2Client) DescribeVpcs(
	ctx context.Context, params *ec2.DescribeVpcsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeVpcsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeVpcsOutput), args.Error(1)
}

func (m *mockEC2Client) DisassociateVpcCidrBlock(
	ctx context.Context, params *ec2.DisassociateVpcCidrBlockInput, _ ...func(*ec2.Options),
) (*ec2.DisassociateVpcCidrBlockOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DisassociateVpcCidrBlockOutput), args.Error(1)
}

func (m *mockEC2Client) DescribeRouteServers(
	ctx context.Context, params *ec2.DescribeRouteServersInput, _ ...func(*ec2.Options),
) (*ec2.DescribeRouteServersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeRouteServersOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteRouteServer(
	ctx context.Context, params *ec2.DeleteRouteServerInput, _ ...func(*ec2.Options),
) (*ec2.DeleteRouteServerOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteRouteServerOutput), args.Error(1)
}

// Fleets

func (m *mockEC2Client) DescribeFleets(
	ctx context.Context, params *ec2.DescribeFleetsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeFleetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeFleetsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteFleets(
	ctx context.Context, params *ec2.DeleteFleetsInput, _ ...func(*ec2.Options),
) (*ec2.DeleteFleetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteFleetsOutput), args.Error(1)
}

// Network Insights Access Scopes

func (m *mockEC2Client) DescribeNetworkInsightsAccessScopes(
	ctx context.Context, params *ec2.DescribeNetworkInsightsAccessScopesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeNetworkInsightsAccessScopesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeNetworkInsightsAccessScopesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteNetworkInsightsAccessScope(
	ctx context.Context, params *ec2.DeleteNetworkInsightsAccessScopeInput, _ ...func(*ec2.Options),
) (*ec2.DeleteNetworkInsightsAccessScopeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteNetworkInsightsAccessScopeOutput), args.Error(1)
}

// VPC Encryption Controls

func (m *mockEC2Client) DescribeVpcEncryptionControls(
	ctx context.Context, params *ec2.DescribeVpcEncryptionControlsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeVpcEncryptionControlsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeVpcEncryptionControlsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteVpcEncryptionControl(
	ctx context.Context, params *ec2.DeleteVpcEncryptionControlInput, _ ...func(*ec2.Options),
) (*ec2.DeleteVpcEncryptionControlOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteVpcEncryptionControlOutput), args.Error(1)
}

func (m *mockEC2Client) DescribePublicIpv4Pools(
	ctx context.Context, params *ec2.DescribePublicIpv4PoolsInput, _ ...func(*ec2.Options),
) (*ec2.DescribePublicIpv4PoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribePublicIpv4PoolsOutput), args.Error(1)
}

func (m *mockEC2Client) DeletePublicIpv4Pool(
	ctx context.Context, params *ec2.DeletePublicIpv4PoolInput, _ ...func(*ec2.Options),
) (*ec2.DeletePublicIpv4PoolOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeletePublicIpv4PoolOutput), args.Error(1)
}

func (m *mockEC2Client) DescribeVpnConcentrators(
	ctx context.Context, params *ec2.DescribeVpnConcentratorsInput, _ ...func(*ec2.Options),
) (*ec2.DescribeVpnConcentratorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeVpnConcentratorsOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteVpnConcentrator(
	ctx context.Context, params *ec2.DeleteVpnConcentratorInput, _ ...func(*ec2.Options),
) (*ec2.DeleteVpnConcentratorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteVpnConcentratorOutput), args.Error(1)
}

func (m *mockEC2Client) DescribeIpamPrefixListResolvers(
	ctx context.Context, params *ec2.DescribeIpamPrefixListResolversInput, _ ...func(*ec2.Options),
) (*ec2.DescribeIpamPrefixListResolversOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeIpamPrefixListResolversOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteIpamPrefixListResolver(
	ctx context.Context, params *ec2.DeleteIpamPrefixListResolverInput, _ ...func(*ec2.Options),
) (*ec2.DeleteIpamPrefixListResolverOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteIpamPrefixListResolverOutput), args.Error(1)
}

func (m *mockEC2Client) DescribeTrafficMirrorFilterRules(
	ctx context.Context, params *ec2.DescribeTrafficMirrorFilterRulesInput, _ ...func(*ec2.Options),
) (*ec2.DescribeTrafficMirrorFilterRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DescribeTrafficMirrorFilterRulesOutput), args.Error(1)
}

func (m *mockEC2Client) DeleteTrafficMirrorFilterRule(
	ctx context.Context, params *ec2.DeleteTrafficMirrorFilterRuleInput, _ ...func(*ec2.Options),
) (*ec2.DeleteTrafficMirrorFilterRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ec2.DeleteTrafficMirrorFilterRuleOutput), args.Error(1)
}

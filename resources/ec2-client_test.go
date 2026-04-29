//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

type ec2IntegrationBase struct {
	suite.Suite
	svc *ec2.Client
	cfg aws.Config
	ts  string
}

func (b *ec2IntegrationBase) listerOpts() *nuke.ListerOpts {
	return &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-west-2",
		},
		Config: &b.cfg,
		Logger: logrus.WithField("test", "ec2-integration"),
	}
}

func (b *ec2IntegrationBase) setup() {
	b.ts = fmt.Sprintf("%d", time.Now().UnixNano())
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		b.T().Fatalf("failed to load config: %v", err)
	}
	b.cfg = cfg
	b.svc = ec2.NewFromConfig(cfg)
}

// ---------------------------------------------------------------------------
// Suite 1: Managed Prefix List (standalone, no dependencies)
// ---------------------------------------------------------------------------

type TestEC2ManagedPrefixListSuite struct {
	ec2IntegrationBase
	prefixListId *string
}

func (suite *TestEC2ManagedPrefixListSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	resp, err := suite.svc.CreateManagedPrefixList(ctx, &ec2.CreateManagedPrefixListInput{
		PrefixListName: aws.String(fmt.Sprintf("aws-nuke-test-pl-%s", suite.ts)),
		AddressFamily:  aws.String("IPv4"),
		MaxEntries:     aws.Int32(5),
	})
	if err != nil {
		suite.T().Fatalf("failed to create prefix list: %v", err)
	}
	suite.prefixListId = resp.PrefixList.PrefixListId
}

func (suite *TestEC2ManagedPrefixListSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.prefixListId != nil {
		_, _ = suite.svc.DeleteManagedPrefixList(ctx, &ec2.DeleteManagedPrefixListInput{
			PrefixListId: suite.prefixListId,
		})
	}
}

func (suite *TestEC2ManagedPrefixListSuite) TestList() {
	assertions := assert.New(suite.T())

	lister := EC2ManagedPrefixListLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		pl := r.(*EC2ManagedPrefixList)
		if *pl.PrefixListID == *suite.prefixListId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find prefix list %s", *suite.prefixListId)
}

func (suite *TestEC2ManagedPrefixListSuite) TestRemove() {
	assertions := assert.New(suite.T())

	prefixList := &EC2ManagedPrefixList{
		svc:          suite.svc,
		PrefixListID: suite.prefixListId,
	}

	err := prefixList.Remove(context.TODO())
	assertions.NoError(err)

	suite.prefixListId = nil
}

func TestEC2ManagedPrefixListIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2ManagedPrefixListSuite))
}

// ---------------------------------------------------------------------------
// Suite 2: VPC-based resources (Flow Log, VPC CIDR Block)
// ---------------------------------------------------------------------------

type TestEC2VPCResourcesSuite struct {
	ec2IntegrationBase
	vpcId       *string
	flowLogId   *string
	cidrAssocId *string
}

func (suite *TestEC2VPCResourcesSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	// Create VPC
	vpcResp, err := suite.svc.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.200.0.0/16"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVpc,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-vpc-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create VPC: %v", err)
	}
	suite.vpcId = vpcResp.Vpc.VpcId

	// Create Flow Log (to CloudWatch Logs, no destination needed for reject-only)
	flResp, err := suite.svc.CreateFlowLogs(ctx, &ec2.CreateFlowLogsInput{
		ResourceIds:              []string{*suite.vpcId},
		ResourceType:             ec2types.FlowLogsResourceTypeVpc,
		TrafficType:              ec2types.TrafficTypeReject,
		LogDestinationType:       ec2types.LogDestinationTypeCloudWatchLogs,
		LogGroupName:             aws.String(fmt.Sprintf("/aws-nuke-test/flow-log-%s", suite.ts)),
		DeliverLogsPermissionArn: aws.String("arn:aws:iam::role/flowlogsRole"),
	})
	if err != nil {
		// Flow log creation may fail without proper IAM role — that's OK, we skip the test
		suite.T().Logf("warning: failed to create flow log (IAM role may be missing): %v", err)
	} else if len(flResp.FlowLogIds) > 0 {
		suite.flowLogId = aws.String(flResp.FlowLogIds[0])
	}

	// Associate secondary CIDR block
	cidrResp, err := suite.svc.AssociateVpcCidrBlock(ctx, &ec2.AssociateVpcCidrBlockInput{
		VpcId:     suite.vpcId,
		CidrBlock: aws.String("10.201.0.0/16"),
	})
	if err != nil {
		suite.T().Fatalf("failed to associate secondary CIDR: %v", err)
	}
	suite.cidrAssocId = cidrResp.CidrBlockAssociation.AssociationId
}

func (suite *TestEC2VPCResourcesSuite) TearDownSuite() {
	ctx := context.TODO()

	if suite.flowLogId != nil {
		_, _ = suite.svc.DeleteFlowLogs(ctx, &ec2.DeleteFlowLogsInput{
			FlowLogIds: []string{*suite.flowLogId},
		})
	}

	if suite.cidrAssocId != nil {
		_, _ = suite.svc.DisassociateVpcCidrBlock(ctx, &ec2.DisassociateVpcCidrBlockInput{
			AssociationId: suite.cidrAssocId,
		})
	}

	if suite.vpcId != nil {
		_, _ = suite.svc.DeleteVpc(ctx, &ec2.DeleteVpcInput{
			VpcId: suite.vpcId,
		})
	}
}

func (suite *TestEC2VPCResourcesSuite) TestListFlowLogs() {
	if suite.flowLogId == nil {
		suite.T().Skip("flow log was not created (IAM role may be missing)")
	}

	assertions := assert.New(suite.T())

	lister := EC2FlowLogLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		fl := r.(*EC2FlowLog)
		if *fl.FlowLogID == *suite.flowLogId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find flow log %s", *suite.flowLogId)
}

func (suite *TestEC2VPCResourcesSuite) TestListVPCCIDRBlocks() {
	assertions := assert.New(suite.T())

	lister := EC2VPCCIDRBlockLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		cidr := r.(*EC2VPCCIDRBlock)
		if cidr.AssociationId != nil && *cidr.AssociationId == *suite.cidrAssocId {
			found = true
			assertions.Equal("10.201.0.0/16", *cidr.CidrBlock)
			assertions.Equal(false, *cidr.IsIPv6)
			break
		}
	}
	assertions.True(found, "expected to find CIDR association %s", *suite.cidrAssocId)
}

func (suite *TestEC2VPCResourcesSuite) TestRemoveVPCCIDRBlock() {
	assertions := assert.New(suite.T())

	cidrBlock := &EC2VPCCIDRBlock{
		svc:           suite.svc,
		VpcId:         suite.vpcId,
		AssociationId: suite.cidrAssocId,
	}

	err := cidrBlock.Remove(context.TODO())
	assertions.NoError(err)

	suite.cidrAssocId = nil
}

func (suite *TestEC2VPCResourcesSuite) TestRemoveFlowLog() {
	if suite.flowLogId == nil {
		suite.T().Skip("flow log was not created (IAM role may be missing)")
	}

	assertions := assert.New(suite.T())

	flowLog := &EC2FlowLog{
		svc:       suite.svc,
		FlowLogId: suite.flowLogId,
	}

	err := flowLog.Remove(context.TODO())
	assertions.NoError(err)

	suite.flowLogId = nil
}

func TestEC2VPCResourcesIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2VPCResourcesSuite))
}

// ---------------------------------------------------------------------------
// Suite 3: IPAM resources (IPAM, Scope, Pool)
// ---------------------------------------------------------------------------

type TestEC2IPAMSuite struct {
	ec2IntegrationBase
	ipamId      *string
	ipamScopeId *string
	ipamPoolId  *string
}

func (suite *TestEC2IPAMSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	// Create IPAM
	ipamResp, err := suite.svc.CreateIpam(ctx, &ec2.CreateIpamInput{
		OperatingRegions: []ec2types.AddIpamOperatingRegion{
			{RegionName: aws.String("us-west-2")},
		},
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeIpam,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-ipam-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create IPAM: %v", err)
	}
	suite.ipamId = ipamResp.Ipam.IpamId

	// Find the private scope that was auto-created
	scopeResp, err := suite.svc.DescribeIpamScopes(ctx, &ec2.DescribeIpamScopesInput{
		Filters: []ec2types.Filter{
			{Name: aws.String("ipam-id"), Values: []string{*suite.ipamId}},
			{Name: aws.String("is-default"), Values: []string{"false"}},
		},
	})
	if err == nil && len(scopeResp.IpamScopes) > 0 {
		suite.ipamScopeId = scopeResp.IpamScopes[0].IpamScopeId
	}

	// If no non-default scope found, use a default scope for pool creation
	if suite.ipamScopeId == nil {
		scopeResp2, err := suite.svc.DescribeIpamScopes(ctx, &ec2.DescribeIpamScopesInput{
			Filters: []ec2types.Filter{
				{Name: aws.String("ipam-id"), Values: []string{*suite.ipamId}},
			},
		})
		if err != nil || len(scopeResp2.IpamScopes) == 0 {
			suite.T().Fatalf("failed to find any IPAM scope: %v", err)
		}
		suite.ipamScopeId = scopeResp2.IpamScopes[0].IpamScopeId
	}

	// Create IPAM Pool
	poolResp, err := suite.svc.CreateIpamPool(ctx, &ec2.CreateIpamPoolInput{
		IpamScopeId:   suite.ipamScopeId,
		AddressFamily: ec2types.AddressFamilyIpv4,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeIpamPool,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-pool-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create IPAM pool: %v", err)
	}
	suite.ipamPoolId = poolResp.IpamPool.IpamPoolId
}

func (suite *TestEC2IPAMSuite) TearDownSuite() {
	ctx := context.TODO()

	if suite.ipamPoolId != nil {
		_, _ = suite.svc.DeleteIpamPool(ctx, &ec2.DeleteIpamPoolInput{
			IpamPoolId: suite.ipamPoolId,
		})
	}

	if suite.ipamId != nil {
		_, _ = suite.svc.DeleteIpam(ctx, &ec2.DeleteIpamInput{
			IpamId:  suite.ipamId,
			Cascade: aws.Bool(true),
		})
	}
}

func (suite *TestEC2IPAMSuite) TestListIPAMs() {
	assertions := assert.New(suite.T())

	lister := EC2IPAMLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		ipam := r.(*EC2IPAM)
		if *ipam.IpamId == *suite.ipamId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find IPAM %s", *suite.ipamId)
}

func (suite *TestEC2IPAMSuite) TestListIPAMScopes() {
	assertions := assert.New(suite.T())

	lister := EC2IPAMScopeLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		scope := r.(*EC2IPAMScope)
		if *scope.IpamScopeId == *suite.ipamScopeId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find IPAM scope %s", *suite.ipamScopeId)
}

func (suite *TestEC2IPAMSuite) TestListIPAMPools() {
	assertions := assert.New(suite.T())

	lister := EC2IPAMPoolLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		pool := r.(*EC2IPAMPool)
		if *pool.IpamPoolId == *suite.ipamPoolId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find IPAM pool %s", *suite.ipamPoolId)
}

func (suite *TestEC2IPAMSuite) TestRemoveIPAMPool() {
	assertions := assert.New(suite.T())

	pool := &EC2IPAMPool{
		svc:        suite.svc,
		IpamPoolId: suite.ipamPoolId,
	}

	err := pool.Remove(context.TODO())
	assertions.NoError(err)

	suite.ipamPoolId = nil
}

func (suite *TestEC2IPAMSuite) TestRemoveIPAM() {
	assertions := assert.New(suite.T())

	ipam := &EC2IPAM{
		svc:    suite.svc,
		IpamId: suite.ipamId,
	}

	err := ipam.Remove(context.TODO())
	assertions.NoError(err)

	suite.ipamId = nil
}

func TestEC2IPAMIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2IPAMSuite))
}

// ---------------------------------------------------------------------------
// Suite 4: Transit Gateway resources (Route Table, Multicast Domain, Policy Table, Prefix List Reference)
// ---------------------------------------------------------------------------

type TestEC2TGWSuite struct {
	ec2IntegrationBase
	tgwId         *string
	routeTableId  *string
	multicastId   *string
	policyTableId *string
	prefixListId  *string
}

func (suite *TestEC2TGWSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	// Create Transit Gateway
	tgwResp, err := suite.svc.CreateTransitGateway(ctx, &ec2.CreateTransitGatewayInput{
		Description: aws.String(fmt.Sprintf("aws-nuke-test-tgw-%s", suite.ts)),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeTransitGateway,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-tgw-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create TGW: %v", err)
	}
	suite.tgwId = tgwResp.TransitGateway.TransitGatewayId

	// Wait for TGW to become available
	suite.T().Logf("waiting for TGW %s to become available...", *suite.tgwId)
	for i := 0; i < 60; i++ {
		descResp, err := suite.svc.DescribeTransitGateways(ctx, &ec2.DescribeTransitGatewaysInput{
			TransitGatewayIds: []string{*suite.tgwId},
		})
		if err == nil && len(descResp.TransitGateways) > 0 &&
			descResp.TransitGateways[0].State == ec2types.TransitGatewayStateAvailable {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Create TGW Route Table
	rtResp, err := suite.svc.CreateTransitGatewayRouteTable(ctx, &ec2.CreateTransitGatewayRouteTableInput{
		TransitGatewayId: suite.tgwId,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeTransitGatewayRouteTable,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-rt-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create TGW route table: %v", err)
	}
	suite.routeTableId = rtResp.TransitGatewayRouteTable.TransitGatewayRouteTableId

	// Create TGW Multicast Domain
	mdResp, err := suite.svc.CreateTransitGatewayMulticastDomain(ctx, &ec2.CreateTransitGatewayMulticastDomainInput{
		TransitGatewayId: suite.tgwId,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeTransitGatewayMulticastDomain,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-md-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create TGW multicast domain: %v", err)
	}
	suite.multicastId = mdResp.TransitGatewayMulticastDomain.TransitGatewayMulticastDomainId

	// Create TGW Policy Table
	ptResp, err := suite.svc.CreateTransitGatewayPolicyTable(ctx, &ec2.CreateTransitGatewayPolicyTableInput{
		TransitGatewayId: suite.tgwId,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeTransitGatewayPolicyTable,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-pt-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create TGW policy table: %v", err)
	}
	suite.policyTableId = ptResp.TransitGatewayPolicyTable.TransitGatewayPolicyTableId

	// Create a managed prefix list for the prefix list reference test
	plResp, err := suite.svc.CreateManagedPrefixList(ctx, &ec2.CreateManagedPrefixListInput{
		PrefixListName: aws.String(fmt.Sprintf("aws-nuke-test-tgw-pl-%s", suite.ts)),
		AddressFamily:  aws.String("IPv4"),
		MaxEntries:     aws.Int32(5),
		Entries: []ec2types.AddPrefixListEntry{
			{Cidr: aws.String("10.0.0.0/8")},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create prefix list: %v", err)
	}
	suite.prefixListId = plResp.PrefixList.PrefixListId

	// Create prefix list reference on the route table
	_, err = suite.svc.CreateTransitGatewayPrefixListReference(ctx, &ec2.CreateTransitGatewayPrefixListReferenceInput{
		TransitGatewayRouteTableId: suite.routeTableId,
		PrefixListId:               suite.prefixListId,
		Blackhole:                  aws.Bool(true),
	})
	if err != nil {
		suite.T().Logf("warning: failed to create prefix list reference: %v", err)
	}
}

func (suite *TestEC2TGWSuite) TearDownSuite() {
	ctx := context.TODO()

	// Delete prefix list reference
	if suite.routeTableId != nil && suite.prefixListId != nil {
		_, _ = suite.svc.DeleteTransitGatewayPrefixListReference(ctx, &ec2.DeleteTransitGatewayPrefixListReferenceInput{
			TransitGatewayRouteTableId: suite.routeTableId,
			PrefixListId:               suite.prefixListId,
		})
	}

	// Delete prefix list
	if suite.prefixListId != nil {
		_, _ = suite.svc.DeleteManagedPrefixList(ctx, &ec2.DeleteManagedPrefixListInput{
			PrefixListId: suite.prefixListId,
		})
	}

	// Delete TGW sub-resources
	if suite.policyTableId != nil {
		_, _ = suite.svc.DeleteTransitGatewayPolicyTable(ctx, &ec2.DeleteTransitGatewayPolicyTableInput{
			TransitGatewayPolicyTableId: suite.policyTableId,
		})
	}
	if suite.multicastId != nil {
		_, _ = suite.svc.DeleteTransitGatewayMulticastDomain(ctx, &ec2.DeleteTransitGatewayMulticastDomainInput{
			TransitGatewayMulticastDomainId: suite.multicastId,
		})
	}
	if suite.routeTableId != nil {
		_, _ = suite.svc.DeleteTransitGatewayRouteTable(ctx, &ec2.DeleteTransitGatewayRouteTableInput{
			TransitGatewayRouteTableId: suite.routeTableId,
		})
	}

	// Delete TGW
	if suite.tgwId != nil {
		_, _ = suite.svc.DeleteTransitGateway(ctx, &ec2.DeleteTransitGatewayInput{
			TransitGatewayId: suite.tgwId,
		})
	}
}

func (suite *TestEC2TGWSuite) TestListRouteTable() {
	assertions := assert.New(suite.T())

	lister := EC2TGWRouteTableLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		rt := r.(*EC2TGWRouteTable)
		if *rt.TransitGatewayRouteTableId == *suite.routeTableId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find TGW route table %s", *suite.routeTableId)
}

func (suite *TestEC2TGWSuite) TestListMulticastDomain() {
	assertions := assert.New(suite.T())

	lister := EC2TGWMulticastDomainLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		md := r.(*EC2TGWMulticastDomain)
		if *md.TransitGatewayMulticastDomainId == *suite.multicastId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find TGW multicast domain %s", *suite.multicastId)
}

func (suite *TestEC2TGWSuite) TestListPolicyTable() {
	assertions := assert.New(suite.T())

	lister := EC2TGWPolicyTableLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		pt := r.(*EC2TGWPolicyTable)
		if *pt.TransitGatewayPolicyTableId == *suite.policyTableId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find TGW policy table %s", *suite.policyTableId)
}

func (suite *TestEC2TGWSuite) TestListPrefixListReference() {
	assertions := assert.New(suite.T())

	lister := EC2TGWPrefixListReferenceLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		plr := r.(*EC2TGWPrefixListReference)
		if *plr.TransitGatewayRouteTableId == *suite.routeTableId &&
			*plr.PrefixListId == *suite.prefixListId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find prefix list reference for %s on %s",
		*suite.prefixListId, *suite.routeTableId)
}

func (suite *TestEC2TGWSuite) TestRemovePrefixListReference() {
	assertions := assert.New(suite.T())

	ref := &EC2TGWPrefixListReference{
		svc:                        suite.svc,
		TransitGatewayRouteTableId: suite.routeTableId,
		PrefixListId:               suite.prefixListId,
	}

	err := ref.Remove(context.TODO())
	assertions.NoError(err)
}

func (suite *TestEC2TGWSuite) TestRemoveRouteTable() {
	assertions := assert.New(suite.T())

	rt := &EC2TGWRouteTable{
		svc:                        suite.svc,
		TransitGatewayRouteTableId: suite.routeTableId,
	}

	err := rt.Remove(context.TODO())
	assertions.NoError(err)

	suite.routeTableId = nil
}

func (suite *TestEC2TGWSuite) TestRemoveMulticastDomain() {
	assertions := assert.New(suite.T())

	md := &EC2TGWMulticastDomain{
		svc:                             suite.svc,
		TransitGatewayMulticastDomainId: suite.multicastId,
	}

	err := md.Remove(context.TODO())
	assertions.NoError(err)

	suite.multicastId = nil
}

func (suite *TestEC2TGWSuite) TestRemovePolicyTable() {
	assertions := assert.New(suite.T())

	pt := &EC2TGWPolicyTable{
		svc:                         suite.svc,
		TransitGatewayPolicyTableId: suite.policyTableId,
	}

	err := pt.Remove(context.TODO())
	assertions.NoError(err)

	suite.policyTableId = nil
}

func TestEC2TGWIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2TGWSuite))
}

// ---------------------------------------------------------------------------
// Suite 5: Traffic Mirror resources (Filter, Target, Session)
// ---------------------------------------------------------------------------

type TestEC2TrafficMirrorSuite struct {
	ec2IntegrationBase
	vpcId     *string
	subnetId  *string
	eniId     *string
	eniId2    *string
	filterId  *string
	targetId  *string
	sessionId *string
}

func (suite *TestEC2TrafficMirrorSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	// Create VPC
	vpcResp, err := suite.svc.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.210.0.0/16"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create VPC: %v", err)
	}
	suite.vpcId = vpcResp.Vpc.VpcId

	// Create subnet
	subnetResp, err := suite.svc.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:     suite.vpcId,
		CidrBlock: aws.String("10.210.1.0/24"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create subnet: %v", err)
	}
	suite.subnetId = subnetResp.Subnet.SubnetId

	// Create two ENIs (one for source, one for target)
	eni1Resp, err := suite.svc.CreateNetworkInterface(ctx, &ec2.CreateNetworkInterfaceInput{
		SubnetId: suite.subnetId,
	})
	if err != nil {
		suite.T().Fatalf("failed to create ENI 1: %v", err)
	}
	suite.eniId = eni1Resp.NetworkInterface.NetworkInterfaceId

	eni2Resp, err := suite.svc.CreateNetworkInterface(ctx, &ec2.CreateNetworkInterfaceInput{
		SubnetId: suite.subnetId,
	})
	if err != nil {
		suite.T().Fatalf("failed to create ENI 2: %v", err)
	}
	suite.eniId2 = eni2Resp.NetworkInterface.NetworkInterfaceId

	// Create traffic mirror filter
	filterResp, err := suite.svc.CreateTrafficMirrorFilter(ctx, &ec2.CreateTrafficMirrorFilterInput{
		Description: aws.String(fmt.Sprintf("aws-nuke-test-filter-%s", suite.ts)),
	})
	if err != nil {
		suite.T().Fatalf("failed to create traffic mirror filter: %v", err)
	}
	suite.filterId = filterResp.TrafficMirrorFilter.TrafficMirrorFilterId

	// Create traffic mirror target (ENI-based)
	targetResp, err := suite.svc.CreateTrafficMirrorTarget(ctx, &ec2.CreateTrafficMirrorTargetInput{
		NetworkInterfaceId: suite.eniId2,
		Description:        aws.String(fmt.Sprintf("aws-nuke-test-target-%s", suite.ts)),
	})
	if err != nil {
		suite.T().Fatalf("failed to create traffic mirror target: %v", err)
	}
	suite.targetId = targetResp.TrafficMirrorTarget.TrafficMirrorTargetId

	// Create traffic mirror session
	sessionResp, err := suite.svc.CreateTrafficMirrorSession(ctx, &ec2.CreateTrafficMirrorSessionInput{
		NetworkInterfaceId:    suite.eniId,
		TrafficMirrorTargetId: suite.targetId,
		TrafficMirrorFilterId: suite.filterId,
		SessionNumber:         aws.Int32(1),
		Description:           aws.String(fmt.Sprintf("aws-nuke-test-session-%s", suite.ts)),
	})
	if err != nil {
		suite.T().Fatalf("failed to create traffic mirror session: %v", err)
	}
	suite.sessionId = sessionResp.TrafficMirrorSession.TrafficMirrorSessionId
}

func (suite *TestEC2TrafficMirrorSuite) TearDownSuite() {
	ctx := context.TODO()

	if suite.sessionId != nil {
		_, _ = suite.svc.DeleteTrafficMirrorSession(ctx, &ec2.DeleteTrafficMirrorSessionInput{
			TrafficMirrorSessionId: suite.sessionId,
		})
	}
	if suite.targetId != nil {
		_, _ = suite.svc.DeleteTrafficMirrorTarget(ctx, &ec2.DeleteTrafficMirrorTargetInput{
			TrafficMirrorTargetId: suite.targetId,
		})
	}
	if suite.filterId != nil {
		_, _ = suite.svc.DeleteTrafficMirrorFilter(ctx, &ec2.DeleteTrafficMirrorFilterInput{
			TrafficMirrorFilterId: suite.filterId,
		})
	}
	if suite.eniId != nil {
		_, _ = suite.svc.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: suite.eniId,
		})
	}
	if suite.eniId2 != nil {
		_, _ = suite.svc.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: suite.eniId2,
		})
	}
	if suite.subnetId != nil {
		_, _ = suite.svc.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
			SubnetId: suite.subnetId,
		})
	}
	if suite.vpcId != nil {
		_, _ = suite.svc.DeleteVpc(ctx, &ec2.DeleteVpcInput{
			VpcId: suite.vpcId,
		})
	}
}

func (suite *TestEC2TrafficMirrorSuite) TestListFilters() {
	assertions := assert.New(suite.T())

	lister := EC2TrafficMirrorFilterLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		f := r.(*EC2TrafficMirrorFilter)
		if *f.TrafficMirrorFilterId == *suite.filterId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find traffic mirror filter %s", *suite.filterId)
}

func (suite *TestEC2TrafficMirrorSuite) TestListTargets() {
	assertions := assert.New(suite.T())

	lister := EC2TrafficMirrorTargetLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		tgt := r.(*EC2TrafficMirrorTarget)
		if *tgt.TrafficMirrorTargetId == *suite.targetId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find traffic mirror target %s", *suite.targetId)
}

func (suite *TestEC2TrafficMirrorSuite) TestListSessions() {
	assertions := assert.New(suite.T())

	lister := EC2TrafficMirrorSessionLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		s := r.(*EC2TrafficMirrorSession)
		if *s.TrafficMirrorSessionId == *suite.sessionId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find traffic mirror session %s", *suite.sessionId)
}

func (suite *TestEC2TrafficMirrorSuite) TestRemoveSession() {
	assertions := assert.New(suite.T())

	session := &EC2TrafficMirrorSession{
		svc:                    suite.svc,
		TrafficMirrorSessionId: suite.sessionId,
	}

	err := session.Remove(context.TODO())
	assertions.NoError(err)

	suite.sessionId = nil
}

func (suite *TestEC2TrafficMirrorSuite) TestRemoveTarget() {
	assertions := assert.New(suite.T())

	target := &EC2TrafficMirrorTarget{
		svc:                   suite.svc,
		TrafficMirrorTargetId: suite.targetId,
	}

	err := target.Remove(context.TODO())
	assertions.NoError(err)

	suite.targetId = nil
}

func (suite *TestEC2TrafficMirrorSuite) TestRemoveFilter() {
	assertions := assert.New(suite.T())

	filter := &EC2TrafficMirrorFilter{
		svc:                   suite.svc,
		TrafficMirrorFilterId: suite.filterId,
	}

	err := filter.Remove(context.TODO())
	assertions.NoError(err)

	suite.filterId = nil
}

func TestEC2TrafficMirrorIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2TrafficMirrorSuite))
}

// ---------------------------------------------------------------------------
// Suite 6: Network Insights resources (Path, Analysis)
// ---------------------------------------------------------------------------

type TestEC2NetworkInsightsSuite struct {
	ec2IntegrationBase
	vpcId      *string
	subnetId   *string
	eniId      *string
	eniId2     *string
	pathId     *string
	analysisId *string
}

func (suite *TestEC2NetworkInsightsSuite) SetupSuite() {
	suite.setup()
	ctx := context.TODO()

	// Create VPC
	vpcResp, err := suite.svc.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.220.0.0/16"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create VPC: %v", err)
	}
	suite.vpcId = vpcResp.Vpc.VpcId

	// Create subnet
	subnetResp, err := suite.svc.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:     suite.vpcId,
		CidrBlock: aws.String("10.220.1.0/24"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create subnet: %v", err)
	}
	suite.subnetId = subnetResp.Subnet.SubnetId

	// Create two ENIs
	eni1Resp, err := suite.svc.CreateNetworkInterface(ctx, &ec2.CreateNetworkInterfaceInput{
		SubnetId: suite.subnetId,
	})
	if err != nil {
		suite.T().Fatalf("failed to create ENI 1: %v", err)
	}
	suite.eniId = eni1Resp.NetworkInterface.NetworkInterfaceId

	eni2Resp, err := suite.svc.CreateNetworkInterface(ctx, &ec2.CreateNetworkInterfaceInput{
		SubnetId: suite.subnetId,
	})
	if err != nil {
		suite.T().Fatalf("failed to create ENI 2: %v", err)
	}
	suite.eniId2 = eni2Resp.NetworkInterface.NetworkInterfaceId

	// Create Network Insights Path
	pathResp, err := suite.svc.CreateNetworkInsightsPath(ctx, &ec2.CreateNetworkInsightsPathInput{
		Source:      suite.eniId,
		Destination: suite.eniId2,
		Protocol:    ec2types.ProtocolTcp,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeNetworkInsightsPath,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-path-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create network insights path: %v", err)
	}
	suite.pathId = pathResp.NetworkInsightsPath.NetworkInsightsPathId

	// Start a Network Insights Analysis
	analysisResp, err := suite.svc.StartNetworkInsightsAnalysis(ctx, &ec2.StartNetworkInsightsAnalysisInput{
		NetworkInsightsPathId: suite.pathId,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeNetworkInsightsAnalysis,
				Tags: []ec2types.Tag{
					{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("aws-nuke-test-analysis-%s", suite.ts))},
				},
			},
		},
	})
	if err != nil {
		suite.T().Logf("warning: failed to start network insights analysis: %v", err)
	} else {
		suite.analysisId = analysisResp.NetworkInsightsAnalysis.NetworkInsightsAnalysisId
	}
}

func (suite *TestEC2NetworkInsightsSuite) TearDownSuite() {
	ctx := context.TODO()

	if suite.analysisId != nil {
		_, _ = suite.svc.DeleteNetworkInsightsAnalysis(ctx, &ec2.DeleteNetworkInsightsAnalysisInput{
			NetworkInsightsAnalysisId: suite.analysisId,
		})
	}
	if suite.pathId != nil {
		_, _ = suite.svc.DeleteNetworkInsightsPath(ctx, &ec2.DeleteNetworkInsightsPathInput{
			NetworkInsightsPathId: suite.pathId,
		})
	}
	if suite.eniId != nil {
		_, _ = suite.svc.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: suite.eniId,
		})
	}
	if suite.eniId2 != nil {
		_, _ = suite.svc.DeleteNetworkInterface(ctx, &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: suite.eniId2,
		})
	}
	if suite.subnetId != nil {
		_, _ = suite.svc.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{
			SubnetId: suite.subnetId,
		})
	}
	if suite.vpcId != nil {
		_, _ = suite.svc.DeleteVpc(ctx, &ec2.DeleteVpcInput{
			VpcId: suite.vpcId,
		})
	}
}

func (suite *TestEC2NetworkInsightsSuite) TestListPaths() {
	assertions := assert.New(suite.T())

	lister := EC2NetworkInsightsPathLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		p := r.(*EC2NetworkInsightsPath)
		if *p.NetworkInsightsPathId == *suite.pathId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find network insights path %s", *suite.pathId)
}

func (suite *TestEC2NetworkInsightsSuite) TestListAnalyses() {
	if suite.analysisId == nil {
		suite.T().Skip("analysis was not created")
	}

	assertions := assert.New(suite.T())

	lister := EC2NetworkInsightsAnalysisLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	assertions.NoError(err)

	found := false
	for _, r := range resources {
		a := r.(*EC2NetworkInsightsAnalysis)
		if *a.NetworkInsightsAnalysisId == *suite.analysisId {
			found = true
			break
		}
	}
	assertions.True(found, "expected to find network insights analysis %s", *suite.analysisId)
}

func (suite *TestEC2NetworkInsightsSuite) TestRemoveAnalysis() {
	if suite.analysisId == nil {
		suite.T().Skip("analysis was not created")
	}

	assertions := assert.New(suite.T())

	analysis := &EC2NetworkInsightsAnalysis{
		svc:                       suite.svc,
		NetworkInsightsAnalysisId: suite.analysisId,
	}

	err := analysis.Remove(context.TODO())
	assertions.NoError(err)

	suite.analysisId = nil
}

func (suite *TestEC2NetworkInsightsSuite) TestRemovePath() {
	assertions := assert.New(suite.T())

	path := &EC2NetworkInsightsPath{
		svc:                   suite.svc,
		NetworkInsightsPathId: suite.pathId,
	}

	err := path.Remove(context.TODO())
	assertions.NoError(err)

	suite.pathId = nil
}

func TestEC2NetworkInsightsIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2NetworkInsightsSuite))
}

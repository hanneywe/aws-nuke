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
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	latticetypes "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestVPCLatticeSuite struct {
	suite.Suite
	svc                *vpclattice.Client
	cfg                aws.Config
	serviceNetworkName string
	serviceNetworkARN  string
	serviceName        string
	serviceARN         string
	targetGroupName    string
	targetGroupARN     string
}

func (suite *TestVPCLatticeSuite) listerOpts() *nuke.ListerOpts {
	return &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-west-2",
		},
		Config: &suite.cfg,
		Logger: logrus.WithField("test", "vpclattice"),
	}
}

func (suite *TestVPCLatticeSuite) SetupSuite() {
	ts := time.Now().UnixNano()
	suite.serviceNetworkName = fmt.Sprintf("aws-nuke-test-sn-%d", ts)
	suite.serviceName = fmt.Sprintf("aws-nuke-test-svc-%d", ts)
	suite.targetGroupName = fmt.Sprintf("aws-nuke-test-tg-%d", ts)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		suite.T().Fatalf("failed to create config: %v", err)
	}
	suite.cfg = cfg
	suite.svc = vpclattice.NewFromConfig(cfg)

	// Create service network
	snResp, err := suite.svc.CreateServiceNetwork(ctx, &vpclattice.CreateServiceNetworkInput{
		Name: aws.String(suite.serviceNetworkName),
	})
	if err != nil {
		suite.T().Fatalf("failed to create service network: %v", err)
	}
	suite.serviceNetworkARN = *snResp.Arn

	// Create service
	svcResp, err := suite.svc.CreateService(ctx, &vpclattice.CreateServiceInput{
		Name: aws.String(suite.serviceName),
	})
	if err != nil {
		suite.T().Fatalf("failed to create service: %v", err)
	}
	suite.serviceARN = *svcResp.Arn

	// Create target group
	tgResp, err := suite.svc.CreateTargetGroup(ctx, &vpclattice.CreateTargetGroupInput{
		Name: aws.String(suite.targetGroupName),
		Type: latticetypes.TargetGroupTypeLambda,
	})
	if err != nil {
		suite.T().Fatalf("failed to create target group: %v", err)
	}
	suite.targetGroupARN = *tgResp.Arn
}

func (suite *TestVPCLatticeSuite) TearDownSuite() {
	ctx := context.TODO()

	// Clean up target group
	_, _ = suite.svc.DeleteTargetGroup(ctx, &vpclattice.DeleteTargetGroupInput{
		TargetGroupIdentifier: aws.String(suite.targetGroupARN),
	})

	// Clean up service
	_, _ = suite.svc.DeleteService(ctx, &vpclattice.DeleteServiceInput{
		ServiceIdentifier: aws.String(suite.serviceARN),
	})

	// Clean up service network
	_, _ = suite.svc.DeleteServiceNetwork(ctx, &vpclattice.DeleteServiceNetworkInput{
		ServiceNetworkIdentifier: aws.String(suite.serviceNetworkARN),
	})
}

func (suite *TestVPCLatticeSuite) TestListServiceNetworks() {
	a := assert.New(suite.T())

	lister := VPCLatticeServiceNetworkLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	a.Nil(err)

	found := false
	for _, r := range resources {
		sn := r.(*VPCLatticeServiceNetwork)
		if *sn.Name == suite.serviceNetworkName {
			found = true
			break
		}
	}
	a.True(found, "expected to find service network %s", suite.serviceNetworkName)
}

func (suite *TestVPCLatticeSuite) TestListServices() {
	a := assert.New(suite.T())

	lister := VPCLatticeServiceLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	a.Nil(err)

	found := false
	for _, r := range resources {
		svc := r.(*VPCLatticeService)
		if *svc.Name == suite.serviceName {
			found = true
			break
		}
	}
	a.True(found, "expected to find service %s", suite.serviceName)
}

func (suite *TestVPCLatticeSuite) TestListTargetGroups() {
	a := assert.New(suite.T())

	lister := VPCLatticeTargetGroupLister{}
	resources, err := lister.List(context.TODO(), suite.listerOpts())
	a.Nil(err)

	found := false
	for _, r := range resources {
		tg := r.(*VPCLatticeTargetGroup)
		if *tg.Name == suite.targetGroupName {
			found = true
			break
		}
	}
	a.True(found, "expected to find target group %s", suite.targetGroupName)
}

func (suite *TestVPCLatticeSuite) TestRemoveTargetGroup() {
	a := assert.New(suite.T())

	tg := &VPCLatticeTargetGroup{
		svc:  suite.svc,
		ARN:  aws.String(suite.targetGroupARN),
		Name: aws.String(suite.targetGroupName),
	}

	err := tg.Remove(context.TODO())
	a.Nil(err)

	// Clear ARN so TearDownSuite doesn't try to delete again
	suite.targetGroupARN = ""
}

func TestVPCLatticeIntegration(t *testing.T) {
	suite.Run(t, new(TestVPCLatticeSuite))
}

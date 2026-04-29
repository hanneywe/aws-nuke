//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/backupgateway"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestBackupGatewayHypervisorSuite struct {
	suite.Suite
	svc           *backupgateway.Client
	cfg           config.Config
	hypervisorArn *string
}

func (suite *TestBackupGatewayHypervisorSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.cfg = cfg
	suite.svc = backupgateway.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())

	resp, err := suite.svc.ImportHypervisorConfiguration(ctx, &backupgateway.ImportHypervisorConfigurationInput{
		Name: ptr.String(name),
		Host: ptr.String("vsphere.example.com"),
	})
	if err != nil {
		suite.T().Fatalf("failed to import hypervisor: %v", err)
	}
	suite.hypervisorArn = resp.HypervisorArn
}

func (suite *TestBackupGatewayHypervisorSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteHypervisor(ctx, &backupgateway.DeleteHypervisorInput{
		HypervisorArn: suite.hypervisorArn,
	})
}

func (suite *TestBackupGatewayHypervisorSuite) TestList() {
	a := assert.New(suite.T())

	awsCfg := suite.cfg
	lister := BackupGatewayHypervisorLister{}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-east-1",
		},
		Config: &awsCfg,
		Logger: logrus.WithField("test", "backupgateway-hypervisor"),
	})
	a.Nil(err)
	a.Greater(len(resources), 0)
}

func (suite *TestBackupGatewayHypervisorSuite) TestRemove() {
	a := assert.New(suite.T())

	hv := &BackupGatewayHypervisor{
		svc:           suite.svc,
		HypervisorArn: suite.hypervisorArn,
	}

	err := hv.Remove(context.TODO())
	a.Nil(err)
}

func TestBackupGatewayHypervisorIntegration(t *testing.T) {
	suite.Run(t, new(TestBackupGatewayHypervisorSuite))
}

//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"
)

type TestRoute53RecoveryControlConfigControlPanelSuite struct {
	suite.Suite
	svc *route53recoverycontrolconfig.Client
}

func (suite *TestRoute53RecoveryControlConfigControlPanelSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = route53recoverycontrolconfig.NewFromConfig(cfg)
}

func (suite *TestRoute53RecoveryControlConfigControlPanelSuite) TestList() {
	a := assert.New(suite.T())
	lister := Route53RecoveryControlConfigControlPanelLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	_ = resources
}

func TestRoute53RecoveryControlConfigControlPanelIntegration(t *testing.T) {
	suite.Run(t, new(TestRoute53RecoveryControlConfigControlPanelSuite))
}

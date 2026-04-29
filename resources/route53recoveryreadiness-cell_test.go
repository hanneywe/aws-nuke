//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"
)

type TestRoute53RecoveryReadinessCellSuite struct {
	suite.Suite
	svc *route53recoveryreadiness.Client
}

func (suite *TestRoute53RecoveryReadinessCellSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = route53recoveryreadiness.NewFromConfig(cfg)
}

func (suite *TestRoute53RecoveryReadinessCellSuite) TestList() {
	a := assert.New(suite.T())
	lister := Route53RecoveryReadinessCellLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryReadinessListerOpts)
	a.NoError(err)
	_ = resources
}

func TestRoute53RecoveryReadinessCellIntegration(t *testing.T) {
	suite.Run(t, new(TestRoute53RecoveryReadinessCellSuite))
}

//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type TestECRPullThroughCacheRuleSuite struct {
	suite.Suite
	svc *ecr.Client
}

func (suite *TestECRPullThroughCacheRuleSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = ecr.NewFromConfig(cfg)
}

func (suite *TestECRPullThroughCacheRuleSuite) TestList() {
	a := assert.New(suite.T())
	lister := ECRPullThroughCacheRuleLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestECRPullThroughCacheRuleIntegration(t *testing.T) {
	suite.Run(t, new(TestECRPullThroughCacheRuleSuite))
}

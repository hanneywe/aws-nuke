//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rbin"
)

type TestRbinRuleSuite struct {
	suite.Suite
	svc *rbin.Client
}

func (s *TestRbinRuleSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = rbin.NewFromConfig(cfg)
}

func (s *TestRbinRuleSuite) TestList() {
	assertions := assert.New(s.T())

	lister := RbinRuleLister{}
	resources, err := lister.List(context.TODO(), testRbinListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestRbinRuleIntegration(t *testing.T) {
	suite.Run(t, new(TestRbinRuleSuite))
}

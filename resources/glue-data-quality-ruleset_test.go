//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
)

type TestGlueDataQualityRulesetSuite struct {
	suite.Suite
	svc *glue.Client
}

func (suite *TestGlueDataQualityRulesetSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = glue.NewFromConfig(cfg)
}

func (suite *TestGlueDataQualityRulesetSuite) TestList() {
	a := assert.New(suite.T())
	lister := GlueDataQualityRulesetLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGlueDataQualityRulesetIntegration(t *testing.T) {
	suite.Run(t, new(TestGlueDataQualityRulesetSuite))
}

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

type TestGlueUserDefinedFunctionSuite struct {
	suite.Suite
	svc *glue.Client
}

func (suite *TestGlueUserDefinedFunctionSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = glue.NewFromConfig(cfg)
}

func (suite *TestGlueUserDefinedFunctionSuite) TestList() {
	a := assert.New(suite.T())
	lister := GlueUserDefinedFunctionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGlueUserDefinedFunctionIntegration(t *testing.T) {
	suite.Run(t, new(TestGlueUserDefinedFunctionSuite))
}

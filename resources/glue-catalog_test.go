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

type TestGlueCatalogSuite struct {
	suite.Suite
	svc *glue.Client
}

func (suite *TestGlueCatalogSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = glue.NewFromConfig(cfg)
}

func (suite *TestGlueCatalogSuite) TestList() {
	a := assert.New(suite.T())
	lister := GlueCatalogLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGlueCatalogIntegration(t *testing.T) {
	suite.Run(t, new(TestGlueCatalogSuite))
}

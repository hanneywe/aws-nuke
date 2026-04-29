//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/schemas"
)

type TestSchemasRegistrySuite struct {
	suite.Suite
	svc *schemas.Client
}

func (suite *TestSchemasRegistrySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = schemas.NewFromConfig(cfg)
}

func (suite *TestSchemasRegistrySuite) TestList() {
	a := assert.New(suite.T())
	lister := SchemasRegistryLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSchemasListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSchemasRegistryIntegration(t *testing.T) {
	suite.Run(t, new(TestSchemasRegistrySuite))
}

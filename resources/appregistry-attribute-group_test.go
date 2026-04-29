//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
)

type TestAppRegistryAttributeGroupSuite struct {
	suite.Suite
	svc *servicecatalogappregistry.Client
}

func (suite *TestAppRegistryAttributeGroupSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = servicecatalogappregistry.NewFromConfig(cfg)
}

func (suite *TestAppRegistryAttributeGroupSuite) TestList() {
	a := assert.New(suite.T())
	lister := AppRegistryAttributeGroupLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAppRegistryListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAppRegistryAttributeGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestAppRegistryAttributeGroupSuite))
}

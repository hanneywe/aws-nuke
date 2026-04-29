//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/greengrass"
)

type TestGreengrassCoreDefinitionSuite struct {
	suite.Suite
	svc *greengrass.Client
}

func (suite *TestGreengrassCoreDefinitionSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = greengrass.NewFromConfig(cfg)
}

func (suite *TestGreengrassCoreDefinitionSuite) TestList() {
	a := assert.New(suite.T())
	lister := GreengrassCoreDefinitionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGreengrassCoreDefinitionIntegration(t *testing.T) {
	suite.Run(t, new(TestGreengrassCoreDefinitionSuite))
}

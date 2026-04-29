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

type TestGreengrassConnectorDefinitionSuite struct {
	suite.Suite
	svc *greengrass.Client
}

func (suite *TestGreengrassConnectorDefinitionSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = greengrass.NewFromConfig(cfg)
}

func (suite *TestGreengrassConnectorDefinitionSuite) TestList() {
	a := assert.New(suite.T())
	lister := GreengrassConnectorDefinitionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGreengrassConnectorDefinitionIntegration(t *testing.T) {
	suite.Run(t, new(TestGreengrassConnectorDefinitionSuite))
}

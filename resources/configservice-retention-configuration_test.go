//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
)

type TestConfigServiceRetentionConfigurationSuite struct {
	suite.Suite
	svc *configservice.Client
}

func (suite *TestConfigServiceRetentionConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = configservice.NewFromConfig(cfg)
}

func (suite *TestConfigServiceRetentionConfigurationSuite) TestList() {
	a := assert.New(suite.T())
	lister := ConfigServiceRetentionConfigurationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	_ = resources
}

func TestConfigServiceRetentionConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestConfigServiceRetentionConfigurationSuite))
}

//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
)

type TestAppRunnerObservabilityConfigurationSuite struct {
	suite.Suite
	svc *apprunner.Client
}

func (suite *TestAppRunnerObservabilityConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = apprunner.NewFromConfig(cfg)
}

func (suite *TestAppRunnerObservabilityConfigurationSuite) TestList() {
	a := assert.New(suite.T())
	lister := AppRunnerObservabilityConfigurationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAppRunnerListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAppRunnerObservabilityConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestAppRunnerObservabilityConfigurationSuite))
}

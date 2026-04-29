//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"
)

type TestResilienceHubAppSuite struct {
	suite.Suite
	svc *resiliencehub.Client
}

func (suite *TestResilienceHubAppSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = resiliencehub.NewFromConfig(cfg)
}

func (suite *TestResilienceHubAppSuite) TestList() {
	a := assert.New(suite.T())
	lister := ResilienceHubAppLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testResilienceHubListerOpts)
	a.NoError(err)
	_ = resources
}

func TestResilienceHubAppIntegration(t *testing.T) {
	suite.Run(t, new(TestResilienceHubAppSuite))
}

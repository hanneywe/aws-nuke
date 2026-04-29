//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
)

type TestAmplifyWebhookSuite struct {
	suite.Suite
	svc *amplify.Client
}

func (suite *TestAmplifyWebhookSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = amplify.NewFromConfig(cfg)
}

func (suite *TestAmplifyWebhookSuite) TestList() {
	a := assert.New(suite.T())
	lister := AmplifyWebhookLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAmplifyWebhookIntegration(t *testing.T) {
	suite.Run(t, new(TestAmplifyWebhookSuite))
}

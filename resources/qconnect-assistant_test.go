//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/qconnect"
)

type TestQConnectAssistantSuite struct {
	suite.Suite
	svc *qconnect.Client
}

func (suite *TestQConnectAssistantSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = qconnect.NewFromConfig(cfg)
}

func (suite *TestQConnectAssistantSuite) TestList() {
	a := assert.New(suite.T())
	lister := QConnectAssistantLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	_ = resources
}

func TestQConnectAssistantIntegration(t *testing.T) {
	suite.Run(t, new(TestQConnectAssistantSuite))
}

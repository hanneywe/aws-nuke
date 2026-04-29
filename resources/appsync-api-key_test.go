//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

type TestAppSyncAPIKeySuite struct {
	suite.Suite
	svc *appsync.Client
}

func (suite *TestAppSyncAPIKeySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = appsync.NewFromConfig(cfg)
}

func (suite *TestAppSyncAPIKeySuite) TestList() {
	a := assert.New(suite.T())
	lister := AppSyncAPIKeyLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAppSyncAPIKeyIntegration(t *testing.T) {
	suite.Run(t, new(TestAppSyncAPIKeySuite))
}

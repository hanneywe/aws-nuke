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

type TestAppSyncAPICacheSuite struct {
	suite.Suite
	svc *appsync.Client
}

func (suite *TestAppSyncAPICacheSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = appsync.NewFromConfig(cfg)
}

func (suite *TestAppSyncAPICacheSuite) TestList() {
	a := assert.New(suite.T())
	lister := AppSyncAPICacheLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAppSyncAPICacheIntegration(t *testing.T) {
	suite.Run(t, new(TestAppSyncAPICacheSuite))
}

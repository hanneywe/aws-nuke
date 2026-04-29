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

type TestAppSyncTypeSuite struct {
	suite.Suite
	svc *appsync.Client
}

func (suite *TestAppSyncTypeSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = appsync.NewFromConfig(cfg)
}

func (suite *TestAppSyncTypeSuite) TestList() {
	a := assert.New(suite.T())
	lister := AppSyncTypeLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAppSyncListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAppSyncTypeIntegration(t *testing.T) {
	suite.Run(t, new(TestAppSyncTypeSuite))
}

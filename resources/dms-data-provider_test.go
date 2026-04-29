//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
)

type TestDMSDataProviderSuite struct {
	suite.Suite
	svc *databasemigrationservice.Client
}

func (suite *TestDMSDataProviderSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = databasemigrationservice.NewFromConfig(cfg)
}

func (suite *TestDMSDataProviderSuite) TestList() {
	a := assert.New(suite.T())
	lister := DMSDataProviderLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDMSListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDMSDataProviderIntegration(t *testing.T) {
	suite.Run(t, new(TestDMSDataProviderSuite))
}

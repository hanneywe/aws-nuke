//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

type TestCloudTrailDashboardSuite struct {
	suite.Suite
	svc *cloudtrail.Client
}

func (suite *TestCloudTrailDashboardSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = cloudtrail.NewFromConfig(cfg)
}

func (suite *TestCloudTrailDashboardSuite) TestList() {
	a := assert.New(suite.T())
	lister := CloudTrailDashboardLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testCloudTrailListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCloudTrailDashboardIntegration(t *testing.T) {
	suite.Run(t, new(TestCloudTrailDashboardSuite))
}

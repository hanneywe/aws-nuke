//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestEC2NetworkInsightsAccessScopeSuite struct {
	suite.Suite
	svc *ec2.Client
}

func (suite *TestEC2NetworkInsightsAccessScopeSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = ec2.NewFromConfig(cfg)
}

func (suite *TestEC2NetworkInsightsAccessScopeSuite) TearDownSuite() {
	// No persistent resources to clean up
}

func (suite *TestEC2NetworkInsightsAccessScopeSuite) TestList() {
	assertions := assert.New(suite.T())

	lister := EC2NetworkInsightsAccessScopeLister{}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	assertions.NoError(err)

	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Config: &cfg,
	})

	assertions.NoError(err)
	assertions.NotNil(resources)
}

func TestEC2NetworkInsightsAccessScopeIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2NetworkInsightsAccessScopeSuite))
}

//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/outposts"
)

type TestOutpostsSiteSuite struct {
	suite.Suite
	svc *outposts.Client
}

func (suite *TestOutpostsSiteSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = outposts.NewFromConfig(cfg)
}

func (suite *TestOutpostsSiteSuite) TestList() {
	a := assert.New(suite.T())
	lister := OutpostsSiteLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testOutpostsListerOpts)
	a.NoError(err)
	_ = resources
}

func TestOutpostsSiteIntegration(t *testing.T) {
	suite.Run(t, new(TestOutpostsSiteSuite))
}

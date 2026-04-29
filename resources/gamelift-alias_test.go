//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

type TestGameLiftAliasSuite struct {
	suite.Suite
	svc *gamelift.Client
}

func (suite *TestGameLiftAliasSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = gamelift.NewFromConfig(cfg)
}

func (suite *TestGameLiftAliasSuite) TestList() {
	a := assert.New(suite.T())
	lister := GameLiftAliasLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGameLiftAliasIntegration(t *testing.T) {
	suite.Run(t, new(TestGameLiftAliasSuite))
}

//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codestarconnections"
)

type TestCodeStarConnectionsHostSuite struct {
	suite.Suite
	svc *codestarconnections.Client
}

func (suite *TestCodeStarConnectionsHostSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = codestarconnections.NewFromConfig(cfg)
}

func (suite *TestCodeStarConnectionsHostSuite) TestList() {
	a := assert.New(suite.T())
	lister := CodeStarConnectionsHostLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testCodeStarConnectionsListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCodeStarConnectionsHostIntegration(t *testing.T) {
	suite.Run(t, new(TestCodeStarConnectionsHostSuite))
}

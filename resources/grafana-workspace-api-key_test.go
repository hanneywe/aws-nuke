//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
)

type TestGrafanaWorkspaceAPIKeySuite struct {
	suite.Suite
	svc *grafana.Client
}

func (suite *TestGrafanaWorkspaceAPIKeySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = grafana.NewFromConfig(cfg)
}

func (suite *TestGrafanaWorkspaceAPIKeySuite) TestList() {
	a := assert.New(suite.T())
	lister := GrafanaWorkspaceAPIKeyLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestGrafanaWorkspaceAPIKeyIntegration(t *testing.T) {
	suite.Run(t, new(TestGrafanaWorkspaceAPIKeySuite))
}

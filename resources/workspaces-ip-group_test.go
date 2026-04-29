//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
)

type TestWorkSpacesIpGroupSuite struct {
	suite.Suite
	svc *workspaces.Client
}

func (suite *TestWorkSpacesIpGroupSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = workspaces.NewFromConfig(cfg)
}

func (suite *TestWorkSpacesIpGroupSuite) TestList() {
	a := assert.New(suite.T())
	lister := WorkSpacesIpGroupLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testWorkSpacesV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestWorkSpacesIpGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestWorkSpacesIpGroupSuite))
}

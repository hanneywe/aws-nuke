//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"
)

type TestWorkSpacesWebPortalSuite struct {
	suite.Suite
	svc *workspacesweb.Client
}

func (s *TestWorkSpacesWebPortalSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = workspacesweb.NewFromConfig(cfg)
}

func (s *TestWorkSpacesWebPortalSuite) TestList() {
	assertions := assert.New(s.T())

	lister := WorkSpacesWebPortalLister{}
	resources, err := lister.List(context.TODO(), testWorkSpacesWebListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestWorkSpacesWebPortalIntegration(t *testing.T) {
	suite.Run(t, new(TestWorkSpacesWebPortalSuite))
}

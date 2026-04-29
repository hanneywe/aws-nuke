//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/omics"
)

type TestOmicsWorkflowSuite struct {
	suite.Suite
	svc *omics.Client
}

func (s *TestOmicsWorkflowSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = omics.NewFromConfig(cfg)
}

func (s *TestOmicsWorkflowSuite) TestList() {
	assertions := assert.New(s.T())

	lister := OmicsWorkflowLister{}
	resources, err := lister.List(context.TODO(), testOmicsListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestOmicsWorkflowIntegration(t *testing.T) {
	suite.Run(t, new(TestOmicsWorkflowSuite))
}

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

type TestOmicsRunGroupSuite struct {
	suite.Suite
	svc *omics.Client
}

func (s *TestOmicsRunGroupSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = omics.NewFromConfig(cfg)
}

func (s *TestOmicsRunGroupSuite) TestList() {
	assertions := assert.New(s.T())

	lister := OmicsRunGroupLister{}
	resources, err := lister.List(context.TODO(), testOmicsListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestOmicsRunGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestOmicsRunGroupSuite))
}

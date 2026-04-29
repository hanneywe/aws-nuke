//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

type TestNeptuneDBClusterEndpointSuite struct {
	suite.Suite
	svc *neptune.Client
}

func (s *TestNeptuneDBClusterEndpointSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = neptune.NewFromConfig(cfg)
}

func (s *TestNeptuneDBClusterEndpointSuite) TestList() {
	a := assert.New(s.T())
	lister := NeptuneDBClusterEndpointLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestNeptuneDBClusterEndpointIntegration(t *testing.T) {
	suite.Run(t, new(TestNeptuneDBClusterEndpointSuite))
}

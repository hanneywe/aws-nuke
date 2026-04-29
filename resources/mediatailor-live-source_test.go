//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
)

type TestMediaTailorLiveSourceSuite struct {
	suite.Suite
	svc *mediatailor.Client
}

func (s *TestMediaTailorLiveSourceSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = mediatailor.NewFromConfig(cfg)
}

func (s *TestMediaTailorLiveSourceSuite) TestList() {
	a := assert.New(s.T())
	lister := MediaTailorLiveSourceLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testMediaTailorV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMediaTailorLiveSourceIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaTailorLiveSourceSuite))
}

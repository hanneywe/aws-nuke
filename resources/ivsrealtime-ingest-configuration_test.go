//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"
)

type TestIVSRealtimeIngestConfigurationSuite struct {
	suite.Suite
	svc *ivsrealtime.Client
}

func (s *TestIVSRealtimeIngestConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = ivsrealtime.NewFromConfig(cfg)
}

func (s *TestIVSRealtimeIngestConfigurationSuite) TestList() {
	assertions := assert.New(s.T())
	lister := IVSRealtimeIngestConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	assertions.NoError(err)
	_ = resources
}

func TestIVSRealtimeIngestConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSRealtimeIngestConfigurationSuite))
}

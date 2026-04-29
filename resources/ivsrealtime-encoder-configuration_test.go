//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"
)

type TestIVSRealtimeEncoderConfigurationSuite struct {
	suite.Suite
	svc *ivsrealtime.Client
	arn *string
}

func (s *TestIVSRealtimeEncoderConfigurationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = ivsrealtime.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateEncoderConfiguration(ctx, &ivsrealtime.CreateEncoderConfigurationInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create encoder configuration: %v", err)
	}
	s.arn = resp.EncoderConfiguration.Arn
}

func (s *TestIVSRealtimeEncoderConfigurationSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteEncoderConfiguration(ctx, &ivsrealtime.DeleteEncoderConfigurationInput{
		Arn: s.arn,
	})
}

func (s *TestIVSRealtimeEncoderConfigurationSuite) TestList() {
	a := assert.New(s.T())
	lister := &IVSRealtimeEncoderConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIVSRealtimeEncoderConfigurationSuite) TestRemove() {
	a := assert.New(s.T())
	ec := &IVSRealtimeEncoderConfiguration{svc: s.svc, ARN: s.arn}
	a.NoError(ec.Remove(context.TODO()))
}

func TestIVSRealtimeEncoderConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSRealtimeEncoderConfigurationSuite))
}

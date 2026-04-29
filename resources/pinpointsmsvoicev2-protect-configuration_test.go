//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
)

type TestPinpointSMSVoiceV2ProtectConfigurationSuite struct {
	suite.Suite
	svc *pinpointsmsvoicev2.Client
	id  *string
}

func (s *TestPinpointSMSVoiceV2ProtectConfigurationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = pinpointsmsvoicev2.NewFromConfig(cfg)

	resp, err := s.svc.CreateProtectConfiguration(ctx, &pinpointsmsvoicev2.CreateProtectConfigurationInput{})
	if err != nil {
		s.T().Fatalf("failed to create protect configuration: %v", err)
	}
	s.id = resp.ProtectConfigurationId
}

func (s *TestPinpointSMSVoiceV2ProtectConfigurationSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteProtectConfiguration(ctx, &pinpointsmsvoicev2.DeleteProtectConfigurationInput{
		ProtectConfigurationId: s.id,
	})
}

func (s *TestPinpointSMSVoiceV2ProtectConfigurationSuite) TestList() {
	a := assert.New(s.T())
	lister := &PinpointSMSVoiceV2ProtectConfigurationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestPinpointSMSVoiceV2ProtectConfigurationSuite) TestRemove() {
	a := assert.New(s.T())
	pc := &PinpointSMSVoiceV2ProtectConfiguration{svc: s.svc, ProtectConfigurationID: s.id}
	a.NoError(pc.Remove(context.TODO()))
}

func TestPinpointSMSVoiceV2ProtectConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestPinpointSMSVoiceV2ProtectConfigurationSuite))
}

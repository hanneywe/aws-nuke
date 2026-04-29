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
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
)

type TestPinpointSMSVoiceV2ConfigurationSetSuite struct {
	suite.Suite
	svc  *pinpointsmsvoicev2.Client
	name *string
}

func (s *TestPinpointSMSVoiceV2ConfigurationSetSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = pinpointsmsvoicev2.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateConfigurationSet(ctx, &pinpointsmsvoicev2.CreateConfigurationSetInput{
		ConfigurationSetName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create configuration set: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestPinpointSMSVoiceV2ConfigurationSetSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteConfigurationSet(ctx, &pinpointsmsvoicev2.DeleteConfigurationSetInput{
		ConfigurationSetName: s.name,
	})
}

func (s *TestPinpointSMSVoiceV2ConfigurationSetSuite) TestList() {
	a := assert.New(s.T())
	lister := &PinpointSMSVoiceV2ConfigurationSetLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestPinpointSMSVoiceV2ConfigurationSetSuite) TestRemove() {
	a := assert.New(s.T())
	cs := &PinpointSMSVoiceV2ConfigurationSet{svc: s.svc, ConfigurationSetName: s.name}
	a.NoError(cs.Remove(context.TODO()))
}

func TestPinpointSMSVoiceV2ConfigurationSetIntegration(t *testing.T) {
	suite.Run(t, new(TestPinpointSMSVoiceV2ConfigurationSetSuite))
}

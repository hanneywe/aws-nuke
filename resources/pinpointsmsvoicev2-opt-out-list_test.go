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

type TestPinpointSMSVoiceV2OptOutListSuite struct {
	suite.Suite
	svc  *pinpointsmsvoicev2.Client
	name *string
}

func (s *TestPinpointSMSVoiceV2OptOutListSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = pinpointsmsvoicev2.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateOptOutList(ctx, &pinpointsmsvoicev2.CreateOptOutListInput{
		OptOutListName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create opt-out list: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestPinpointSMSVoiceV2OptOutListSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteOptOutList(ctx, &pinpointsmsvoicev2.DeleteOptOutListInput{
		OptOutListName: s.name,
	})
}

func (s *TestPinpointSMSVoiceV2OptOutListSuite) TestList() {
	a := assert.New(s.T())
	lister := &PinpointSMSVoiceV2OptOutListLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestPinpointSMSVoiceV2OptOutListSuite) TestRemove() {
	a := assert.New(s.T())
	ol := &PinpointSMSVoiceV2OptOutList{svc: s.svc, OptOutListName: s.name}
	a.NoError(ol.Remove(context.TODO()))
}

func TestPinpointSMSVoiceV2OptOutListIntegration(t *testing.T) {
	suite.Run(t, new(TestPinpointSMSVoiceV2OptOutListSuite))
}

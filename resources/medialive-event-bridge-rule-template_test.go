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
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	mltypes "github.com/aws/aws-sdk-go-v2/service/medialive/types"
)

type TestMediaLiveEventBridgeRuleTemplateSuite struct {
	suite.Suite
	svc     *medialive.Client
	id      *string
	groupID *string
}

func (s *TestMediaLiveEventBridgeRuleTemplateSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = medialive.NewFromConfig(cfg)

	groupName := fmt.Sprintf("aws-nuke-test-group-%d", time.Now().UnixNano())
	groupResp, err := s.svc.CreateEventBridgeRuleTemplateGroup(ctx, &medialive.CreateEventBridgeRuleTemplateGroupInput{
		Name: ptr.String(groupName),
	})
	if err != nil {
		s.T().Fatalf("failed to create event bridge rule template group: %v", err)
	}
	s.groupID = groupResp.Id

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateEventBridgeRuleTemplate(ctx, &medialive.CreateEventBridgeRuleTemplateInput{
		Name:            ptr.String(name),
		GroupIdentifier: s.groupID,
		EventType:       mltypes.EventBridgeRuleTemplateEventTypeMedialiveMultiplexAlert,
	})
	if err != nil {
		s.T().Fatalf("failed to create event bridge rule template: %v", err)
	}
	s.id = resp.Id
}

func (s *TestMediaLiveEventBridgeRuleTemplateSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteEventBridgeRuleTemplate(ctx, &medialive.DeleteEventBridgeRuleTemplateInput{
		Identifier: s.id,
	})
	_, _ = s.svc.DeleteEventBridgeRuleTemplateGroup(ctx, &medialive.DeleteEventBridgeRuleTemplateGroupInput{
		Identifier: s.groupID,
	})
}

func (s *TestMediaLiveEventBridgeRuleTemplateSuite) TestList() {
	a := assert.New(s.T())
	lister := &MediaLiveEventBridgeRuleTemplateLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestMediaLiveEventBridgeRuleTemplateSuite) TestRemove() {
	a := assert.New(s.T())
	r := &MediaLiveEventBridgeRuleTemplate{svc: s.svc, ID: s.id}
	a.NoError(r.Remove(context.TODO()))
}

func TestMediaLiveEventBridgeRuleTemplateIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaLiveEventBridgeRuleTemplateSuite))
}

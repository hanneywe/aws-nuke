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
	"github.com/aws/aws-sdk-go-v2/service/ivschat"
)

type TestIVSChatRoomSuite struct {
	suite.Suite
	svc *ivschat.Client
	arn *string
}

func (s *TestIVSChatRoomSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = ivschat.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateRoom(ctx, &ivschat.CreateRoomInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create room: %v", err)
	}
	s.arn = resp.Arn
}

func (s *TestIVSChatRoomSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteRoom(ctx, &ivschat.DeleteRoomInput{
		Identifier: s.arn,
	})
}

func (s *TestIVSChatRoomSuite) TestList() {
	a := assert.New(s.T())
	lister := &IVSChatRoomLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIVSChatRoomSuite) TestRemove() {
	a := assert.New(s.T())
	room := &IVSChatRoom{svc: s.svc, ARN: s.arn}
	a.NoError(room.Remove(context.TODO()))
}

func TestIVSChatRoomIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSChatRoomSuite))
}

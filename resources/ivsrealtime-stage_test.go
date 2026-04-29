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

type TestIVSRealtimeStageSuite struct {
	suite.Suite
	svc *ivsrealtime.Client
	arn *string
}

func (s *TestIVSRealtimeStageSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = ivsrealtime.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateStage(ctx, &ivsrealtime.CreateStageInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create stage: %v", err)
	}
	s.arn = resp.Stage.Arn
}

func (s *TestIVSRealtimeStageSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteStage(ctx, &ivsrealtime.DeleteStageInput{
		Arn: s.arn,
	})
}

func (s *TestIVSRealtimeStageSuite) TestList() {
	a := assert.New(s.T())
	lister := &IVSRealtimeStageLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIVSRealtimeStageSuite) TestRemove() {
	a := assert.New(s.T())
	stage := &IVSRealtimeStage{svc: s.svc, ARN: s.arn}
	a.NoError(stage.Remove(context.TODO()))
}

func TestIVSRealtimeStageIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSRealtimeStageSuite))
}

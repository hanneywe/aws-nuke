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
	"github.com/aws/aws-sdk-go-v2/service/ivs"
)

type TestIVSChannelSuite struct {
	suite.Suite
	svc *ivs.Client
	arn *string
}

func (s *TestIVSChannelSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = ivs.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateChannel(ctx, &ivs.CreateChannelInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create channel: %v", err)
	}
	s.arn = resp.Channel.Arn
}

func (s *TestIVSChannelSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteChannel(ctx, &ivs.DeleteChannelInput{
		Arn: s.arn,
	})
}

func (s *TestIVSChannelSuite) TestList() {
	a := assert.New(s.T())
	lister := &IVSChannelLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIVSListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestIVSChannelSuite) TestRemove() {
	a := assert.New(s.T())
	ch := &IVSChannel{svc: s.svc, ARN: s.arn}
	a.NoError(ch.Remove(context.TODO()))
}

func TestIVSChannelIntegration(t *testing.T) {
	suite.Run(t, new(TestIVSChannelSuite))
}

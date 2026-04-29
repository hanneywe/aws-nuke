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
)

type TestMediaLiveNetworkSuite struct {
	suite.Suite
	svc *medialive.Client
	id  *string
}

func (s *TestMediaLiveNetworkSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = medialive.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateNetwork(ctx, &medialive.CreateNetworkInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create network: %v", err)
	}
	s.id = resp.Id
}

func (s *TestMediaLiveNetworkSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteNetwork(ctx, &medialive.DeleteNetworkInput{
		NetworkId: s.id,
	})
}

func (s *TestMediaLiveNetworkSuite) TestList() {
	a := assert.New(s.T())
	lister := &MediaLiveNetworkLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestMediaLiveNetworkSuite) TestRemove() {
	a := assert.New(s.T())
	r := &MediaLiveNetwork{svc: s.svc, ID: s.id}
	a.NoError(r.Remove(context.TODO()))
}

func TestMediaLiveNetworkIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaLiveNetworkSuite))
}

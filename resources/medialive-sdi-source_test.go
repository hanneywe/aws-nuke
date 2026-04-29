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

type TestMediaLiveSdiSourceSuite struct {
	suite.Suite
	svc *medialive.Client
	id  *string
}

func (s *TestMediaLiveSdiSourceSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = medialive.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateSdiSource(ctx, &medialive.CreateSdiSourceInput{
		Name: ptr.String(name),
		Type: mltypes.SdiSourceTypeSingle,
	})
	if err != nil {
		s.T().Fatalf("failed to create sdi source: %v", err)
	}
	s.id = resp.SdiSource.Id
}

func (s *TestMediaLiveSdiSourceSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteSdiSource(ctx, &medialive.DeleteSdiSourceInput{
		SdiSourceId: s.id,
	})
}

func (s *TestMediaLiveSdiSourceSuite) TestList() {
	a := assert.New(s.T())
	lister := &MediaLiveSdiSourceLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestMediaLiveSdiSourceSuite) TestRemove() {
	a := assert.New(s.T())
	r := &MediaLiveSdiSource{svc: s.svc, ID: s.id}
	a.NoError(r.Remove(context.TODO()))
}

func TestMediaLiveSdiSourceIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaLiveSdiSourceSuite))
}

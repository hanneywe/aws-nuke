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
	"github.com/aws/aws-sdk-go-v2/service/location"
	locationtypes "github.com/aws/aws-sdk-go-v2/service/location/types"
)

type TestLocationServiceMapSuite struct {
	suite.Suite
	svc  *location.Client
	name *string
}

func (s *TestLocationServiceMapSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = location.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateMap(ctx, &location.CreateMapInput{
		MapName:       ptr.String(name),
		Configuration: &locationtypes.MapConfiguration{Style: ptr.String("VectorEsriStreets")},
	})
	if err != nil {
		s.T().Fatalf("failed to create map: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestLocationServiceMapSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteMap(ctx, &location.DeleteMapInput{MapName: s.name})
}

func (s *TestLocationServiceMapSuite) TestList() {
	a := assert.New(s.T())
	lister := &LocationServiceMapLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestLocationServiceMapSuite) TestRemove() {
	a := assert.New(s.T())
	m := &LocationServiceMap{svc: s.svc, MapName: s.name}
	a.NoError(m.Remove(context.TODO()))
}

func TestLocationServiceMapIntegration(t *testing.T) {
	suite.Run(t, new(TestLocationServiceMapSuite))
}

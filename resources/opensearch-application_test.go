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
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

type TestOpenSearchApplicationSuite struct {
	suite.Suite
	svc   *opensearch.Client
	appID *string
}

func (s *TestOpenSearchApplicationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = opensearch.NewFromConfig(cfg)

	appName := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateApplication(ctx, &opensearch.CreateApplicationInput{
		Name: ptr.String(appName),
	})
	if err != nil {
		s.T().Fatalf("failed to create application: %v", err)
	}
	s.appID = resp.Id
}

func (s *TestOpenSearchApplicationSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.appID != nil {
		_, _ = s.svc.DeleteApplication(ctx, &opensearch.DeleteApplicationInput{
			Id: s.appID,
		})
	}
}

func (s *TestOpenSearchApplicationSuite) TestList() {
	a := assert.New(s.T())
	lister := &OpenSearchApplicationLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testOpenSearchListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestOpenSearchApplicationSuite) TestRemove() {
	a := assert.New(s.T())
	app := &OpenSearchApplication{svc: s.svc, Id: s.appID}
	a.NoError(app.Remove(context.TODO()))
}

func TestOpenSearchApplicationIntegration(t *testing.T) {
	suite.Run(t, new(TestOpenSearchApplicationSuite))
}

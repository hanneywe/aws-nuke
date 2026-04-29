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
	"github.com/aws/aws-sdk-go-v2/service/kendraranking"
)

type TestKendraRankingRescoreExecutionPlanSuite struct {
	suite.Suite
	svc *kendraranking.Client
	id  *string
}

func (s *TestKendraRankingRescoreExecutionPlanSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = kendraranking.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateRescoreExecutionPlan(ctx, &kendraranking.CreateRescoreExecutionPlanInput{
		Name: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create rescore execution plan: %v", err)
	}
	s.id = resp.Id
}

func (s *TestKendraRankingRescoreExecutionPlanSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteRescoreExecutionPlan(ctx, &kendraranking.DeleteRescoreExecutionPlanInput{
		Id: s.id,
	})
}

func (s *TestKendraRankingRescoreExecutionPlanSuite) TestList() {
	a := assert.New(s.T())
	lister := &KendraRankingRescoreExecutionPlanLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testKendraRankingListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestKendraRankingRescoreExecutionPlanSuite) TestRemove() {
	a := assert.New(s.T())
	plan := &KendraRankingRescoreExecutionPlan{svc: s.svc, Id: s.id}
	a.NoError(plan.Remove(context.TODO()))
}

func TestKendraRankingRescoreExecutionPlanIntegration(t *testing.T) {
	suite.Run(t, new(TestKendraRankingRescoreExecutionPlanSuite))
}

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
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

type TestSageMakerExperimentSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerExperimentSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateExperiment(ctx, &sagemaker.CreateExperimentInput{
		ExperimentName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create experiment: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerExperimentSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteExperiment(ctx, &sagemaker.DeleteExperimentInput{
		ExperimentName: s.name,
	})
}

func (s *TestSageMakerExperimentSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerExperimentLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerExperimentSuite) TestRemove() {
	a := assert.New(s.T())
	experiment := &SageMakerExperiment{svc: s.svc, ExperimentName: s.name}
	a.NoError(experiment.Remove(context.TODO()))
}

func TestSageMakerExperimentIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerExperimentSuite))
}

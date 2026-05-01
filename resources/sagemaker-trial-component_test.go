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

type TestSageMakerTrialComponentSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerTrialComponentSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateTrialComponent(ctx, &sagemaker.CreateTrialComponentInput{
		TrialComponentName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create trial component: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerTrialComponentSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteTrialComponent(ctx, &sagemaker.DeleteTrialComponentInput{
		TrialComponentName: s.name,
	})
}

func (s *TestSageMakerTrialComponentSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerTrialComponentLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerTrialComponentSuite) TestRemove() {
	a := assert.New(s.T())
	component := &SageMakerTrialComponent{svc: s.svc, TrialComponentName: s.name}
	a.NoError(component.Remove(context.TODO()))
}

func TestSageMakerTrialComponentIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerTrialComponentSuite))
}

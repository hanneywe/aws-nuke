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

type TestSageMakerModelCardSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerModelCardSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateModelCard(ctx, &sagemaker.CreateModelCardInput{
		ModelCardName:   ptr.String(name),
		Content:         ptr.String(`{"model_overview":{"model_description":"test"}}`),
		ModelCardStatus: "Draft",
	})
	if err != nil {
		s.T().Fatalf("failed to create model card: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerModelCardSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteModelCard(ctx, &sagemaker.DeleteModelCardInput{
		ModelCardName: s.name,
	})
}

func (s *TestSageMakerModelCardSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerModelCardLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerModelCardSuite) TestRemove() {
	a := assert.New(s.T())
	card := &SageMakerModelCard{svc: s.svc, ModelCardName: s.name}
	a.NoError(card.Remove(context.TODO()))
}

func TestSageMakerModelCardIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerModelCardSuite))
}

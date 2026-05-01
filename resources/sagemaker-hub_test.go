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

type TestSageMakerHubSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerHubSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateHub(ctx, &sagemaker.CreateHubInput{
		HubName:        ptr.String(name),
		HubDescription: ptr.String("test hub for aws-nuke"),
	})
	if err != nil {
		s.T().Fatalf("failed to create hub: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerHubSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteHub(ctx, &sagemaker.DeleteHubInput{
		HubName: s.name,
	})
}

func (s *TestSageMakerHubSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerHubLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerHubSuite) TestRemove() {
	a := assert.New(s.T())
	hub := &SageMakerHub{svc: s.svc, HubName: s.name}
	a.NoError(hub.Remove(context.TODO()))
}

func TestSageMakerHubIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerHubSuite))
}

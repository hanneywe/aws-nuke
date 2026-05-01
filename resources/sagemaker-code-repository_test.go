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
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

type TestSageMakerCodeRepositorySuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerCodeRepositorySuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateCodeRepository(ctx, &sagemaker.CreateCodeRepositoryInput{
		CodeRepositoryName: ptr.String(name),
		GitConfig: &sagemakertypes.GitConfig{
			RepositoryUrl: ptr.String("https://github.com/example/repo.git"),
		},
	})
	if err != nil {
		s.T().Fatalf("failed to create code repository: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerCodeRepositorySuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteCodeRepository(ctx, &sagemaker.DeleteCodeRepositoryInput{
		CodeRepositoryName: s.name,
	})
}

func (s *TestSageMakerCodeRepositorySuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerCodeRepositoryLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerCodeRepositorySuite) TestRemove() {
	a := assert.New(s.T())
	repo := &SageMakerCodeRepository{svc: s.svc, CodeRepositoryName: s.name}
	a.NoError(repo.Remove(context.TODO()))
}

func TestSageMakerCodeRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerCodeRepositorySuite))
}

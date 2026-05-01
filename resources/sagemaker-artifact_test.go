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

type TestSageMakerArtifactSuite struct {
	suite.Suite
	svc *sagemaker.Client
	arn *string
}

func (s *TestSageMakerArtifactSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := s.svc.CreateArtifact(ctx, &sagemaker.CreateArtifactInput{
		ArtifactName: ptr.String(name),
		ArtifactType: ptr.String("test"),
		Source: &sagemakertypes.ArtifactSource{
			SourceUri: ptr.String(fmt.Sprintf("s3://test-bucket/%s", name)),
		},
	})
	if err != nil {
		s.T().Fatalf("failed to create artifact: %v", err)
	}
	s.arn = resp.ArtifactArn
}

func (s *TestSageMakerArtifactSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteArtifact(ctx, &sagemaker.DeleteArtifactInput{
		ArtifactArn: s.arn,
	})
}

func (s *TestSageMakerArtifactSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerArtifactLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerArtifactSuite) TestRemove() {
	a := assert.New(s.T())
	artifact := &SageMakerArtifact{svc: s.svc, ArtifactArn: s.arn}
	a.NoError(artifact.Remove(context.TODO()))
}

func TestSageMakerArtifactIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerArtifactSuite))
}

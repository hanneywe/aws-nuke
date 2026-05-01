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

type TestSageMakerModelPackageGroupSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerModelPackageGroupSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateModelPackageGroup(ctx, &sagemaker.CreateModelPackageGroupInput{
		ModelPackageGroupName: ptr.String(name),
	})
	if err != nil {
		s.T().Fatalf("failed to create model package group: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerModelPackageGroupSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteModelPackageGroup(ctx, &sagemaker.DeleteModelPackageGroupInput{
		ModelPackageGroupName: s.name,
	})
}

func (s *TestSageMakerModelPackageGroupSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerModelPackageGroupLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerModelPackageGroupSuite) TestRemove() {
	a := assert.New(s.T())
	group := &SageMakerModelPackageGroup{svc: s.svc, ModelPackageGroupName: s.name}
	a.NoError(group.Remove(context.TODO()))
}

func TestSageMakerModelPackageGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerModelPackageGroupSuite))
}

//go:build integration

package resources

import (
	"context"
	"encoding/base64"
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

type TestSageMakerStudioLifecycleConfigSuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerStudioLifecycleConfigSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateStudioLifecycleConfig(ctx, &sagemaker.CreateStudioLifecycleConfigInput{
		StudioLifecycleConfigName:    ptr.String(name),
		StudioLifecycleConfigAppType: sagemakertypes.StudioLifecycleConfigAppTypeJupyterServer,
		StudioLifecycleConfigContent: ptr.String(base64.StdEncoding.EncodeToString([]byte("echo test"))),
	})
	if err != nil {
		s.T().Fatalf("failed to create studio lifecycle config: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerStudioLifecycleConfigSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteStudioLifecycleConfig(ctx, &sagemaker.DeleteStudioLifecycleConfigInput{
		StudioLifecycleConfigName: s.name,
	})
}

func (s *TestSageMakerStudioLifecycleConfigSuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerStudioLifecycleConfigLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerStudioLifecycleConfigSuite) TestRemove() {
	a := assert.New(s.T())
	lc := &SageMakerStudioLifecycleConfig{svc: s.svc, StudioLifecycleConfigName: s.name}
	a.NoError(lc.Remove(context.TODO()))
}

func TestSageMakerStudioLifecycleConfigIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerStudioLifecycleConfigSuite))
}

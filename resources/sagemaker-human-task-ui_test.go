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

type TestSageMakerHumanTaskUISuite struct {
	suite.Suite
	svc  *sagemaker.Client
	name *string
}

func (s *TestSageMakerHumanTaskUISuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = sagemaker.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	_, err = s.svc.CreateHumanTaskUi(ctx, &sagemaker.CreateHumanTaskUiInput{
		HumanTaskUiName: ptr.String(name),
		UiTemplate: &sagemakertypes.UiTemplate{
			Content: ptr.String("<html><body>test</body></html>"),
		},
	})
	if err != nil {
		s.T().Fatalf("failed to create human task UI: %v", err)
	}
	s.name = ptr.String(name)
}

func (s *TestSageMakerHumanTaskUISuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteHumanTaskUi(ctx, &sagemaker.DeleteHumanTaskUiInput{
		HumanTaskUiName: s.name,
	})
}

func (s *TestSageMakerHumanTaskUISuite) TestList() {
	a := assert.New(s.T())
	lister := &SageMakerHumanTaskUILister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestSageMakerHumanTaskUISuite) TestRemove() {
	a := assert.New(s.T())
	ui := &SageMakerHumanTaskUI{svc: s.svc, HumanTaskUIName: s.name}
	a.NoError(ui.Remove(context.TODO()))
}

func TestSageMakerHumanTaskUIIntegration(t *testing.T) {
	suite.Run(t, new(TestSageMakerHumanTaskUISuite))
}

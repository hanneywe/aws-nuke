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
	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
)

type TestDeviceFarmTestGridProjectSuite struct {
	suite.Suite
	svc        *devicefarm.Client
	projectArn *string
}

func (s *TestDeviceFarmTestGridProjectSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = devicefarm.NewFromConfig(cfg)

	projectName := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	createOutput, err := s.svc.CreateTestGridProject(ctx, &devicefarm.CreateTestGridProjectInput{
		Name: ptr.String(projectName),
	})
	if err != nil {
		s.T().Fatalf("failed to create test grid project: %v", err)
	}

	s.projectArn = createOutput.TestGridProject.Arn
}

func (s *TestDeviceFarmTestGridProjectSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteTestGridProject(ctx, &devicefarm.DeleteTestGridProjectInput{
		ProjectArn: s.projectArn,
	})
}

func (s *TestDeviceFarmTestGridProjectSuite) TestList() {
	assertions := assert.New(s.T())

	lister := DeviceFarmTestGridProjectLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)

	assertions.NoError(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestDeviceFarmTestGridProjectSuite) TestRemove() {
	assertions := assert.New(s.T())

	testGridProject := DeviceFarmTestGridProject{
		svc:  s.svc,
		Arn:  s.projectArn,
		Name: ptr.String("test-grid-project"),
	}

	err := testGridProject.Remove(context.TODO())
	assertions.NoError(err)
}

func TestDeviceFarmTestGridProjectIntegration(t *testing.T) {
	suite.Run(t, new(TestDeviceFarmTestGridProjectSuite))
}

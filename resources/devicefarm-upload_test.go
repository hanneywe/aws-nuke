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

type TestDeviceFarmUploadSuite struct {
	suite.Suite
	svc        *devicefarm.Client
	projectArn *string
	uploadArn  *string
}

func (s *TestDeviceFarmUploadSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = devicefarm.NewFromConfig(cfg)

	projectName := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	projectOutput, err := s.svc.CreateProject(ctx, &devicefarm.CreateProjectInput{
		Name: ptr.String(projectName),
	})
	if err != nil {
		s.T().Fatalf("failed to create project: %v", err)
	}

	s.projectArn = projectOutput.Project.Arn

	uploadOutput, err := s.svc.CreateUpload(ctx, &devicefarm.CreateUploadInput{
		Name:       ptr.String("test-upload.apk"),
		ProjectArn: s.projectArn,
		Type:       "ANDROID_APP",
	})
	if err != nil {
		s.T().Fatalf("failed to create upload: %v", err)
	}

	s.uploadArn = uploadOutput.Upload.Arn
}

func (s *TestDeviceFarmUploadSuite) TearDownSuite() {
	ctx := context.TODO()
	if s.uploadArn != nil {
		_, _ = s.svc.DeleteUpload(ctx, &devicefarm.DeleteUploadInput{
			Arn: s.uploadArn,
		})
	}
	if s.projectArn != nil {
		_, _ = s.svc.DeleteProject(ctx, &devicefarm.DeleteProjectInput{
			Arn: s.projectArn,
		})
	}
}

func (s *TestDeviceFarmUploadSuite) TestList() {
	assertions := assert.New(s.T())

	lister := DeviceFarmUploadLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)

	assertions.NoError(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestDeviceFarmUploadSuite) TestRemove() {
	assertions := assert.New(s.T())

	upload := DeviceFarmUpload{
		svc:  s.svc,
		Arn:  s.uploadArn,
		Name: ptr.String("test-upload.apk"),
	}

	err := upload.Remove(context.TODO())
	assertions.NoError(err)
}

func TestDeviceFarmUploadIntegration(t *testing.T) {
	suite.Run(t, new(TestDeviceFarmUploadSuite))
}

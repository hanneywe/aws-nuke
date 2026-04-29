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

type TestDeviceFarmInstanceProfileSuite struct {
	suite.Suite
	svc                *devicefarm.Client
	instanceProfileArn *string
}

func (s *TestDeviceFarmInstanceProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = devicefarm.NewFromConfig(cfg)

	profileName := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	createOutput, err := s.svc.CreateInstanceProfile(ctx, &devicefarm.CreateInstanceProfileInput{
		Name: ptr.String(profileName),
	})
	if err != nil {
		s.T().Fatalf("failed to create instance profile: %v", err)
	}

	s.instanceProfileArn = createOutput.InstanceProfile.Arn
}

func (s *TestDeviceFarmInstanceProfileSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteInstanceProfile(ctx, &devicefarm.DeleteInstanceProfileInput{
		Arn: s.instanceProfileArn,
	})
}

func (s *TestDeviceFarmInstanceProfileSuite) TestList() {
	assertions := assert.New(s.T())

	lister := DeviceFarmInstanceProfileLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)

	assertions.NoError(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestDeviceFarmInstanceProfileSuite) TestRemove() {
	assertions := assert.New(s.T())

	instanceProfile := DeviceFarmInstanceProfile{
		svc:  s.svc,
		Arn:  s.instanceProfileArn,
		Name: ptr.String("test-instance-profile"),
	}

	err := instanceProfile.Remove(context.TODO())
	assertions.NoError(err)
}

func TestDeviceFarmInstanceProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestDeviceFarmInstanceProfileSuite))
}

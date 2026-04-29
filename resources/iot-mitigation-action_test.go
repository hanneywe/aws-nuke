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
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

type TestIoTMitigationActionSuite struct {
	suite.Suite
	svc        *iot.Client
	actionName *string
}

func (s *TestIoTMitigationActionSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = iot.NewFromConfig(cfg)

	s.actionName = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateMitigationAction(ctx, &iot.CreateMitigationActionInput{
		ActionName: s.actionName,
		RoleArn:    ptr.String("arn:aws:iam::123456789012:role/test-role"),
		ActionParams: &iottypes.MitigationActionParams{
			EnableIoTLoggingParams: &iottypes.EnableIoTLoggingParams{
				RoleArnForLogging: ptr.String("arn:aws:iam::123456789012:role/test-role"),
				LogLevel:          iottypes.LogLevelError,
			},
		},
	})
	if err != nil {
		s.T().Fatalf("failed to create mitigation action: %v", err)
	}
}

func (s *TestIoTMitigationActionSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = s.svc.DeleteMitigationAction(ctx, &iot.DeleteMitigationActionInput{
		ActionName: s.actionName,
	})
}

func (s *TestIoTMitigationActionSuite) TestList() {
	assertions := assert.New(s.T())
	lister := &IoTMitigationActionLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Greater(len(resources), 0)
}

func (s *TestIoTMitigationActionSuite) TestRemove() {
	assertions := assert.New(s.T())
	action := &IoTMitigationAction{svc: s.svc, ActionName: s.actionName}
	assertions.NoError(action.Remove(context.TODO()))
}

func TestIoTMitigationActionIntegration(t *testing.T) {
	suite.Run(t, new(TestIoTMitigationActionSuite))
}

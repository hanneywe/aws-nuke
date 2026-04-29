//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type TestCloudWatchLogsRetentionPolicySuite struct {
	suite.Suite
	svc *cloudwatchlogs.Client
}

func (s *TestCloudWatchLogsRetentionPolicySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = cloudwatchlogs.NewFromConfig(cfg)
}

func (s *TestCloudWatchLogsRetentionPolicySuite) TestList() {
	a := assert.New(s.T())
	lister := CloudWatchLogsRetentionPolicyLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testCloudWatchLogsV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCloudWatchLogsRetentionPolicyIntegration(t *testing.T) {
	suite.Run(t, new(TestCloudWatchLogsRetentionPolicySuite))
}

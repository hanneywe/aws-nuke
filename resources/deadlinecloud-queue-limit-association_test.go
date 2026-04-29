//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
)

type TestDeadlineCloudQueueLimitAssociationSuite struct {
	suite.Suite
	svc *deadline.Client
}

func (suite *TestDeadlineCloudQueueLimitAssociationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = deadline.NewFromConfig(cfg)
}

func (suite *TestDeadlineCloudQueueLimitAssociationSuite) TearDownSuite() {
	// Cleanup handled by farm/queue/limit teardown
}

func (suite *TestDeadlineCloudQueueLimitAssociationSuite) TestList() {
	a := assert.New(suite.T())

	lister := DeadlineCloudQueueLimitAssociationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDeadlineCloudQueueLimitAssociationIntegration(t *testing.T) {
	suite.Run(t, new(TestDeadlineCloudQueueLimitAssociationSuite))
}

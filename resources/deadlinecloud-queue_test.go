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

type TestDeadlineCloudQueueSuite struct {
	suite.Suite
	svc     *deadline.Client
	farmId  *string
	queueId *string
}

func (suite *TestDeadlineCloudQueueSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = deadline.NewFromConfig(cfg)
}

func (suite *TestDeadlineCloudQueueSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.queueId != nil && suite.farmId != nil {
		_, _ = suite.svc.DeleteQueue(ctx, &deadline.DeleteQueueInput{
			FarmId:  suite.farmId,
			QueueId: suite.queueId,
		})
	}
	if suite.farmId != nil {
		_, _ = suite.svc.DeleteFarm(ctx, &deadline.DeleteFarmInput{
			FarmId: suite.farmId,
		})
	}
}

func (suite *TestDeadlineCloudQueueSuite) TestList() {
	a := assert.New(suite.T())

	lister := DeadlineCloudQueueLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDeadlineCloudQueueIntegration(t *testing.T) {
	suite.Run(t, new(TestDeadlineCloudQueueSuite))
}

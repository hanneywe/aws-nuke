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

type TestDeadlineCloudLimitSuite struct {
	suite.Suite
	svc     *deadline.Client
	farmId  *string
	limitId *string
}

func (suite *TestDeadlineCloudLimitSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = deadline.NewFromConfig(cfg)
}

func (suite *TestDeadlineCloudLimitSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.limitId != nil && suite.farmId != nil {
		_, _ = suite.svc.DeleteLimit(ctx, &deadline.DeleteLimitInput{
			FarmId:  suite.farmId,
			LimitId: suite.limitId,
		})
	}
	if suite.farmId != nil {
		_, _ = suite.svc.DeleteFarm(ctx, &deadline.DeleteFarmInput{
			FarmId: suite.farmId,
		})
	}
}

func (suite *TestDeadlineCloudLimitSuite) TestList() {
	a := assert.New(suite.T())

	lister := DeadlineCloudLimitLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDeadlineCloudLimitIntegration(t *testing.T) {
	suite.Run(t, new(TestDeadlineCloudLimitSuite))
}

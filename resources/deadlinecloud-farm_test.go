//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
)

type TestDeadlineCloudFarmSuite struct {
	suite.Suite
	svc    *deadline.Client
	farmId *string
}

func (suite *TestDeadlineCloudFarmSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = deadline.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := suite.svc.CreateFarm(ctx, &deadline.CreateFarmInput{
		DisplayName: &name,
	})
	if err != nil {
		suite.T().Fatalf("failed to create test farm: %v", err)
	}
	suite.farmId = resp.FarmId
}

func (suite *TestDeadlineCloudFarmSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.farmId != nil {
		_, _ = suite.svc.DeleteFarm(ctx, &deadline.DeleteFarmInput{
			FarmId: suite.farmId,
		})
	}
}

func (suite *TestDeadlineCloudFarmSuite) TestList() {
	a := assert.New(suite.T())

	lister := DeadlineCloudFarmLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestDeadlineCloudFarmSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := DeadlineCloudFarm{
		svc:    suite.svc,
		FarmId: suite.farmId,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.farmId = nil
}

func TestDeadlineCloudFarmIntegration(t *testing.T) {
	suite.Run(t, new(TestDeadlineCloudFarmSuite))
}

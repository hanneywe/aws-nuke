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

type TestDeadlineCloudStorageProfileSuite struct {
	suite.Suite
	svc              *deadline.Client
	farmId           *string
	storageProfileId *string
}

func (suite *TestDeadlineCloudStorageProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = deadline.NewFromConfig(cfg)
}

func (suite *TestDeadlineCloudStorageProfileSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.storageProfileId != nil && suite.farmId != nil {
		_, _ = suite.svc.DeleteStorageProfile(ctx, &deadline.DeleteStorageProfileInput{
			FarmId:           suite.farmId,
			StorageProfileId: suite.storageProfileId,
		})
	}
	if suite.farmId != nil {
		_, _ = suite.svc.DeleteFarm(ctx, &deadline.DeleteFarmInput{
			FarmId: suite.farmId,
		})
	}
}

func (suite *TestDeadlineCloudStorageProfileSuite) TestList() {
	a := assert.New(suite.T())

	lister := DeadlineCloudStorageProfileLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDeadlineCloudListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDeadlineCloudStorageProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestDeadlineCloudStorageProfileSuite))
}

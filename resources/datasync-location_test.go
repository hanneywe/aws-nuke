//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/datasync"
)

type TestDataSyncLocationSuite struct {
	suite.Suite
	svc         *datasync.Client
	locationArn *string
}

func (suite *TestDataSyncLocationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = datasync.NewFromConfig(cfg)
}

func (suite *TestDataSyncLocationSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.locationArn != nil {
		_, _ = suite.svc.DeleteLocation(ctx, &datasync.DeleteLocationInput{
			LocationArn: suite.locationArn,
		})
	}
}

func (suite *TestDataSyncLocationSuite) TestList() {
	a := assert.New(suite.T())

	lister := DataSyncLocationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDataSyncListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDataSyncLocationIntegration(t *testing.T) {
	suite.Run(t, new(TestDataSyncLocationSuite))
}

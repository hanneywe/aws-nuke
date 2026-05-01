//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
)

type TestBatchConsumableResourceSuite struct {
	suite.Suite
	svc *batch.Client
}

func (suite *TestBatchConsumableResourceSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = batch.NewFromConfig(cfg)
}

func (suite *TestBatchConsumableResourceSuite) TestList() {
	a := assert.New(suite.T())
	lister := BatchConsumableResourceLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	_ = resources
}

func TestBatchConsumableResourceIntegration(t *testing.T) {
	suite.Run(t, new(TestBatchConsumableResourceSuite))
}

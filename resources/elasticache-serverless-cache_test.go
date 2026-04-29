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
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

type TestElasticacheServerlessCacheSuite struct {
	suite.Suite
	svc                 *elasticache.Client
	serverlessCacheName *string
}

func (suite *TestElasticacheServerlessCacheSuite) SetupSuite() {
	suite.serverlessCacheName = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = elasticache.NewFromConfig(cfg)

	_, err = suite.svc.CreateServerlessCache(ctx, &elasticache.CreateServerlessCacheInput{
		ServerlessCacheName: suite.serverlessCacheName,
		Engine:              ptr.String("redis"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create test serverless cache: %v", err)
	}
}

func (suite *TestElasticacheServerlessCacheSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteServerlessCache(ctx, &elasticache.DeleteServerlessCacheInput{
		ServerlessCacheName: suite.serverlessCacheName,
	})
}

func (suite *TestElasticacheServerlessCacheSuite) TestList() {
	assertions := assert.New(suite.T())

	lister := ElasticacheServerlessCacheLister{
		svc: suite.svc,
	}

	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)

	assertions.Nil(err)
	assertions.Greater(len(resources), 0)
}

func (suite *TestElasticacheServerlessCacheSuite) TestRemove() {
	assertions := assert.New(suite.T())

	serverlessCache := ElasticacheServerlessCache{
		svc:                 suite.svc,
		ServerlessCacheName: suite.serverlessCacheName,
	}

	err := serverlessCache.Remove(context.TODO())
	assertions.Nil(err)
}

func TestElasticacheServerlessCacheIntegration(t *testing.T) {
	suite.Run(t, new(TestElasticacheServerlessCacheSuite))
}

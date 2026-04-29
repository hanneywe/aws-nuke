package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

func Test_Mock_ElasticacheServerlessCache_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockElasticacheClient)

	mockClient.
		On("DescribeServerlessCaches", mock.Anything, mock.Anything).
		Return(&elasticache.DescribeServerlessCachesOutput{
			ServerlessCaches: []elasticachetypes.ServerlessCache{
				{
					ServerlessCacheName: ptr.String("my-serverless-cache"),
					ARN:                 ptr.String("arn:aws:elasticache:us-east-1:123456789012:serverlesscache:my-serverless-cache"),
					Status:              ptr.String("available"),
				},
			},
		}, nil)

	lister := &ElasticacheServerlessCacheLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	serverlessCache := resources[0].(*ElasticacheServerlessCache)
	assertions.Equal("my-serverless-cache", *serverlessCache.ServerlessCacheName)
	assertions.Equal("arn:aws:elasticache:us-east-1:123456789012:serverlesscache:my-serverless-cache", *serverlessCache.ARN)
	assertions.Equal("available", *serverlessCache.Status)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheServerlessCache_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockElasticacheClient)

	mockClient.
		On("DescribeServerlessCaches", mock.Anything, mock.Anything).
		Return(&elasticache.DescribeServerlessCachesOutput{
			ServerlessCaches: []elasticachetypes.ServerlessCache{},
		}, nil)

	lister := &ElasticacheServerlessCacheLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheServerlessCache_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockElasticacheClient)

	serverlessCache := &ElasticacheServerlessCache{
		svc:                 mockClient,
		ServerlessCacheName: ptr.String("my-serverless-cache"),
		ARN:                 ptr.String("arn:aws:elasticache:us-east-1:123456789012:serverlesscache:my-serverless-cache"),
		Status:              ptr.String("available"),
	}

	mockClient.
		On("DeleteServerlessCache", mock.Anything, &elasticache.DeleteServerlessCacheInput{
			ServerlessCacheName: serverlessCache.ServerlessCacheName,
		}).
		Return(&elasticache.DeleteServerlessCacheOutput{}, nil)

	err := serverlessCache.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheServerlessCache_Properties(t *testing.T) {
	assertions := assert.New(t)

	serverlessCache := ElasticacheServerlessCache{
		ServerlessCacheName: ptr.String("my-serverless-cache"),
		ARN:                 ptr.String("arn:aws:elasticache:us-east-1:123456789012:serverlesscache:my-serverless-cache"),
		Status:              ptr.String("available"),
	}

	properties := serverlessCache.Properties()

	assertions.Equal("my-serverless-cache", properties.Get("ServerlessCacheName"))
	assertions.Equal("arn:aws:elasticache:us-east-1:123456789012:serverlesscache:my-serverless-cache", properties.Get("ARN"))
	assertions.Equal("available", properties.Get("Status"))
}

func Test_Mock_ElasticacheServerlessCache_String(t *testing.T) {
	assertions := assert.New(t)

	serverlessCache := ElasticacheServerlessCache{
		ServerlessCacheName: ptr.String("my-serverless-cache"),
	}

	assertions.Equal("my-serverless-cache", serverlessCache.String())
}

func Test_Mock_ElasticacheServerlessCache_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	serverlessCache := ElasticacheServerlessCache{
		ServerlessCacheName: ptr.String("my-serverless-cache"),
		Status:              ptr.String("deleting"),
	}

	err := serverlessCache.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleting")
}

func Test_Mock_ElasticacheServerlessCache_Filter_Available(t *testing.T) {
	assertions := assert.New(t)

	serverlessCache := ElasticacheServerlessCache{
		ServerlessCacheName: ptr.String("my-serverless-cache"),
		Status:              ptr.String("available"),
	}

	err := serverlessCache.Filter()
	assertions.NoError(err)
}

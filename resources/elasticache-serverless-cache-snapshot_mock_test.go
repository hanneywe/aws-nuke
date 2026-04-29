package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

func Test_Mock_ElastiCacheServerlessCacheSnapshot_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	now := time.Now()
	mockClient.On("DescribeServerlessCacheSnapshots",
		mock.Anything, mock.Anything).
		Return(&elasticache.DescribeServerlessCacheSnapshotsOutput{
			ServerlessCacheSnapshots: []elasticachetypes.ServerlessCacheSnapshot{
				{
					ServerlessCacheSnapshotName: ptr.String("snap-1"),
					ARN:                         ptr.String("arn:snap-1"),
					Status:                      ptr.String("available"),
					SnapshotType:                ptr.String("manual"),
					CreateTime:                  &now,
				},
			},
		}, nil)

	lister := &ElastiCacheServerlessCacheSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ElastiCacheServerlessCacheSnapshot)
	a.Equal("snap-1", *r.ServerlessCacheSnapshotName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElastiCacheServerlessCacheSnapshot_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	mockClient.On("DescribeServerlessCacheSnapshots",
		mock.Anything, mock.Anything).
		Return(&elasticache.DescribeServerlessCacheSnapshotsOutput{
			ServerlessCacheSnapshots: []elasticachetypes.ServerlessCacheSnapshot{},
		}, nil)

	lister := &ElastiCacheServerlessCacheSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElastiCacheServerlessCacheSnapshot_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	r := &ElastiCacheServerlessCacheSnapshot{
		svc:                         mockClient,
		ServerlessCacheSnapshotName: ptr.String("snap-1"),
	}

	mockClient.On("DeleteServerlessCacheSnapshot", mock.Anything,
		&elasticache.DeleteServerlessCacheSnapshotInput{
			ServerlessCacheSnapshotName: r.ServerlessCacheSnapshotName,
		}).Return(&elasticache.DeleteServerlessCacheSnapshotOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ElastiCacheServerlessCacheSnapshot_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ElastiCacheServerlessCacheSnapshot{
		ServerlessCacheSnapshotName: ptr.String("snap-1"),
		ARN:                         ptr.String("arn:snap-1"),
		Status:                      ptr.String("available"),
		SnapshotType:                ptr.String("manual"),
	}
	props := r.Properties()
	a.Equal("snap-1", props.Get("ServerlessCacheSnapshotName"))
	a.Equal("arn:snap-1", props.Get("ARN"))
	a.Equal("available", props.Get("Status"))
	a.Equal("manual", props.Get("SnapshotType"))
}

func Test_Mock_ElastiCacheServerlessCacheSnapshot_String(t *testing.T) {
	a := assert.New(t)
	r := &ElastiCacheServerlessCacheSnapshot{
		ServerlessCacheSnapshotName: ptr.String("snap-1"),
	}
	a.Equal("snap-1", r.String())
}

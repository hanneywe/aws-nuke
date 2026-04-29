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

func Test_Mock_ElasticacheSnapshot_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	mockClient.
		On("DescribeSnapshots", mock.Anything, mock.Anything).
		Return(&elasticache.DescribeSnapshotsOutput{
			Snapshots: []elasticachetypes.Snapshot{
				{
					SnapshotName:   ptr.String("my-snapshot"),
					CacheClusterId: ptr.String("my-cluster"),
					SnapshotStatus: ptr.String("available"),
					ARN:            ptr.String("arn:aws:elasticache:us-east-1:123456789012:snapshot:my-snapshot"),
				},
			},
		}, nil)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(&elasticache.ListTagsForResourceOutput{
			TagList: []elasticachetypes.Tag{
				{Key: ptr.String("env"), Value: ptr.String("test")},
			},
		}, nil)

	lister := &ElasticacheSnapshotLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	snapshot := resources[0].(*ElasticacheSnapshot)
	a.Equal("my-snapshot", *snapshot.SnapshotName)
	a.Equal("my-cluster", *snapshot.CacheClusterID)
	a.Equal("available", *snapshot.Status)
	a.Equal("test", snapshot.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheSnapshot_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	mockClient.
		On("DescribeSnapshots", mock.Anything, mock.Anything).
		Return(&elasticache.DescribeSnapshotsOutput{
			Snapshots: []elasticachetypes.Snapshot{},
		}, nil)

	lister := &ElasticacheSnapshotLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testElasticacheListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheSnapshot_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockElasticacheClient)

	snapshot := &ElasticacheSnapshot{
		svc:          mockClient,
		SnapshotName: ptr.String("my-snapshot"),
		Status:       ptr.String("available"),
	}

	mockClient.
		On("DeleteSnapshot", mock.Anything, &elasticache.DeleteSnapshotInput{
			SnapshotName: snapshot.SnapshotName,
		}).
		Return(&elasticache.DeleteSnapshotOutput{}, nil)

	err := snapshot.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ElasticacheSnapshot_Properties(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName:   ptr.String("my-snapshot"),
		CacheClusterID: ptr.String("my-cluster"),
		Status:         ptr.String("available"),
		Tags:           map[string]string{"env": "test"},
	}

	props := snapshot.Properties()

	a.Equal("my-snapshot", props.Get("SnapshotName"))
	a.Equal("my-cluster", props.Get("CacheClusterID"))
	a.Equal("available", props.Get("Status"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_ElasticacheSnapshot_String(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("my-snapshot"),
	}

	a.Equal("my-snapshot", snapshot.String())
}

func Test_Mock_ElasticacheSnapshot_Filter_Creating(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("my-snapshot"),
		Status:       ptr.String("creating"),
	}

	err := snapshot.Filter()
	a.Error(err)
	a.Contains(err.Error(), "transient state")
}

func Test_Mock_ElasticacheSnapshot_Filter_Deleting(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("my-snapshot"),
		Status:       ptr.String("deleting"),
	}

	err := snapshot.Filter()
	a.Error(err)
	a.Contains(err.Error(), "transient state")
}

func Test_Mock_ElasticacheSnapshot_Filter_Restoring(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("my-snapshot"),
		Status:       ptr.String("restoring"),
	}

	err := snapshot.Filter()
	a.Error(err)
	a.Contains(err.Error(), "transient state")
}

func Test_Mock_ElasticacheSnapshot_Filter_AutomaticPrefix(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("automatic.my-cluster-2024-01-01"),
		Status:       ptr.String("available"),
	}

	err := snapshot.Filter()
	a.Error(err)
	a.Contains(err.Error(), "automatic snapshot")
}

func Test_Mock_ElasticacheSnapshot_Filter_Available(t *testing.T) {
	a := assert.New(t)

	snapshot := ElasticacheSnapshot{
		SnapshotName: ptr.String("my-snapshot"),
		Status:       ptr.String("available"),
	}

	err := snapshot.Filter()
	a.NoError(err)
}

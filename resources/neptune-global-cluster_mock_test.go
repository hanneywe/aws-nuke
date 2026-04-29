package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	neptunetypes "github.com/aws/aws-sdk-go-v2/service/neptune/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

func Test_Mock_NeptuneGlobalCluster_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeGlobalClusters", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeGlobalClustersOutput{
				GlobalClusters: []neptunetypes.GlobalCluster{
					{
						GlobalClusterIdentifier: ptr.String("neptune-global-cluster-1"),
						GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1"),
						Status:                  ptr.String("available"),
						DeletionProtection:      ptr.Bool(false),
					},
				},
			}, nil,
		)

	lister := &NeptuneGlobalClusterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	globalCluster := resources[0].(*NeptuneGlobalCluster)
	assertions.Equal("neptune-global-cluster-1", *globalCluster.GlobalClusterIdentifier)
	assertions.Equal("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1", *globalCluster.GlobalClusterArn)
	assertions.Equal("available", *globalCluster.Status)
	assertions.Equal(false, *globalCluster.DeletionProtection)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGlobalCluster_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeGlobalClusters", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeGlobalClustersOutput{
				GlobalClusters: []neptunetypes.GlobalCluster{},
			}, nil,
		)

	lister := &NeptuneGlobalClusterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGlobalCluster_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", false)

	globalCluster := &NeptuneGlobalCluster{
		svc:                     mockClient,
		settings:                settings,
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(false),
	}

	mockClient.
		On("DeleteGlobalCluster", mock.Anything,
			&neptune.DeleteGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
			},
		).
		Return(&neptune.DeleteGlobalClusterOutput{}, nil)

	err := globalCluster.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGlobalCluster_Remove_WithDeletionProtection(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", true)

	globalCluster := &NeptuneGlobalCluster{
		svc:                     mockClient,
		settings:                settings,
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(true),
	}

	mockClient.
		On("ModifyGlobalCluster", mock.Anything,
			&neptune.ModifyGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
				DeletionProtection:      aws.Bool(false),
			},
		).
		Return(&neptune.ModifyGlobalClusterOutput{}, nil)

	mockClient.
		On("DeleteGlobalCluster", mock.Anything,
			&neptune.DeleteGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
			},
		).
		Return(&neptune.DeleteGlobalClusterOutput{}, nil)

	err := globalCluster.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGlobalCluster_Properties(t *testing.T) {
	assertions := assert.New(t)

	globalCluster := NeptuneGlobalCluster{
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(false),
	}

	properties := globalCluster.Properties()

	assertions.Equal("neptune-global-cluster-1", properties.Get("GlobalClusterIdentifier"))
	assertions.Equal("arn:aws:rds::123456789012:global-cluster:neptune-global-cluster-1", properties.Get("GlobalClusterArn"))
	assertions.Equal("available", properties.Get("Status"))
	assertions.Equal("false", properties.Get("DeletionProtection"))
}

func Test_Mock_NeptuneGlobalCluster_String(t *testing.T) {
	assertions := assert.New(t)

	globalCluster := NeptuneGlobalCluster{
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-1"),
	}

	assertions.Equal("neptune-global-cluster-1", globalCluster.String())
}

func Test_Mock_NeptuneGlobalCluster_Filter(t *testing.T) {
	assertions := assert.New(t)

	// "deleting" status should be filtered
	deletingCluster := NeptuneGlobalCluster{
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-deleting"),
		Status:                  ptr.String("deleting"),
	}
	err := deletingCluster.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleting")

	// "available" status should not be filtered
	availableCluster := NeptuneGlobalCluster{
		GlobalClusterIdentifier: ptr.String("neptune-global-cluster-available"),
		Status:                  ptr.String("available"),
	}
	err = availableCluster.Filter()
	assertions.NoError(err)
}

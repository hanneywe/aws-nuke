package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	docdbtypes "github.com/aws/aws-sdk-go-v2/service/docdb/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testDocDBV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_DocDBGlobalCluster_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDocDBV2Client)

	mockClient.
		On("DescribeGlobalClusters", mock.Anything, mock.Anything).
		Return(
			&docdb.DescribeGlobalClustersOutput{
				GlobalClusters: []docdbtypes.GlobalCluster{
					{
						GlobalClusterIdentifier: ptr.String("global-cluster-1"),
						GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:global-cluster-1"),
						Status:                  ptr.String("available"),
						DeletionProtection:      ptr.Bool(false),
					},
				},
			}, nil,
		)

	lister := &DocDBGlobalClusterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDocDBV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	globalCluster := resources[0].(*DocDBGlobalCluster)
	assertions.Equal("global-cluster-1", *globalCluster.GlobalClusterIdentifier)
	assertions.Equal("arn:aws:rds::123456789012:global-cluster:global-cluster-1", *globalCluster.GlobalClusterArn)
	assertions.Equal("available", *globalCluster.Status)
	assertions.Equal(false, *globalCluster.DeletionProtection)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DocDBGlobalCluster_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDocDBV2Client)

	mockClient.
		On("DescribeGlobalClusters", mock.Anything, mock.Anything).
		Return(
			&docdb.DescribeGlobalClustersOutput{
				GlobalClusters: []docdbtypes.GlobalCluster{},
			}, nil,
		)

	lister := &DocDBGlobalClusterLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDocDBV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DocDBGlobalCluster_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDocDBV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", false)

	globalCluster := &DocDBGlobalCluster{
		svc:                     mockClient,
		settings:                settings,
		GlobalClusterIdentifier: ptr.String("global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(false),
	}

	mockClient.
		On(
			"DeleteGlobalCluster",
			mock.Anything,
			&docdb.DeleteGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
			},
		).
		Return(&docdb.DeleteGlobalClusterOutput{}, nil)

	err := globalCluster.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DocDBGlobalCluster_Remove_WithDeletionProtection(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDocDBV2Client)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", true)

	globalCluster := &DocDBGlobalCluster{
		svc:                     mockClient,
		settings:                settings,
		GlobalClusterIdentifier: ptr.String("global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(true),
	}

	mockClient.
		On("ModifyGlobalCluster", mock.Anything,
			&docdb.ModifyGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
				DeletionProtection:      aws.Bool(false),
			},
		).
		Return(&docdb.ModifyGlobalClusterOutput{}, nil)

	mockClient.
		On("DeleteGlobalCluster", mock.Anything,
			&docdb.DeleteGlobalClusterInput{
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
			},
		).
		Return(&docdb.DeleteGlobalClusterOutput{}, nil)

	err := globalCluster.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DocDBGlobalCluster_Properties(t *testing.T) {
	assertions := assert.New(t)

	globalCluster := DocDBGlobalCluster{
		GlobalClusterIdentifier: ptr.String("global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:global-cluster-1"),
		Status:                  ptr.String("available"),
		DeletionProtection:      ptr.Bool(false),
	}

	properties := globalCluster.Properties()

	assertions.Equal("global-cluster-1", properties.Get("GlobalClusterIdentifier"))
	assertions.Equal("arn:aws:rds::123456789012:global-cluster:global-cluster-1", properties.Get("GlobalClusterArn"))
	assertions.Equal("available", properties.Get("Status"))
	assertions.Equal("false", properties.Get("DeletionProtection"))
}

func Test_Mock_DocDBGlobalCluster_String(t *testing.T) {
	assertions := assert.New(t)

	globalCluster := DocDBGlobalCluster{
		GlobalClusterIdentifier: ptr.String("global-cluster-1"),
		GlobalClusterArn:        ptr.String("arn:aws:rds::123456789012:global-cluster:global-cluster-1"),
	}

	assertions.Equal("global-cluster-1", globalCluster.String())
}

func Test_Mock_DocDBGlobalCluster_Filter(t *testing.T) {
	assertions := assert.New(t)

	// "deleting" status should be filtered
	deletingCluster := DocDBGlobalCluster{
		GlobalClusterIdentifier: ptr.String("global-cluster-deleting"),
		Status:                  ptr.String("deleting"),
	}
	err := deletingCluster.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleting")

	// "available" status should not be filtered
	availableCluster := DocDBGlobalCluster{
		GlobalClusterIdentifier: ptr.String("global-cluster-available"),
		Status:                  ptr.String("available"),
	}
	err = availableCluster.Filter()
	assertions.NoError(err)
}

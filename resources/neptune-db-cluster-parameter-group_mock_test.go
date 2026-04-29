package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
	neptunetypes "github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func Test_Mock_NeptuneDBClusterParameterGroup_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBClusterParameterGroups", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBClusterParameterGroupsOutput{
				DBClusterParameterGroups: []neptunetypes.DBClusterParameterGroup{
					{
						DBClusterParameterGroupName: ptr.String("custom-cluster-params"),
						DBClusterParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:cpg:custom-cluster-params"),
					},
				},
			}, nil,
		)

	lister := &NeptuneDBClusterParameterGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	clusterParameterGroup := resources[0].(*NeptuneDBClusterParameterGroup)
	assertions.Equal("custom-cluster-params", *clusterParameterGroup.DBClusterParameterGroupName)
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:cpg:custom-cluster-params", *clusterParameterGroup.DBClusterParameterGroupArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterParameterGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBClusterParameterGroups", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBClusterParameterGroupsOutput{
				DBClusterParameterGroups: []neptunetypes.DBClusterParameterGroup{},
			}, nil,
		)

	lister := &NeptuneDBClusterParameterGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterParameterGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	clusterParameterGroup := &NeptuneDBClusterParameterGroup{
		svc:                         mockClient,
		DBClusterParameterGroupName: ptr.String("custom-cluster-params"),
		DBClusterParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:cpg:custom-cluster-params"),
	}

	mockClient.
		On("DeleteDBClusterParameterGroup", mock.Anything,
			&neptune.DeleteDBClusterParameterGroupInput{
				DBClusterParameterGroupName: clusterParameterGroup.DBClusterParameterGroupName,
			},
		).
		Return(&neptune.DeleteDBClusterParameterGroupOutput{}, nil)

	err := clusterParameterGroup.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterParameterGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	clusterParameterGroup := NeptuneDBClusterParameterGroup{
		DBClusterParameterGroupName: ptr.String("custom-cluster-params"),
		DBClusterParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:cpg:custom-cluster-params"),
	}

	properties := clusterParameterGroup.Properties()

	assertions.Equal("custom-cluster-params", properties.Get("DBClusterParameterGroupName"))
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:cpg:custom-cluster-params", properties.Get("DBClusterParameterGroupArn"))
}

func Test_Mock_NeptuneDBClusterParameterGroup_String(t *testing.T) {
	assertions := assert.New(t)

	clusterParameterGroup := NeptuneDBClusterParameterGroup{
		DBClusterParameterGroupName: ptr.String("custom-cluster-params"),
	}

	assertions.Equal("custom-cluster-params", clusterParameterGroup.String())
}

func Test_Mock_NeptuneDBClusterParameterGroup_Filter(t *testing.T) {
	assertions := assert.New(t)

	// Default cluster parameter group should be filtered
	defaultClusterParameterGroup := NeptuneDBClusterParameterGroup{
		DBClusterParameterGroupName: ptr.String("default.neptune1"),
	}
	err := defaultClusterParameterGroup.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "default")

	// Custom cluster parameter group should not be filtered
	customClusterParameterGroup := NeptuneDBClusterParameterGroup{
		DBClusterParameterGroupName: ptr.String("custom-cluster-params"),
	}
	err = customClusterParameterGroup.Filter()
	assertions.NoError(err)
}

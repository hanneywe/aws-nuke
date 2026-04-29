package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
	neptunetypes "github.com/aws/aws-sdk-go-v2/service/neptune/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testNeptuneV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_NeptuneDBParameterGroup_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBParameterGroups", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBParameterGroupsOutput{
				DBParameterGroups: []neptunetypes.DBParameterGroup{
					{
						DBParameterGroupName: ptr.String("custom-neptune-params"),
						DBParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:pg:custom-neptune-params"),
					},
				},
			}, nil,
		)

	lister := &NeptuneDBParameterGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	parameterGroup := resources[0].(*NeptuneDBParameterGroup)
	assertions.Equal("custom-neptune-params", *parameterGroup.DBParameterGroupName)
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:pg:custom-neptune-params", *parameterGroup.DBParameterGroupArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBParameterGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBParameterGroups", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBParameterGroupsOutput{
				DBParameterGroups: []neptunetypes.DBParameterGroup{},
			}, nil,
		)

	lister := &NeptuneDBParameterGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBParameterGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	parameterGroup := &NeptuneDBParameterGroup{
		svc:                  mockClient,
		DBParameterGroupName: ptr.String("custom-neptune-params"),
		DBParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:pg:custom-neptune-params"),
	}

	mockClient.
		On("DeleteDBParameterGroup", mock.Anything,
			&neptune.DeleteDBParameterGroupInput{
				DBParameterGroupName: parameterGroup.DBParameterGroupName,
			},
		).
		Return(&neptune.DeleteDBParameterGroupOutput{}, nil)

	err := parameterGroup.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBParameterGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	parameterGroup := NeptuneDBParameterGroup{
		DBParameterGroupName: ptr.String("custom-neptune-params"),
		DBParameterGroupArn:  ptr.String("arn:aws:rds:us-east-1:123456789012:pg:custom-neptune-params"),
	}

	properties := parameterGroup.Properties()

	assertions.Equal("custom-neptune-params", properties.Get("DBParameterGroupName"))
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:pg:custom-neptune-params", properties.Get("DBParameterGroupArn"))
}

func Test_Mock_NeptuneDBParameterGroup_String(t *testing.T) {
	assertions := assert.New(t)

	parameterGroup := NeptuneDBParameterGroup{
		DBParameterGroupName: ptr.String("custom-neptune-params"),
	}

	assertions.Equal("custom-neptune-params", parameterGroup.String())
}

func Test_Mock_NeptuneDBParameterGroup_Filter(t *testing.T) {
	assertions := assert.New(t)

	// Default parameter group should be filtered
	defaultParameterGroup := NeptuneDBParameterGroup{
		DBParameterGroupName: ptr.String("default.neptune1"),
	}
	err := defaultParameterGroup.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "default")

	// Custom parameter group should not be filtered
	customParameterGroup := NeptuneDBParameterGroup{
		DBParameterGroupName: ptr.String("custom-neptune-params"),
	}
	err = customParameterGroup.Filter()
	assertions.NoError(err)
}

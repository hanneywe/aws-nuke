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

func Test_Mock_NeptuneDBClusterEndpoint_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBClusterEndpoints", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBClusterEndpointsOutput{
				DBClusterEndpoints: []neptunetypes.DBClusterEndpoint{
					{
						DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
						DBClusterIdentifier:         ptr.String("my-cluster"),
						EndpointType:                ptr.String("CUSTOM"),
						Status:                      ptr.String("available"),
					},
				},
			}, nil,
		)

	lister := &NeptuneDBClusterEndpointLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	endpoint := resources[0].(*NeptuneDBClusterEndpoint)
	assertions.Equal("my-endpoint", *endpoint.DBClusterEndpointIdentifier)
	assertions.Equal("my-cluster", *endpoint.DBClusterIdentifier)
	assertions.Equal("CUSTOM", *endpoint.EndpointType)
	assertions.Equal("available", *endpoint.Status)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterEndpoint_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeDBClusterEndpoints", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeDBClusterEndpointsOutput{
				DBClusterEndpoints: []neptunetypes.DBClusterEndpoint{},
			}, nil,
		)

	lister := &NeptuneDBClusterEndpointLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterEndpoint_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	endpoint := &NeptuneDBClusterEndpoint{
		svc:                         mockClient,
		DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
		DBClusterIdentifier:         ptr.String("my-cluster"),
		EndpointType:                ptr.String("CUSTOM"),
		Status:                      ptr.String("available"),
	}

	mockClient.
		On("DeleteDBClusterEndpoint", mock.Anything,
			&neptune.DeleteDBClusterEndpointInput{
				DBClusterEndpointIdentifier: endpoint.DBClusterEndpointIdentifier,
			},
		).
		Return(&neptune.DeleteDBClusterEndpointOutput{}, nil)

	err := endpoint.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneDBClusterEndpoint_Properties(t *testing.T) {
	assertions := assert.New(t)

	endpoint := NeptuneDBClusterEndpoint{
		DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
		DBClusterIdentifier:         ptr.String("my-cluster"),
		EndpointType:                ptr.String("CUSTOM"),
		Status:                      ptr.String("available"),
	}

	properties := endpoint.Properties()

	assertions.Equal("my-endpoint", properties.Get("DBClusterEndpointIdentifier"))
	assertions.Equal("my-cluster", properties.Get("DBClusterIdentifier"))
	assertions.Equal("CUSTOM", properties.Get("EndpointType"))
	assertions.Equal("available", properties.Get("Status"))
}

func Test_Mock_NeptuneDBClusterEndpoint_String(t *testing.T) {
	assertions := assert.New(t)

	endpoint := NeptuneDBClusterEndpoint{
		DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
	}

	assertions.Equal("my-endpoint", endpoint.String())
}

func Test_Mock_NeptuneDBClusterEndpoint_Filter_Deleting(t *testing.T) {
	assertions := assert.New(t)

	endpoint := NeptuneDBClusterEndpoint{
		DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
		Status:                      ptr.String("deleting"),
	}
	err := endpoint.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleting")
}

func Test_Mock_NeptuneDBClusterEndpoint_Filter_Available(t *testing.T) {
	assertions := assert.New(t)

	endpoint := NeptuneDBClusterEndpoint{
		DBClusterEndpointIdentifier: ptr.String("my-endpoint"),
		Status:                      ptr.String("available"),
	}
	err := endpoint.Filter()
	assertions.NoError(err)
}

package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

func Test_Mock_ConnectRoutingProfile_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("test-instanceid")},
			},
		}, nil)

	mockClient.On("ListRoutingProfiles", mock.Anything, mock.Anything).
		Return(&connect.ListRoutingProfilesOutput{
			RoutingProfileSummaryList: []connecttypes.RoutingProfileSummary{
				{Id: ptr.String("test-routingprofileid"), Name: ptr.String("test-name")},
			},
		}, nil)

	lister := &ConnectRoutingProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ConnectRoutingProfile)
	a.Equal("test-instanceid", *r.InstanceID)
	a.Equal("test-routingprofileid", *r.RoutingProfileID)
	a.Equal("test-name", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectRoutingProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectRoutingProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectRoutingProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	r := &ConnectRoutingProfile{
		svc:              mockClient,
		InstanceID:       ptr.String("test-instanceid"),
		RoutingProfileID: ptr.String("test-routingprofileid"),
	}

	mockClient.On("DeleteRoutingProfile", mock.Anything,
		&connect.DeleteRoutingProfileInput{
			InstanceId:       r.InstanceID,
			RoutingProfileId: r.RoutingProfileID,
		}).Return(&connect.DeleteRoutingProfileOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectRoutingProfile_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ConnectRoutingProfile{
		InstanceID:       ptr.String("test-instanceid"),
		RoutingProfileID: ptr.String("test-routingprofileid"),
		Name:             ptr.String("test-name"),
	}
	props := r.Properties()
	a.Equal("test-instanceid", props.Get("InstanceId"))
	a.Equal("test-routingprofileid", props.Get("RoutingProfileId"))
	a.Equal("test-name", props.Get("Name"))
}

func Test_Mock_ConnectRoutingProfile_String(t *testing.T) {
	a := assert.New(t)
	r := &ConnectRoutingProfile{
		Name: ptr.String("test-name"),
	}
	a.Equal("test-name", r.String())
}

func Test_Mock_ConnectRoutingProfile_Filter(t *testing.T) {
	a := assert.New(t)

	r := &ConnectRoutingProfile{
		Name: ptr.String("Basic Routing Profile"),
	}
	a.Error(r.Filter())

	r2 := &ConnectRoutingProfile{
		Name: ptr.String("Custom Routing Profile"),
	}
	a.NoError(r2.Filter())
}

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

func Test_Mock_ConnectHoursOfOperation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("i-12345")},
			},
		}, nil)

	mockClient.
		On("ListHoursOfOperations", mock.Anything, mock.Anything).
		Return(&connect.ListHoursOfOperationsOutput{
			HoursOfOperationSummaryList: []connecttypes.HoursOfOperationSummary{
				{
					Id:   ptr.String("hoo-12345"),
					Name: ptr.String("Custom Hours"),
				},
			},
		}, nil)

	lister := &ConnectHoursOfOperationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	hoo := resources[0].(*ConnectHoursOfOperation)
	a.Equal("hoo-12345", *hoo.ID)
	a.Equal("Custom Hours", *hoo.Name)
	a.Equal("i-12345", *hoo.InstanceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectHoursOfOperation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectHoursOfOperationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectHoursOfOperation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	hoo := &ConnectHoursOfOperation{
		svc:        mockClient,
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("hoo-12345"),
		Name:       ptr.String("Custom Hours"),
	}

	mockClient.
		On("DeleteHoursOfOperation", mock.Anything, &connect.DeleteHoursOfOperationInput{
			InstanceId:         hoo.InstanceID,
			HoursOfOperationId: hoo.ID,
		}).
		Return(&connect.DeleteHoursOfOperationOutput{}, nil)

	err := hoo.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectHoursOfOperation_Filter_Default(t *testing.T) {
	a := assert.New(t)

	hoo := &ConnectHoursOfOperation{
		Name: ptr.String("Basic Hours"),
	}

	err := hoo.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete default hours of operation")
}

func Test_Mock_ConnectHoursOfOperation_Filter_Custom(t *testing.T) {
	a := assert.New(t)

	hoo := &ConnectHoursOfOperation{
		Name: ptr.String("Custom Hours"),
	}

	err := hoo.Filter()
	a.NoError(err)
}

func Test_Mock_ConnectHoursOfOperation_Properties(t *testing.T) {
	a := assert.New(t)

	hoo := ConnectHoursOfOperation{
		InstanceID: ptr.String("i-12345"),
		ID:         ptr.String("hoo-12345"),
		Name:       ptr.String("Custom Hours"),
	}

	props := hoo.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("hoo-12345", props.Get("Id"))
	a.Equal("Custom Hours", props.Get("Name"))
}

func Test_Mock_ConnectHoursOfOperation_String(t *testing.T) {
	a := assert.New(t)

	hoo := ConnectHoursOfOperation{
		Name: ptr.String("Custom Hours"),
	}

	a.Equal("Custom Hours", hoo.String())
}

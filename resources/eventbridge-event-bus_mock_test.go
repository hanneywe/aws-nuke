package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func Test_Mock_EventBridgeEventBus_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.On("ListEventBuses", mock.Anything, mock.Anything).
		Return(&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{
				{Name: ptr.String("my-bus"), Arn: ptr.String("arn:aws:events:us-east-1:123456789012:event-bus/my-bus")},
			},
		}, nil)

	lister := &EventBridgeEventBusLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*EventBridgeEventBus)
	a.Equal("my-bus", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeEventBus_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.On("ListEventBuses", mock.Anything, mock.Anything).
		Return(&eventbridge.ListEventBusesOutput{
			EventBuses: []eventbridgetypes.EventBus{},
		}, nil)

	lister := &EventBridgeEventBusLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeEventBus_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	r := &EventBridgeEventBus{
		svc:  mockClient,
		Name: ptr.String("my-bus"),
	}

	mockClient.On("DeleteEventBus", mock.Anything,
		&eventbridge.DeleteEventBusInput{
			Name: r.Name,
		}).Return(&eventbridge.DeleteEventBusOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeEventBus_Properties(t *testing.T) {
	a := assert.New(t)
	r := &EventBridgeEventBus{
		Name: ptr.String("my-bus"),
		Arn:  ptr.String("arn:aws:events:us-east-1:123456789012:event-bus/my-bus"),
	}
	props := r.Properties()
	a.Equal("my-bus", props.Get("Name"))
	a.Equal("arn:aws:events:us-east-1:123456789012:event-bus/my-bus", props.Get("Arn"))
}

func Test_Mock_EventBridgeEventBus_String(t *testing.T) {
	a := assert.New(t)
	r := &EventBridgeEventBus{
		Name: ptr.String("my-bus"),
	}
	a.Equal("my-bus", r.String())
}

func Test_Mock_EventBridgeEventBus_Filter(t *testing.T) {
	a := assert.New(t)

	r := &EventBridgeEventBus{
		Name: ptr.String("default"),
	}
	a.Error(r.Filter())
	a.Contains(r.Filter().Error(), "cannot delete default event bus")

	r2 := &EventBridgeEventBus{
		Name: ptr.String("my-custom-bus"),
	}
	a.NoError(r2.Filter())
}

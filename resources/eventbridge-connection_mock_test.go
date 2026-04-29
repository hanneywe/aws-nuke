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

func Test_Mock_EventBridgeConnection_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.
		On("ListConnections", mock.Anything, mock.Anything).
		Return(
			&eventbridge.ListConnectionsOutput{
				Connections: []eventbridgetypes.Connection{
					{
						Name:            ptr.String("test-connection"),
						ConnectionArn:   ptr.String("arn:aws:events:us-east-1:123456789012:connection/test-connection"),
						ConnectionState: eventbridgetypes.ConnectionStateAuthorized,
					},
				},
			}, nil,
		)

	lister := &EventBridgeConnectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	connection := resources[0].(*EventBridgeConnection)
	assertions.Equal("test-connection", *connection.Name)
	assertions.Equal("arn:aws:events:us-east-1:123456789012:connection/test-connection", *connection.ConnectionArn)
	assertions.Equal("AUTHORIZED", *connection.ConnectionState)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeConnection_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.
		On("ListConnections", mock.Anything, mock.Anything).
		Return(
			&eventbridge.ListConnectionsOutput{
				Connections: []eventbridgetypes.Connection{},
			}, nil,
		)

	lister := &EventBridgeConnectionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeConnection_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	connection := &EventBridgeConnection{
		svc:  mockClient,
		Name: ptr.String("test-connection"),
	}

	mockClient.
		On(
			"DeleteConnection",
			mock.Anything,
			&eventbridge.DeleteConnectionInput{
				Name: connection.Name,
			},
		).
		Return(&eventbridge.DeleteConnectionOutput{}, nil)

	err := connection.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeConnection_Properties(t *testing.T) {
	assertions := assert.New(t)

	connection := EventBridgeConnection{
		Name:            ptr.String("test-connection"),
		ConnectionArn:   ptr.String("arn:aws:events:us-east-1:123456789012:connection/test-connection"),
		ConnectionState: ptr.String("AUTHORIZED"),
	}

	properties := connection.Properties()

	assertions.Equal("test-connection", properties.Get("Name"))
	assertions.Equal("arn:aws:events:us-east-1:123456789012:connection/test-connection", properties.Get("ConnectionArn"))
	assertions.Equal("AUTHORIZED", properties.Get("ConnectionState"))
}

func Test_Mock_EventBridgeConnection_String(t *testing.T) {
	assertions := assert.New(t)

	connection := EventBridgeConnection{
		Name: ptr.String("test-connection"),
	}

	assertions.Equal("test-connection", connection.String())
}

func Test_Mock_EventBridgeConnection_Filter(t *testing.T) {
	assertions := assert.New(t)

	deletingConnection := EventBridgeConnection{ConnectionState: ptr.String("DELETING")}
	assertions.Error(deletingConnection.Filter())

	authorizedConnection := EventBridgeConnection{ConnectionState: ptr.String("AUTHORIZED")}
	assertions.NoError(authorizedConnection.Filter())

	activeConnection := EventBridgeConnection{ConnectionState: ptr.String("ACTIVE")}
	assertions.NoError(activeConnection.Filter())
}

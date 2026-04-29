package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	codeconnectionstypes "github.com/aws/aws-sdk-go-v2/service/codeconnections/types"
)

func Test_Mock_CodeConnectionsConnection_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)

	mockClient.On("ListConnections", mock.Anything, mock.Anything).
		Return(&codeconnections.ListConnectionsOutput{
			Connections: []codeconnectionstypes.Connection{
				{
					ConnectionArn:    ptr.String("arn:aws:codeconnections:us-east-1:123456789012:connection/test-id"),
					ConnectionName:   ptr.String("test-connection"),
					ConnectionStatus: codeconnectionstypes.ConnectionStatusAvailable,
					ProviderType:     codeconnectionstypes.ProviderTypeGithub,
				},
			},
		}, nil)

	lister := &CodeConnectionsConnectionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CodeConnectionsConnection)
	a.Equal("arn:aws:codeconnections:us-east-1:123456789012:connection/test-id", *r.ConnectionArn)
	a.Equal("test-connection", *r.ConnectionName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsConnection_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)

	mockClient.On("ListConnections", mock.Anything, mock.Anything).
		Return(&codeconnections.ListConnectionsOutput{
			Connections: []codeconnectionstypes.Connection{},
		}, nil)

	lister := &CodeConnectionsConnectionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsConnection_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)

	r := &CodeConnectionsConnection{
		svc:           mockClient,
		ConnectionArn: ptr.String("arn:aws:codeconnections:us-east-1:123456789012:connection/test-id"),
	}

	mockClient.On("DeleteConnection", mock.Anything,
		&codeconnections.DeleteConnectionInput{
			ConnectionArn: r.ConnectionArn,
		}).Return(&codeconnections.DeleteConnectionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsConnection_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CodeConnectionsConnection{
		ConnectionArn:    ptr.String("arn:aws:codeconnections:us-east-1:123456789012:connection/test-id"),
		ConnectionName:   ptr.String("test-connection"),
		ConnectionStatus: codeconnectionstypes.ConnectionStatusAvailable,
		ProviderType:     codeconnectionstypes.ProviderTypeGithub,
	}
	props := r.Properties()
	a.Equal("arn:aws:codeconnections:us-east-1:123456789012:connection/test-id", props.Get("ConnectionArn"))
	a.Equal("test-connection", props.Get("ConnectionName"))
}

func Test_Mock_CodeConnectionsConnection_String(t *testing.T) {
	a := assert.New(t)
	r := &CodeConnectionsConnection{
		ConnectionName: ptr.String("test-connection"),
	}
	a.Equal("test-connection", r.String())
}

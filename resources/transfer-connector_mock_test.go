package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
	transfertypes "github.com/aws/aws-sdk-go-v2/service/transfer/types"
)

func Test_Mock_TransferConnector_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListConnectors", mock.Anything, mock.Anything).
		Return(&transfer.ListConnectorsOutput{
			Connectors: []transfertypes.ListedConnector{
				{ConnectorId: ptr.String("test-value"), Url: ptr.String("test-value")},
			},
		}, nil)

	lister := &TransferConnectorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*TransferConnector)
	a.Equal("test-value", *r.ConnectorID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferConnector_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListConnectors", mock.Anything, mock.Anything).
		Return(&transfer.ListConnectorsOutput{
			Connectors: []transfertypes.ListedConnector{},
		}, nil)

	lister := &TransferConnectorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferConnector_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	r := &TransferConnector{
		svc:         mockClient,
		ConnectorID: ptr.String("test-connectorid"),
	}

	mockClient.On("DeleteConnector", mock.Anything,
		&transfer.DeleteConnectorInput{
			ConnectorId: r.ConnectorID,
		}).Return(&transfer.DeleteConnectorOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferConnector_Properties(t *testing.T) {
	a := assert.New(t)
	r := &TransferConnector{
		ConnectorID: ptr.String("test-connectorid"),
		URL:         ptr.String("test-url"),
	}
	props := r.Properties()
	a.Equal("test-connectorid", props.Get("ConnectorID"))
	a.Equal("test-url", props.Get("URL"))
}

func Test_Mock_TransferConnector_String(t *testing.T) {
	a := assert.New(t)
	r := &TransferConnector{
		ConnectorID: ptr.String("test-connectorid"),
	}
	a.Equal("test-connectorid", r.String())
}

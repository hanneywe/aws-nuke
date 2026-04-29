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

func Test_Mock_ConnectIntegrationAssociation_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("test-instanceid")},
			},
		}, nil)

	mockClient.On("ListIntegrationAssociations", mock.Anything, mock.Anything).
		Return(&connect.ListIntegrationAssociationsOutput{
			IntegrationAssociationSummaryList: []connecttypes.IntegrationAssociationSummary{
				{IntegrationAssociationId: ptr.String("test-integrationassociationid"), IntegrationArn: ptr.String("test-integrationarn")},
			},
		}, nil)

	lister := &ConnectIntegrationAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ConnectIntegrationAssociation)
	a.Equal("test-instanceid", *r.InstanceID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectIntegrationAssociation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectIntegrationAssociationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectIntegrationAssociation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	r := &ConnectIntegrationAssociation{
		svc:                      mockClient,
		InstanceID:               ptr.String("test-instanceid"),
		IntegrationAssociationID: ptr.String("test-integrationassociationid"),
	}

	mockClient.On("DeleteIntegrationAssociation", mock.Anything,
		&connect.DeleteIntegrationAssociationInput{
			InstanceId:               r.InstanceID,
			IntegrationAssociationId: r.IntegrationAssociationID,
		}).Return(&connect.DeleteIntegrationAssociationOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectIntegrationAssociation_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ConnectIntegrationAssociation{
		InstanceID:               ptr.String("test-instanceid"),
		IntegrationAssociationID: ptr.String("test-integrationassociationid"),
		IntegrationARN:           ptr.String("test-integrationarn"),
	}
	props := r.Properties()
	a.Equal("test-instanceid", props.Get("InstanceId"))
	a.Equal("test-integrationassociationid", props.Get("IntegrationAssociationId"))
	a.Equal("test-integrationarn", props.Get("IntegrationArn"))
}

func Test_Mock_ConnectIntegrationAssociation_String(t *testing.T) {
	a := assert.New(t)
	r := &ConnectIntegrationAssociation{
		IntegrationAssociationID: ptr.String("test-integrationassociationid"),
	}
	a.Equal("test-integrationassociationid", r.String())
}

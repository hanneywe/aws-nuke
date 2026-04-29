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

func Test_Mock_ConnectPredefinedAttribute_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("instance-123")},
			},
		}, nil)

	mockClient.On("ListPredefinedAttributes", mock.Anything, mock.Anything).
		Return(&connect.ListPredefinedAttributesOutput{
			PredefinedAttributeSummaryList: []connecttypes.PredefinedAttributeSummary{
				{Name: ptr.String("Department")},
			},
		}, nil)

	lister := &ConnectPredefinedAttributeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	predefinedAttribute := resources[0].(*ConnectPredefinedAttribute)
	assertions.Equal("Department", *predefinedAttribute.Name)
	assertions.Equal("instance-123", *predefinedAttribute.InstanceID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectPredefinedAttribute_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("instance-123")},
			},
		}, nil)

	mockClient.On("ListPredefinedAttributes", mock.Anything, mock.Anything).
		Return(&connect.ListPredefinedAttributesOutput{
			PredefinedAttributeSummaryList: []connecttypes.PredefinedAttributeSummary{},
		}, nil)

	lister := &ConnectPredefinedAttributeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectPredefinedAttribute_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockConnectClient)

	predefinedAttribute := &ConnectPredefinedAttribute{
		svc:        mockClient,
		Name:       ptr.String("Department"),
		InstanceID: ptr.String("instance-123"),
	}

	mockClient.On("DeletePredefinedAttribute", mock.Anything, &connect.DeletePredefinedAttributeInput{
		InstanceId: predefinedAttribute.InstanceID,
		Name:       predefinedAttribute.Name,
	}).Return(&connect.DeletePredefinedAttributeOutput{}, nil)

	err := predefinedAttribute.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectPredefinedAttribute_Properties(t *testing.T) {
	assertions := assert.New(t)

	predefinedAttribute := ConnectPredefinedAttribute{
		Name:       ptr.String("Department"),
		InstanceID: ptr.String("instance-123"),
	}

	properties := predefinedAttribute.Properties()
	assertions.Equal("Department", properties.Get("Name"))
	assertions.Equal("instance-123", properties.Get("InstanceId"))
}

func Test_Mock_ConnectPredefinedAttribute_String(t *testing.T) {
	assertions := assert.New(t)
	predefinedAttribute := ConnectPredefinedAttribute{Name: ptr.String("Department")}
	assertions.Equal("Department", predefinedAttribute.String())
}

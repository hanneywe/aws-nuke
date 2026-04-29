package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"
	greengrasstypes "github.com/aws/aws-sdk-go-v2/service/greengrass/types"
)

func Test_Mock_GreengrassDeviceDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListDeviceDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListDeviceDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-device")},
			},
		}, nil)

	lister := &GreengrassDeviceDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassDeviceDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassDeviceDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListDeviceDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListDeviceDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassDeviceDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassDeviceDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassDeviceDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteDeviceDefinition", mock.Anything, &greengrass.DeleteDeviceDefinitionInput{
		DeviceDefinitionId: r.ID,
	}).Return(&greengrass.DeleteDeviceDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassDeviceDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassDeviceDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-device")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-device", props.Get("Name"))
}

func Test_Mock_GreengrassDeviceDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassDeviceDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

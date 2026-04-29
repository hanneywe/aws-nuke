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

func Test_Mock_GreengrassResourceDefinition_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListResourceDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListResourceDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("res-def-1"), Name: ptr.String("my-resource-def")},
			},
		}, nil)

	lister := &GreengrassResourceDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	resourceDef := resources[0].(*GreengrassResourceDefinition)
	assertions.Equal("res-def-1", *resourceDef.ID)
	assertions.Equal("my-resource-def", *resourceDef.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassResourceDefinition_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListResourceDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListResourceDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassResourceDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassResourceDefinition_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockGreengrassClient)

	resourceDef := &GreengrassResourceDefinition{
		svc: mockClient,
		ID:  ptr.String("res-def-1"),
	}

	mockClient.On("DeleteResourceDefinition", mock.Anything, &greengrass.DeleteResourceDefinitionInput{
		ResourceDefinitionId: resourceDef.ID,
	}).Return(&greengrass.DeleteResourceDefinitionOutput{}, nil)

	assertions.NoError(resourceDef.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassResourceDefinition_Properties(t *testing.T) {
	assertions := assert.New(t)

	resourceDef := GreengrassResourceDefinition{
		ID:   ptr.String("res-def-1"),
		Name: ptr.String("my-resource-def"),
	}

	properties := resourceDef.Properties()
	assertions.Equal("res-def-1", properties.Get("Id"))
	assertions.Equal("my-resource-def", properties.Get("Name"))
}

func Test_Mock_GreengrassResourceDefinition_String(t *testing.T) {
	assertions := assert.New(t)

	resourceDef := GreengrassResourceDefinition{ID: ptr.String("res-def-1")}
	assertions.Equal("res-def-1", resourceDef.String())
}

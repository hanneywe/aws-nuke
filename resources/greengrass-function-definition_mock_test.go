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

func Test_Mock_GreengrassFunctionDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListFunctionDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListFunctionDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-function")},
			},
		}, nil)

	lister := &GreengrassFunctionDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassFunctionDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassFunctionDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListFunctionDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListFunctionDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassFunctionDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassFunctionDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassFunctionDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteFunctionDefinition", mock.Anything, &greengrass.DeleteFunctionDefinitionInput{
		FunctionDefinitionId: r.ID,
	}).Return(&greengrass.DeleteFunctionDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassFunctionDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassFunctionDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-function")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-function", props.Get("Name"))
}

func Test_Mock_GreengrassFunctionDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassFunctionDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

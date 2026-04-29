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

func Test_Mock_GreengrassConnectorDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListConnectorDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListConnectorDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-connector")},
			},
		}, nil)

	lister := &GreengrassConnectorDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassConnectorDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassConnectorDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListConnectorDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListConnectorDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassConnectorDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassConnectorDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassConnectorDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteConnectorDefinition", mock.Anything, &greengrass.DeleteConnectorDefinitionInput{
		ConnectorDefinitionId: r.ID,
	}).Return(&greengrass.DeleteConnectorDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassConnectorDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassConnectorDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-connector")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-connector", props.Get("Name"))
}

func Test_Mock_GreengrassConnectorDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassConnectorDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

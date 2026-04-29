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

func Test_Mock_GreengrassCoreDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListCoreDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListCoreDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-core")},
			},
		}, nil)

	lister := &GreengrassCoreDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassCoreDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassCoreDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListCoreDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListCoreDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassCoreDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassCoreDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassCoreDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteCoreDefinition", mock.Anything, &greengrass.DeleteCoreDefinitionInput{
		CoreDefinitionId: r.ID,
	}).Return(&greengrass.DeleteCoreDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassCoreDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassCoreDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-core")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-core", props.Get("Name"))
}

func Test_Mock_GreengrassCoreDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassCoreDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

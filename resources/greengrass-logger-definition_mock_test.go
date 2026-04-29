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

func Test_Mock_GreengrassLoggerDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListLoggerDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListLoggerDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-logger")},
			},
		}, nil)

	lister := &GreengrassLoggerDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassLoggerDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassLoggerDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListLoggerDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListLoggerDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassLoggerDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassLoggerDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassLoggerDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteLoggerDefinition", mock.Anything, &greengrass.DeleteLoggerDefinitionInput{
		LoggerDefinitionId: r.ID,
	}).Return(&greengrass.DeleteLoggerDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassLoggerDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassLoggerDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-logger")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-logger", props.Get("Name"))
}

func Test_Mock_GreengrassLoggerDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassLoggerDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

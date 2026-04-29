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

func Test_Mock_GreengrassSubscriptionDefinition_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListSubscriptionDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListSubscriptionDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{
				{Id: ptr.String("def-1"), Name: ptr.String("my-subscription")},
			},
		}, nil)

	lister := &GreengrassSubscriptionDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("def-1", *resources[0].(*GreengrassSubscriptionDefinition).ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassSubscriptionDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	mockClient.On("ListSubscriptionDefinitions", mock.Anything, mock.Anything).
		Return(&greengrass.ListSubscriptionDefinitionsOutput{
			Definitions: []greengrasstypes.DefinitionInformation{},
		}, nil)

	lister := &GreengrassSubscriptionDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGreengrassListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassSubscriptionDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGreengrassClient)

	r := &GreengrassSubscriptionDefinition{svc: mockClient, ID: ptr.String("def-1")}
	mockClient.On("DeleteSubscriptionDefinition", mock.Anything, &greengrass.DeleteSubscriptionDefinitionInput{
		SubscriptionDefinitionId: r.ID,
	}).Return(&greengrass.DeleteSubscriptionDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GreengrassSubscriptionDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := GreengrassSubscriptionDefinition{ID: ptr.String("def-1"), Name: ptr.String("my-subscription")}
	props := r.Properties()
	a.Equal("def-1", props.Get("Id"))
	a.Equal("my-subscription", props.Get("Name"))
}

func Test_Mock_GreengrassSubscriptionDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := GreengrassSubscriptionDefinition{ID: ptr.String("def-1")}
	a.Equal("def-1", r.String())
}

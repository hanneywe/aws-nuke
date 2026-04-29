package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
	"github.com/aws/aws-sdk-go-v2/service/applicationsignals/types"
)

func Test_Mock_ApplicationSignalsGroupingConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationSignalsClient)

	mockClient.On("ListGroupingAttributeDefinitions", mock.Anything, mock.Anything).
		Return(&applicationsignals.ListGroupingAttributeDefinitionsOutput{
			GroupingAttributeDefinitions: []types.GroupingAttributeDefinition{
				{
					GroupingName: ptr.String("BusinessUnit"),
				},
				{
					GroupingName: ptr.String("Environment"),
				},
			},
		}, nil)

	lister := &ApplicationSignalsGroupingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*ApplicationSignalsGroupingConfiguration)
	a.Equal("custom (2 definitions)", r.GroupingType)
	a.Equal("BusinessUnit, Environment", r.Definitions)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ApplicationSignalsGroupingConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationSignalsClient)

	mockClient.On("ListGroupingAttributeDefinitions", mock.Anything, mock.Anything).
		Return(&applicationsignals.ListGroupingAttributeDefinitionsOutput{
			GroupingAttributeDefinitions: []types.GroupingAttributeDefinition{},
		}, nil)

	lister := &ApplicationSignalsGroupingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ApplicationSignalsGroupingConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationSignalsClient)

	r := &ApplicationSignalsGroupingConfiguration{
		svc: mockClient,
	}

	mockClient.On("DeleteGroupingConfiguration", mock.Anything,
		&applicationsignals.DeleteGroupingConfigurationInput{}).
		Return(&applicationsignals.DeleteGroupingConfigurationOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ApplicationSignalsGroupingConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := &ApplicationSignalsGroupingConfiguration{
		GroupingType: "custom (2 definitions)",
		Definitions:  "BusinessUnit, Environment",
	}
	props := r.Properties()
	a.Equal("custom (2 definitions)", props.Get("GroupingType"))
	a.Equal("BusinessUnit, Environment", props.Get("Definitions"))
}

func Test_Mock_ApplicationSignalsGroupingConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &ApplicationSignalsGroupingConfiguration{}
	a.Equal("ApplicationSignalsGroupingConfiguration", r.String())
}

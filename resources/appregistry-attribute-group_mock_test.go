package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
	appregistrytypes "github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry/types"
)

func Test_Mock_AppRegistryAttributeGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRegistryClient)
	mockClient.On("ListAttributeGroups", mock.Anything, mock.Anything).
		Return(&servicecatalogappregistry.ListAttributeGroupsOutput{
			AttributeGroups: []appregistrytypes.AttributeGroupSummary{
				{Id: ptr.String("ag-1"), Name: ptr.String("my-group")},
			},
		}, nil)
	lister := &AppRegistryAttributeGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRegistryListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-group", *resources[0].(*AppRegistryAttributeGroup).Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRegistryAttributeGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRegistryClient)
	mockClient.On("ListAttributeGroups", mock.Anything, mock.Anything).
		Return(&servicecatalogappregistry.ListAttributeGroupsOutput{AttributeGroups: []appregistrytypes.AttributeGroupSummary{}}, nil)
	lister := &AppRegistryAttributeGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRegistryListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRegistryAttributeGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRegistryClient)
	r := &AppRegistryAttributeGroup{svc: mockClient, ID: ptr.String("ag-1"), Name: ptr.String("my-group")}
	mockClient.On("DeleteAttributeGroup", mock.Anything, &servicecatalogappregistry.DeleteAttributeGroupInput{AttributeGroup: r.ID}).
		Return(&servicecatalogappregistry.DeleteAttributeGroupOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRegistryAttributeGroup_Properties(t *testing.T) {
	a := assert.New(t)
	r := AppRegistryAttributeGroup{ID: ptr.String("ag-1"), Name: ptr.String("my-group")}
	a.Equal("ag-1", r.Properties().Get("Id"))
	a.Equal("my-group", r.Properties().Get("Name"))
}

func Test_Mock_AppRegistryAttributeGroup_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-group", (&AppRegistryAttributeGroup{Name: ptr.String("my-group")}).String())
}

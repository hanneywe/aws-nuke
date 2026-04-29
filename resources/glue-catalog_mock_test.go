package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueCatalog_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("GetCatalogs", mock.Anything, mock.Anything).
		Return(&glue.GetCatalogsOutput{
			CatalogList: []gluetypes.Catalog{
				{CatalogId: ptr.String("cat-123"), Name: ptr.String("my-catalog")},
			},
		}, nil)
	lister := &GlueCatalogLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-catalog", resources[0].(*GlueCatalog).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCatalog_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("GetCatalogs", mock.Anything, mock.Anything).
		Return(&glue.GetCatalogsOutput{CatalogList: []gluetypes.Catalog{}}, nil)
	lister := &GlueCatalogLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCatalog_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueCatalog{svc: mockClient, CatalogID: ptr.String("cat-123"), Name: ptr.String("my-catalog")}
	mockClient.On("DeleteCatalog", mock.Anything, &glue.DeleteCatalogInput{
		CatalogId: r.CatalogID,
	}).Return(&glue.DeleteCatalogOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCatalog_Filter_Default(t *testing.T) {
	a := assert.New(t)
	r := GlueCatalog{CatalogID: ptr.String("123456789012"), Name: ptr.String("123456789012")}
	a.Error(r.Filter())
}

func Test_Mock_GlueCatalog_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := GlueCatalog{CatalogID: ptr.String("cat-123"), Name: ptr.String("my-catalog")}
	a.NoError(r.Filter())
}

func Test_Mock_GlueCatalog_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueCatalog{CatalogID: ptr.String("cat-123"), Name: ptr.String("my-catalog")}
	a.Equal("cat-123", r.Properties().Get("CatalogId"))
	a.Equal("my-catalog", r.Properties().Get("Name"))
}

func Test_Mock_GlueCatalog_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-catalog", (&GlueCatalog{Name: ptr.String("my-catalog")}).String())
}

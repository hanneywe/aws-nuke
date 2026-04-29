package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"
)

func Test_Mock_IoTSiteWiseAssetModelCompositeModel_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	mockClient.On("ListAssetModels", mock.Anything, mock.Anything).
		Return(&iotsitewise.ListAssetModelsOutput{
			AssetModelSummaries: []iotsitewisetypes.AssetModelSummary{
				{Id: ptr.String("test-assetmodelid")},
			},
		}, nil)

	mockClient.On("ListAssetModelCompositeModels", mock.Anything, mock.Anything).
		Return(&iotsitewise.ListAssetModelCompositeModelsOutput{
			AssetModelCompositeModelSummaries: []iotsitewisetypes.AssetModelCompositeModelSummary{
				{Id: ptr.String("test-assetmodelcompositemodelid"), Name: ptr.String("test-name"), Type: ptr.String("test-type")},
			},
		}, nil)

	lister := &IoTSiteWiseAssetModelCompositeModelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIotsitewiseListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*IoTSiteWiseAssetModelCompositeModel)
	a.Equal("test-assetmodelid", *r.AssetModelID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseAssetModelCompositeModel_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	mockClient.On("ListAssetModels", mock.Anything, mock.Anything).
		Return(&iotsitewise.ListAssetModelsOutput{
			AssetModelSummaries: []iotsitewisetypes.AssetModelSummary{},
		}, nil)

	lister := &IoTSiteWiseAssetModelCompositeModelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIotsitewiseListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseAssetModelCompositeModel_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	r := &IoTSiteWiseAssetModelCompositeModel{
		svc:                        mockClient,
		AssetModelID:               ptr.String("test-assetmodelid"),
		AssetModelCompositeModelID: ptr.String("test-assetmodelcompositemodelid"),
	}

	mockClient.On("DeleteAssetModelCompositeModel", mock.Anything,
		&iotsitewise.DeleteAssetModelCompositeModelInput{
			AssetModelId:               r.AssetModelID,
			AssetModelCompositeModelId: r.AssetModelCompositeModelID,
		}).Return(&iotsitewise.DeleteAssetModelCompositeModelOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseAssetModelCompositeModel_Properties(t *testing.T) {
	a := assert.New(t)
	r := &IoTSiteWiseAssetModelCompositeModel{
		AssetModelID:               ptr.String("test-assetmodelid"),
		AssetModelCompositeModelID: ptr.String("test-assetmodelcompositemodelid"),
		Name:                       ptr.String("test-name"),
		Type:                       ptr.String("test-type"),
	}
	props := r.Properties()
	a.Equal("test-assetmodelid", props.Get("AssetModelId"))
	a.Equal("test-assetmodelcompositemodelid", props.Get("AssetModelCompositeModelId"))
	a.Equal("test-name", props.Get("Name"))
	a.Equal("test-type", props.Get("Type"))
}

func Test_Mock_IoTSiteWiseAssetModelCompositeModel_String(t *testing.T) {
	a := assert.New(t)
	r := &IoTSiteWiseAssetModelCompositeModel{
		AssetModelCompositeModelID: ptr.String("test-assetmodelcompositemodelid"),
	}
	a.Equal("test-assetmodelcompositemodelid", r.String())
}

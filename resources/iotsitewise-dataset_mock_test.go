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

func Test_Mock_IoTSiteWiseDataset_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	mockClient.On("ListDatasets", mock.Anything, mock.Anything).
		Return(&iotsitewise.ListDatasetsOutput{
			DatasetSummaries: []iotsitewisetypes.DatasetSummary{
				{
					Id:   ptr.String("test-dataset-id"),
					Name: ptr.String("test-dataset"),
				},
			},
		}, nil)

	lister := &IoTSiteWiseDatasetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIotsitewiseListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*IoTSiteWiseDataset)
	a.Equal("test-dataset-id", *r.DatasetID)
	a.Equal("test-dataset", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseDataset_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	mockClient.On("ListDatasets", mock.Anything, mock.Anything).
		Return(&iotsitewise.ListDatasetsOutput{
			DatasetSummaries: []iotsitewisetypes.DatasetSummary{},
		}, nil)

	lister := &IoTSiteWiseDatasetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIotsitewiseListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseDataset_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIotsitewiseClient)

	r := &IoTSiteWiseDataset{
		svc:       mockClient,
		DatasetID: ptr.String("test-dataset-id"),
	}

	mockClient.On("DeleteDataset", mock.Anything,
		&iotsitewise.DeleteDatasetInput{
			DatasetId: r.DatasetID,
		}).Return(&iotsitewise.DeleteDatasetOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTSiteWiseDataset_Properties(t *testing.T) {
	a := assert.New(t)
	r := &IoTSiteWiseDataset{
		DatasetID: ptr.String("test-dataset-id"),
		Name:      ptr.String("test-dataset"),
	}
	props := r.Properties()
	a.Equal("test-dataset-id", props.Get("DatasetId"))
	a.Equal("test-dataset", props.Get("Name"))
}

func Test_Mock_IoTSiteWiseDataset_String(t *testing.T) {
	a := assert.New(t)
	r := &IoTSiteWiseDataset{
		DatasetID: ptr.String("test-dataset-id"),
	}
	a.Equal("test-dataset-id", r.String())
}

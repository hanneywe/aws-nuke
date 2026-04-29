package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/dataexchange"
	dataexchangetypes "github.com/aws/aws-sdk-go-v2/service/dataexchange/types"
)

func Test_Mock_DataExchangeDataSet_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	mockClient.On("ListDataSets", mock.Anything, mock.Anything).
		Return(&dataexchange.ListDataSetsOutput{
			DataSets: []dataexchangetypes.DataSetEntry{
				{Id: ptr.String("test-value"), Name: ptr.String("test-value"), Arn: ptr.String("test-value")},
			},
		}, nil)

	lister := &DataExchangeDataSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataexchangeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*DataExchangeDataSet)
	a.Equal("test-value", *r.DataSetID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeDataSet_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	mockClient.On("ListDataSets", mock.Anything, mock.Anything).
		Return(&dataexchange.ListDataSetsOutput{
			DataSets: []dataexchangetypes.DataSetEntry{},
		}, nil)

	lister := &DataExchangeDataSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataexchangeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeDataSet_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	r := &DataExchangeDataSet{
		svc:       mockClient,
		DataSetID: ptr.String("test-datasetid"),
	}

	mockClient.On("DeleteDataSet", mock.Anything,
		&dataexchange.DeleteDataSetInput{
			DataSetId: r.DataSetID,
		}).Return(&dataexchange.DeleteDataSetOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeDataSet_Properties(t *testing.T) {
	a := assert.New(t)
	r := &DataExchangeDataSet{
		DataSetID: ptr.String("test-datasetid"),
		Name:      ptr.String("test-name"),
		ARN:       ptr.String("test-arn"),
	}
	props := r.Properties()
	a.Equal("test-datasetid", props.Get("DataSetId"))
	a.Equal("test-name", props.Get("Name"))
	a.Equal("test-arn", props.Get("Arn"))
}

func Test_Mock_DataExchangeDataSet_String(t *testing.T) {
	a := assert.New(t)
	r := &DataExchangeDataSet{
		DataSetID: ptr.String("test-datasetid"),
	}
	a.Equal("test-datasetid", r.String())
}

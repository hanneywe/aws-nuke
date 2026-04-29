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

func Test_Mock_DataExchangeRevision_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	mockClient.On("ListDataSets", mock.Anything, mock.Anything).
		Return(&dataexchange.ListDataSetsOutput{
			DataSets: []dataexchangetypes.DataSetEntry{
				{Id: ptr.String("test-datasetid")},
			},
		}, nil)

	mockClient.On("ListDataSetRevisions", mock.Anything, mock.Anything).
		Return(&dataexchange.ListDataSetRevisionsOutput{
			Revisions: []dataexchangetypes.RevisionEntry{
				{Id: ptr.String("test-revisionid"), Arn: ptr.String("test-arn"), Comment: ptr.String("test-comment")},
			},
		}, nil)

	lister := &DataExchangeRevisionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataexchangeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*DataExchangeRevision)
	a.Equal("test-datasetid", *r.DataSetID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeRevision_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	mockClient.On("ListDataSets", mock.Anything, mock.Anything).
		Return(&dataexchange.ListDataSetsOutput{
			DataSets: []dataexchangetypes.DataSetEntry{},
		}, nil)

	lister := &DataExchangeRevisionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataexchangeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeRevision_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataexchangeClient)

	r := &DataExchangeRevision{
		svc:        mockClient,
		DataSetID:  ptr.String("test-datasetid"),
		RevisionID: ptr.String("test-revisionid"),
	}

	mockClient.On("DeleteRevision", mock.Anything,
		&dataexchange.DeleteRevisionInput{
			DataSetId:  r.DataSetID,
			RevisionId: r.RevisionID,
		}).Return(&dataexchange.DeleteRevisionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_DataExchangeRevision_Properties(t *testing.T) {
	a := assert.New(t)
	r := &DataExchangeRevision{
		DataSetID:  ptr.String("test-datasetid"),
		RevisionID: ptr.String("test-revisionid"),
		ARN:        ptr.String("test-arn"),
		Comment:    ptr.String("test-comment"),
	}
	props := r.Properties()
	a.Equal("test-datasetid", props.Get("DataSetId"))
	a.Equal("test-revisionid", props.Get("RevisionId"))
	a.Equal("test-arn", props.Get("Arn"))
	a.Equal("test-comment", props.Get("Comment"))
}

func Test_Mock_DataExchangeRevision_String(t *testing.T) {
	a := assert.New(t)
	r := &DataExchangeRevision{
		RevisionID: ptr.String("test-revisionid"),
	}
	a.Equal("test-revisionid", r.String())
}

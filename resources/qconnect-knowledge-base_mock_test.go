package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"
	qconnecttypes "github.com/aws/aws-sdk-go-v2/service/qconnect/types"
)

func Test_Mock_QConnectKnowledgeBase_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListKnowledgeBases", mock.Anything, mock.Anything).
		Return(&qconnect.ListKnowledgeBasesOutput{
			KnowledgeBaseSummaries: []qconnecttypes.KnowledgeBaseSummary{
				{KnowledgeBaseId: ptr.String("kb-12345"), Name: ptr.String("my-kb")},
			},
		}, nil)
	lister := &QConnectKnowledgeBaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	kb := resources[0].(*QConnectKnowledgeBase)
	a.Equal("kb-12345", *kb.KnowledgeBaseID)
	a.Equal("my-kb", *kb.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectKnowledgeBase_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListKnowledgeBases", mock.Anything, mock.Anything).
		Return(&qconnect.ListKnowledgeBasesOutput{KnowledgeBaseSummaries: []qconnecttypes.KnowledgeBaseSummary{}}, nil)
	lister := &QConnectKnowledgeBaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectKnowledgeBase_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	kb := &QConnectKnowledgeBase{
		svc:             mockClient,
		KnowledgeBaseID: ptr.String("kb-12345"),
	}
	mockClient.On("DeleteKnowledgeBase", mock.Anything, &qconnect.DeleteKnowledgeBaseInput{
		KnowledgeBaseId: kb.KnowledgeBaseID,
	}).Return(&qconnect.DeleteKnowledgeBaseOutput{}, nil)
	a.NoError(kb.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectKnowledgeBase_Properties(t *testing.T) {
	a := assert.New(t)
	kb := QConnectKnowledgeBase{
		KnowledgeBaseID: ptr.String("kb-12345"),
		Name:            ptr.String("my-kb"),
	}
	a.Equal("kb-12345", kb.Properties().Get("KnowledgeBaseId"))
	a.Equal("my-kb", kb.Properties().Get("Name"))
}

func Test_Mock_QConnectKnowledgeBase_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("kb-12345", (&QConnectKnowledgeBase{KnowledgeBaseID: ptr.String("kb-12345")}).String())
}

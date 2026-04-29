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

func Test_Mock_QConnectAssistant_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListAssistants", mock.Anything, mock.Anything).
		Return(&qconnect.ListAssistantsOutput{
			AssistantSummaries: []qconnecttypes.AssistantSummary{
				{AssistantId: ptr.String("ast-12345"), Name: ptr.String("my-assistant")},
			},
		}, nil)
	lister := &QConnectAssistantLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	ast := resources[0].(*QConnectAssistant)
	a.Equal("ast-12345", *ast.AssistantID)
	a.Equal("my-assistant", *ast.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAssistant_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListAssistants", mock.Anything, mock.Anything).
		Return(&qconnect.ListAssistantsOutput{AssistantSummaries: []qconnecttypes.AssistantSummary{}}, nil)
	lister := &QConnectAssistantLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAssistant_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	ast := &QConnectAssistant{svc: mockClient, AssistantID: ptr.String("ast-12345")}
	mockClient.On("DeleteAssistant", mock.Anything, &qconnect.DeleteAssistantInput{AssistantId: ast.AssistantID}).
		Return(&qconnect.DeleteAssistantOutput{}, nil)
	a.NoError(ast.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAssistant_Properties(t *testing.T) {
	a := assert.New(t)
	ast := QConnectAssistant{AssistantID: ptr.String("ast-12345"), Name: ptr.String("my-assistant")}
	a.Equal("ast-12345", ast.Properties().Get("AssistantId"))
	a.Equal("my-assistant", ast.Properties().Get("Name"))
}

func Test_Mock_QConnectAssistant_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-assistant", (&QConnectAssistant{Name: ptr.String("my-assistant")}).String())
}

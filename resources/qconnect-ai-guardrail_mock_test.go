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

func Test_Mock_QConnectAIGuardrail_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListAssistants", mock.Anything, mock.Anything).
		Return(&qconnect.ListAssistantsOutput{
			AssistantSummaries: []qconnecttypes.AssistantSummary{
				{AssistantId: ptr.String("ast-12345"), Name: ptr.String("my-assistant")},
			},
		}, nil)
	mockClient.On("ListAIGuardrails", mock.Anything, mock.Anything).
		Return(&qconnect.ListAIGuardrailsOutput{
			AiGuardrailSummaries: []qconnecttypes.AIGuardrailSummary{
				{AiGuardrailId: ptr.String("grd-001"), AssistantId: ptr.String("ast-12345"), Name: ptr.String("my-guardrail")},
			},
		}, nil)
	lister := &QConnectAIGuardrailLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	grd := resources[0].(*QConnectAIGuardrail)
	a.Equal("grd-001", *grd.AIGuardrailID)
	a.Equal("ast-12345", *grd.AssistantID)
	a.Equal("my-guardrail", *grd.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAIGuardrail_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	mockClient.On("ListAssistants", mock.Anything, mock.Anything).
		Return(&qconnect.ListAssistantsOutput{AssistantSummaries: []qconnecttypes.AssistantSummary{}}, nil)
	lister := &QConnectAIGuardrailLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testQConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAIGuardrail_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockQConnectClient)
	grd := &QConnectAIGuardrail{
		svc:           mockClient,
		AssistantID:   ptr.String("ast-12345"),
		AIGuardrailID: ptr.String("grd-001"),
	}
	mockClient.On("DeleteAIGuardrail", mock.Anything, &qconnect.DeleteAIGuardrailInput{
		AiGuardrailId: grd.AIGuardrailID,
		AssistantId:   grd.AssistantID,
	}).Return(&qconnect.DeleteAIGuardrailOutput{}, nil)
	a.NoError(grd.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_QConnectAIGuardrail_Properties(t *testing.T) {
	a := assert.New(t)
	grd := QConnectAIGuardrail{
		AssistantID:   ptr.String("ast-12345"),
		AIGuardrailID: ptr.String("grd-001"),
		Name:          ptr.String("my-guardrail"),
	}
	a.Equal("ast-12345", grd.Properties().Get("AssistantId"))
	a.Equal("grd-001", grd.Properties().Get("AIGuardrailId"))
	a.Equal("my-guardrail", grd.Properties().Get("Name"))
}

func Test_Mock_QConnectAIGuardrail_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("grd-001", (&QConnectAIGuardrail{AIGuardrailID: ptr.String("grd-001")}).String())
}

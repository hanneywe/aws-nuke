package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockQConnectClient struct {
	mock.Mock
}

func (m *mockQConnectClient) ListAssistants(
	ctx context.Context, params *qconnect.ListAssistantsInput,
	_ ...func(*qconnect.Options),
) (*qconnect.ListAssistantsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.ListAssistantsOutput), args.Error(1)
}

func (m *mockQConnectClient) DeleteAssistant(
	ctx context.Context, params *qconnect.DeleteAssistantInput,
	_ ...func(*qconnect.Options),
) (*qconnect.DeleteAssistantOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.DeleteAssistantOutput), args.Error(1)
}

func (m *mockQConnectClient) ListAIGuardrails(
	ctx context.Context, params *qconnect.ListAIGuardrailsInput,
	_ ...func(*qconnect.Options),
) (*qconnect.ListAIGuardrailsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.ListAIGuardrailsOutput), args.Error(1)
}

func (m *mockQConnectClient) DeleteAIGuardrail(
	ctx context.Context, params *qconnect.DeleteAIGuardrailInput,
	_ ...func(*qconnect.Options),
) (*qconnect.DeleteAIGuardrailOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.DeleteAIGuardrailOutput), args.Error(1)
}

func (m *mockQConnectClient) ListKnowledgeBases(
	ctx context.Context, params *qconnect.ListKnowledgeBasesInput,
	_ ...func(*qconnect.Options),
) (*qconnect.ListKnowledgeBasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.ListKnowledgeBasesOutput), args.Error(1)
}

func (m *mockQConnectClient) DeleteKnowledgeBase(
	ctx context.Context, params *qconnect.DeleteKnowledgeBaseInput,
	_ ...func(*qconnect.Options),
) (*qconnect.DeleteKnowledgeBaseOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*qconnect.DeleteKnowledgeBaseOutput), args.Error(1)
}

var testQConnectListerOpts = &nuke.ListerOpts{}

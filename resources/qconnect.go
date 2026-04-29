package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"
)

// QConnectClient is the interface for the QConnect SDK client methods.
type QConnectClient interface {
	ListAssistants(ctx context.Context, params *qconnect.ListAssistantsInput,
		optFns ...func(*qconnect.Options)) (*qconnect.ListAssistantsOutput, error)
	DeleteAssistant(ctx context.Context, params *qconnect.DeleteAssistantInput,
		optFns ...func(*qconnect.Options)) (*qconnect.DeleteAssistantOutput, error)
	ListAIGuardrails(ctx context.Context, params *qconnect.ListAIGuardrailsInput,
		optFns ...func(*qconnect.Options)) (*qconnect.ListAIGuardrailsOutput, error)
	DeleteAIGuardrail(ctx context.Context, params *qconnect.DeleteAIGuardrailInput,
		optFns ...func(*qconnect.Options)) (*qconnect.DeleteAIGuardrailOutput, error)
	ListKnowledgeBases(ctx context.Context, params *qconnect.ListKnowledgeBasesInput,
		optFns ...func(*qconnect.Options)) (*qconnect.ListKnowledgeBasesOutput, error)
	DeleteKnowledgeBase(ctx context.Context, params *qconnect.DeleteKnowledgeBaseInput,
		optFns ...func(*qconnect.Options)) (*qconnect.DeleteKnowledgeBaseOutput, error)
}

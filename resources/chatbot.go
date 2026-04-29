package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/chatbot"
)

// ChatbotClient is the interface for the Chatbot SDK client methods.
type ChatbotClient interface {
	ListCustomActions(ctx context.Context, params *chatbot.ListCustomActionsInput,
		optFns ...func(*chatbot.Options)) (*chatbot.ListCustomActionsOutput, error)
	GetCustomAction(ctx context.Context, params *chatbot.GetCustomActionInput,
		optFns ...func(*chatbot.Options)) (*chatbot.GetCustomActionOutput, error)
	DeleteCustomAction(ctx context.Context, params *chatbot.DeleteCustomActionInput,
		optFns ...func(*chatbot.Options)) (*chatbot.DeleteCustomActionOutput, error)
}

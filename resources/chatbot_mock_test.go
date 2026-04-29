package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/chatbot"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockChatbotClient struct {
	mock.Mock
}

func (m *mockChatbotClient) ListCustomActions(ctx context.Context,
	params *chatbot.ListCustomActionsInput,
	_ ...func(*chatbot.Options)) (*chatbot.ListCustomActionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chatbot.ListCustomActionsOutput), args.Error(1)
}

func (m *mockChatbotClient) GetCustomAction(ctx context.Context,
	params *chatbot.GetCustomActionInput,
	_ ...func(*chatbot.Options)) (*chatbot.GetCustomActionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chatbot.GetCustomActionOutput), args.Error(1)
}

func (m *mockChatbotClient) DeleteCustomAction(ctx context.Context,
	params *chatbot.DeleteCustomActionInput,
	_ ...func(*chatbot.Options)) (*chatbot.DeleteCustomActionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chatbot.DeleteCustomActionOutput), args.Error(1)
}

var testChatbotListerOpts = &nuke.ListerOpts{}

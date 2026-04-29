package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
)

type mockLexV2Client struct {
	mock.Mock
}

func (m *mockLexV2Client) ListBots(
	ctx context.Context, params *lexmodelsv2.ListBotsInput,
	_ ...func(*lexmodelsv2.Options),
) (*lexmodelsv2.ListBotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lexmodelsv2.ListBotsOutput), args.Error(1)
}

func (m *mockLexV2Client) DeleteBot(
	ctx context.Context, params *lexmodelsv2.DeleteBotInput,
	_ ...func(*lexmodelsv2.Options),
) (*lexmodelsv2.DeleteBotOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lexmodelsv2.DeleteBotOutput), args.Error(1)
}

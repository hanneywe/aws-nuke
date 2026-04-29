package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
)

// LexV2Client is the interface for the Lex V2 SDK client methods.
type LexV2Client interface {
	ListBots(ctx context.Context, params *lexmodelsv2.ListBotsInput,
		optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListBotsOutput, error)
	DeleteBot(ctx context.Context, params *lexmodelsv2.DeleteBotInput,
		optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DeleteBotOutput, error)
}

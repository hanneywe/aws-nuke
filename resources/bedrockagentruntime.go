package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
)

// BedrockagentruntimeClient is the interface for the bedrockagentruntime SDK client methods.
type BedrockagentruntimeClient interface {
	ListSessions(ctx context.Context, params *bedrockagentruntime.ListSessionsInput,
		optFns ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.ListSessionsOutput, error)
	EndSession(ctx context.Context, params *bedrockagentruntime.EndSessionInput,
		optFns ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.EndSessionOutput, error)
	DeleteSession(ctx context.Context, params *bedrockagentruntime.DeleteSessionInput,
		optFns ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.DeleteSessionOutput, error)
}

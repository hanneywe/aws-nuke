package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codestarconnections"
)

// CodeStarConnectionsClient is the interface for the CodeStarConnections SDK v2 client methods.
type CodeStarConnectionsClient interface {
	ListHosts(ctx context.Context, params *codestarconnections.ListHostsInput,
		optFns ...func(*codestarconnections.Options)) (*codestarconnections.ListHostsOutput, error)
	DeleteHost(ctx context.Context, params *codestarconnections.DeleteHostInput,
		optFns ...func(*codestarconnections.Options)) (*codestarconnections.DeleteHostOutput, error)
}

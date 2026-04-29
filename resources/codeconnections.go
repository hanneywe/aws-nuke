package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
)

type CodeConnectionsClient interface {
	ListRepositoryLinks(ctx context.Context, params *codeconnections.ListRepositoryLinksInput,
		optFns ...func(*codeconnections.Options)) (*codeconnections.ListRepositoryLinksOutput, error)
	DeleteRepositoryLink(ctx context.Context, params *codeconnections.DeleteRepositoryLinkInput,
		optFns ...func(*codeconnections.Options)) (*codeconnections.DeleteRepositoryLinkOutput, error)
	ListConnections(ctx context.Context, params *codeconnections.ListConnectionsInput,
		optFns ...func(*codeconnections.Options)) (*codeconnections.ListConnectionsOutput, error)
	DeleteConnection(ctx context.Context, params *codeconnections.DeleteConnectionInput,
		optFns ...func(*codeconnections.Options)) (*codeconnections.DeleteConnectionOutput, error)
}

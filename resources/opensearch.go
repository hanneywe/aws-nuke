package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

// OpenSearchClient is the interface for the OpenSearch SDK client methods used by all OpenSearch resources.
// It enables mock testing of List and Remove operations.
type OpenSearchClient interface {
	ListApplications(ctx context.Context, params *opensearch.ListApplicationsInput,
		optFns ...func(*opensearch.Options)) (*opensearch.ListApplicationsOutput, error)
	DeleteApplication(ctx context.Context, params *opensearch.DeleteApplicationInput,
		optFns ...func(*opensearch.Options)) (*opensearch.DeleteApplicationOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dataexchange"
)

// DataexchangeClient is the interface for the dataexchange SDK client methods.
type DataexchangeClient interface {
	ListDataSets(ctx context.Context, params *dataexchange.ListDataSetsInput,
		optFns ...func(*dataexchange.Options)) (*dataexchange.ListDataSetsOutput, error)
	DeleteDataSet(ctx context.Context, params *dataexchange.DeleteDataSetInput,
		optFns ...func(*dataexchange.Options)) (*dataexchange.DeleteDataSetOutput, error)
	ListDataSetRevisions(ctx context.Context, params *dataexchange.ListDataSetRevisionsInput,
		optFns ...func(*dataexchange.Options)) (*dataexchange.ListDataSetRevisionsOutput, error)
	DeleteRevision(ctx context.Context, params *dataexchange.DeleteRevisionInput,
		optFns ...func(*dataexchange.Options)) (*dataexchange.DeleteRevisionOutput, error)
}

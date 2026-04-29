package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
)

type ResourceExplorer2Client interface {
	ListIndexes(ctx context.Context, params *resourceexplorer2.ListIndexesInput,
		optFns ...func(*resourceexplorer2.Options)) (*resourceexplorer2.ListIndexesOutput, error)
	ListTagsForResource(ctx context.Context, params *resourceexplorer2.ListTagsForResourceInput,
		optFns ...func(*resourceexplorer2.Options)) (*resourceexplorer2.ListTagsForResourceOutput, error)
	DeleteIndex(ctx context.Context, params *resourceexplorer2.DeleteIndexInput,
		optFns ...func(*resourceexplorer2.Options)) (*resourceexplorer2.DeleteIndexOutput, error)
}

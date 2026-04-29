package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
)

// IotsitewiseClient is the interface for the iotsitewise SDK client methods.
type IotsitewiseClient interface {
	ListAssetModels(ctx context.Context, params *iotsitewise.ListAssetModelsInput,
		optFns ...func(*iotsitewise.Options)) (*iotsitewise.ListAssetModelsOutput, error)
	ListAssetModelCompositeModels(ctx context.Context, params *iotsitewise.ListAssetModelCompositeModelsInput,
		optFns ...func(*iotsitewise.Options)) (*iotsitewise.ListAssetModelCompositeModelsOutput, error)
	DeleteAssetModelCompositeModel(ctx context.Context, params *iotsitewise.DeleteAssetModelCompositeModelInput,
		optFns ...func(*iotsitewise.Options)) (*iotsitewise.DeleteAssetModelCompositeModelOutput, error)
	ListDatasets(ctx context.Context, params *iotsitewise.ListDatasetsInput,
		optFns ...func(*iotsitewise.Options)) (*iotsitewise.ListDatasetsOutput, error)
	DeleteDataset(ctx context.Context, params *iotsitewise.DeleteDatasetInput,
		optFns ...func(*iotsitewise.Options)) (*iotsitewise.DeleteDatasetOutput, error)
}

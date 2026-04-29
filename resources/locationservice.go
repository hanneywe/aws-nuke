package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/location"
)

// LocationServiceClient is the interface for the Location Service SDK client methods.
type LocationServiceClient interface {
	ListMaps(ctx context.Context, params *location.ListMapsInput,
		optFns ...func(*location.Options)) (*location.ListMapsOutput, error)
	DeleteMap(ctx context.Context, params *location.DeleteMapInput,
		optFns ...func(*location.Options)) (*location.DeleteMapOutput, error)
	ListRouteCalculators(ctx context.Context, params *location.ListRouteCalculatorsInput,
		optFns ...func(*location.Options)) (*location.ListRouteCalculatorsOutput, error)
	DeleteRouteCalculator(ctx context.Context, params *location.DeleteRouteCalculatorInput,
		optFns ...func(*location.Options)) (*location.DeleteRouteCalculatorOutput, error)
	ListKeys(ctx context.Context, params *location.ListKeysInput,
		optFns ...func(*location.Options)) (*location.ListKeysOutput, error)
	DeleteKey(ctx context.Context, params *location.DeleteKeyInput,
		optFns ...func(*location.Options)) (*location.DeleteKeyOutput, error)

	ListPlaceIndexes(ctx context.Context, params *location.ListPlaceIndexesInput,
		optFns ...func(*location.Options)) (*location.ListPlaceIndexesOutput, error)
	DeletePlaceIndex(ctx context.Context, params *location.DeletePlaceIndexInput,
		optFns ...func(*location.Options)) (*location.DeletePlaceIndexOutput, error)
}

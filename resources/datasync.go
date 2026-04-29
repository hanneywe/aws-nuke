package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/datasync"
)

// DataSyncClient is the interface for the DataSync SDK client methods.
type DataSyncClient interface {
	ListLocations(ctx context.Context, params *datasync.ListLocationsInput,
		optFns ...func(*datasync.Options)) (*datasync.ListLocationsOutput, error)
	DeleteLocation(ctx context.Context, params *datasync.DeleteLocationInput,
		optFns ...func(*datasync.Options)) (*datasync.DeleteLocationOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
)

// DMSClient is the interface for the DMS SDK v2 client methods.
type DMSClient interface {
	DescribeDataProviders(ctx context.Context, params *databasemigrationservice.DescribeDataProvidersInput,
		optFns ...func(*databasemigrationservice.Options)) (*databasemigrationservice.DescribeDataProvidersOutput, error)
	DeleteDataProvider(ctx context.Context, params *databasemigrationservice.DeleteDataProviderInput,
		optFns ...func(*databasemigrationservice.Options)) (*databasemigrationservice.DeleteDataProviderOutput, error)
}

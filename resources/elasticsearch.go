package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
)

// ElasticsearchserviceClient is the interface for the elasticsearchservice SDK client methods.
type ElasticsearchserviceClient interface {
	DescribePackages(ctx context.Context, params *elasticsearchservice.DescribePackagesInput,
		optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.DescribePackagesOutput, error)
	DeletePackage(ctx context.Context, params *elasticsearchservice.DeletePackageInput,
		optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.DeletePackageOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
)

// CodeartifactClient is the interface for the codeartifact SDK client methods.
type CodeartifactClient interface {
	ListDomains(ctx context.Context, params *codeartifact.ListDomainsInput,
		optFns ...func(*codeartifact.Options)) (*codeartifact.ListDomainsOutput, error)
	ListRepositoriesInDomain(ctx context.Context, params *codeartifact.ListRepositoriesInDomainInput,
		optFns ...func(*codeartifact.Options)) (*codeartifact.ListRepositoriesInDomainOutput, error)
	ListPackages(ctx context.Context, params *codeartifact.ListPackagesInput,
		optFns ...func(*codeartifact.Options)) (*codeartifact.ListPackagesOutput, error)
	ListPackageVersions(ctx context.Context, params *codeartifact.ListPackageVersionsInput,
		optFns ...func(*codeartifact.Options)) (*codeartifact.ListPackageVersionsOutput, error)
	DeletePackageVersions(ctx context.Context, params *codeartifact.DeletePackageVersionsInput,
		optFns ...func(*codeartifact.Options)) (*codeartifact.DeletePackageVersionsOutput, error)
}

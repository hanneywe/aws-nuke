package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/panorama"
)

// PanoramaClient is the interface for the Panorama SDK client methods.
type PanoramaClient interface {
	ListPackages(ctx context.Context, params *panorama.ListPackagesInput,
		optFns ...func(*panorama.Options)) (*panorama.ListPackagesOutput, error)
	DeletePackage(ctx context.Context, params *panorama.DeletePackageInput,
		optFns ...func(*panorama.Options)) (*panorama.DeletePackageOutput, error)
}

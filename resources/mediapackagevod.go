package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagevod"
)

// MediaPackageVODClient is the interface for the MediaPackage VOD SDK client methods.
type MediaPackageVODClient interface {
	ListPackagingGroups(ctx context.Context, params *mediapackagevod.ListPackagingGroupsInput,
		optFns ...func(*mediapackagevod.Options)) (*mediapackagevod.ListPackagingGroupsOutput, error)
	DeletePackagingGroup(ctx context.Context, params *mediapackagevod.DeletePackagingGroupInput,
		optFns ...func(*mediapackagevod.Options)) (*mediapackagevod.DeletePackagingGroupOutput, error)
}

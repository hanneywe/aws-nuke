package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
)

// DeviceFarmClient is an interface for the DeviceFarm SDK client methods used by all DeviceFarm resources.
// It enables mock testing of List and Remove operations.
type DeviceFarmClient interface {
	// Listing
	ListInstanceProfiles(ctx context.Context, params *devicefarm.ListInstanceProfilesInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.ListInstanceProfilesOutput, error)
	ListTestGridProjects(ctx context.Context, params *devicefarm.ListTestGridProjectsInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.ListTestGridProjectsOutput, error)
	ListProjects(ctx context.Context, params *devicefarm.ListProjectsInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.ListProjectsOutput, error)
	ListUploads(ctx context.Context, params *devicefarm.ListUploadsInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.ListUploadsOutput, error)

	// Deletion
	DeleteInstanceProfile(ctx context.Context, params *devicefarm.DeleteInstanceProfileInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.DeleteInstanceProfileOutput, error)
	DeleteTestGridProject(ctx context.Context, params *devicefarm.DeleteTestGridProjectInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.DeleteTestGridProjectOutput, error)
	DeleteUpload(ctx context.Context, params *devicefarm.DeleteUploadInput,
		optFns ...func(*devicefarm.Options)) (*devicefarm.DeleteUploadOutput, error)
}

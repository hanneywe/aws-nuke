package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
)

// LightsailClient is the interface for the Lightsail SDK v2 client methods used by new Lightsail resources.
// Existing Lightsail resources use SDK v1; this interface is for new SDK v2 resources only.
type LightsailClient interface {
	GetBuckets(ctx context.Context, params *lightsail.GetBucketsInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetBucketsOutput, error)
	DeleteBucket(ctx context.Context, params *lightsail.DeleteBucketInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteBucketOutput, error)
	GetCertificates(ctx context.Context, params *lightsail.GetCertificatesInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetCertificatesOutput, error)
	DeleteCertificate(ctx context.Context, params *lightsail.DeleteCertificateInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteCertificateOutput, error)
	GetContainerServices(ctx context.Context, params *lightsail.GetContainerServicesInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetContainerServicesOutput, error)
	DeleteContainerService(ctx context.Context, params *lightsail.DeleteContainerServiceInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteContainerServiceOutput, error)
	GetContactMethods(ctx context.Context, params *lightsail.GetContactMethodsInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetContactMethodsOutput, error)
	DeleteContactMethod(ctx context.Context, params *lightsail.DeleteContactMethodInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteContactMethodOutput, error)

	GetBucketAccessKeys(ctx context.Context, params *lightsail.GetBucketAccessKeysInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetBucketAccessKeysOutput, error)
	DeleteBucketAccessKey(ctx context.Context, params *lightsail.DeleteBucketAccessKeyInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteBucketAccessKeyOutput, error)

	GetDiskSnapshots(ctx context.Context, params *lightsail.GetDiskSnapshotsInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetDiskSnapshotsOutput, error)
	DeleteDiskSnapshot(ctx context.Context, params *lightsail.DeleteDiskSnapshotInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteDiskSnapshotOutput, error)

	GetRelationalDatabases(ctx context.Context, params *lightsail.GetRelationalDatabasesInput,
		optFns ...func(*lightsail.Options)) (*lightsail.GetRelationalDatabasesOutput, error)
	DeleteRelationalDatabase(ctx context.Context, params *lightsail.DeleteRelationalDatabaseInput,
		optFns ...func(*lightsail.Options)) (*lightsail.DeleteRelationalDatabaseOutput, error)
}

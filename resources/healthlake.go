package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/healthlake"
)

// HealthLakeClient is an interface for the AWS HealthLake SDK client methods used by all HealthLake resources.
// It enables mock testing of List and Remove operations.
type HealthLakeClient interface {
	ListFHIRDatastores(ctx context.Context, params *healthlake.ListFHIRDatastoresInput,
		optFns ...func(*healthlake.Options)) (*healthlake.ListFHIRDatastoresOutput, error)
	DeleteFHIRDatastore(ctx context.Context, params *healthlake.DeleteFHIRDatastoreInput,
		optFns ...func(*healthlake.Options)) (*healthlake.DeleteFHIRDatastoreOutput, error)
}

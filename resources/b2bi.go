package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/b2bi"
)

// B2BIClient is an interface for the B2BI SDK client methods used by all B2BI resources.
// It enables mock testing of List and Remove operations.
type B2BIClient interface {
	ListProfiles(ctx context.Context, params *b2bi.ListProfilesInput,
		optFns ...func(*b2bi.Options)) (*b2bi.ListProfilesOutput, error)
	DeleteProfile(ctx context.Context, params *b2bi.DeleteProfileInput,
		optFns ...func(*b2bi.Options)) (*b2bi.DeleteProfileOutput, error)
}

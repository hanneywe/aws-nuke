package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/signer"
)

// SignerClient is an interface for the AWS Signer SDK v2 client methods.
type SignerClient interface {
	ListSigningProfiles(ctx context.Context, params *signer.ListSigningProfilesInput,
		optFns ...func(*signer.Options)) (*signer.ListSigningProfilesOutput, error)
	CancelSigningProfile(ctx context.Context, params *signer.CancelSigningProfileInput,
		optFns ...func(*signer.Options)) (*signer.CancelSigningProfileOutput, error)
}

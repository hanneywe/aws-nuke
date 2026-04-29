package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"
)

type CustomerProfilesClient interface {
	ListDomains(ctx context.Context, params *customerprofiles.ListDomainsInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.ListDomainsOutput, error)
	DeleteDomain(ctx context.Context, params *customerprofiles.DeleteDomainInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.DeleteDomainOutput, error)
	ListSegmentDefinitions(ctx context.Context, params *customerprofiles.ListSegmentDefinitionsInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.ListSegmentDefinitionsOutput, error)
	DeleteSegmentDefinition(ctx context.Context, params *customerprofiles.DeleteSegmentDefinitionInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.DeleteSegmentDefinitionOutput, error)
	SearchProfiles(ctx context.Context, params *customerprofiles.SearchProfilesInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.SearchProfilesOutput, error)
	DeleteProfile(ctx context.Context, params *customerprofiles.DeleteProfileInput,
		optFns ...func(*customerprofiles.Options)) (*customerprofiles.DeleteProfileOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// Route53Client is an interface for the Route53 SDK client methods used by all Route53 resources.
// It enables mock testing of List and Remove operations.
type Route53Client interface {
	// Listing
	ListCidrCollections(ctx context.Context, params *route53.ListCidrCollectionsInput,
		optFns ...func(*route53.Options)) (*route53.ListCidrCollectionsOutput, error)

	// Deletion
	DeleteCidrCollection(ctx context.Context, params *route53.DeleteCidrCollectionInput,
		optFns ...func(*route53.Options)) (*route53.DeleteCidrCollectionOutput, error)
	ListReusableDelegationSets(ctx context.Context, params *route53.ListReusableDelegationSetsInput,
		optFns ...func(*route53.Options)) (*route53.ListReusableDelegationSetsOutput, error)
	DeleteReusableDelegationSet(ctx context.Context, params *route53.DeleteReusableDelegationSetInput,
		optFns ...func(*route53.Options)) (*route53.DeleteReusableDelegationSetOutput, error)
}

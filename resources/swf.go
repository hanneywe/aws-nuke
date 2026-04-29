package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/swf"
)

// SWFClient is the interface for the SWF SDK client methods.
type SWFClient interface {
	ListDomains(ctx context.Context, params *swf.ListDomainsInput,
		optFns ...func(*swf.Options)) (*swf.ListDomainsOutput, error)
	DeprecateDomain(ctx context.Context, params *swf.DeprecateDomainInput,
		optFns ...func(*swf.Options)) (*swf.DeprecateDomainOutput, error)
	ListActivityTypes(ctx context.Context, params *swf.ListActivityTypesInput,
		optFns ...func(*swf.Options)) (*swf.ListActivityTypesOutput, error)
	DeprecateActivityType(ctx context.Context, params *swf.DeprecateActivityTypeInput,
		optFns ...func(*swf.Options)) (*swf.DeprecateActivityTypeOutput, error)
}

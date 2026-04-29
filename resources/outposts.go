package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/outposts"
)

// OutpostsClient is the interface for the Outposts SDK client methods.
type OutpostsClient interface {
	ListSites(ctx context.Context, params *outposts.ListSitesInput,
		optFns ...func(*outposts.Options)) (*outposts.ListSitesOutput, error)
	DeleteSite(ctx context.Context, params *outposts.DeleteSiteInput,
		optFns ...func(*outposts.Options)) (*outposts.DeleteSiteOutput, error)
}

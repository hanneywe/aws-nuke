package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"
)

// ResilienceHubClient is the interface for the Resilience Hub SDK client methods.
type ResilienceHubClient interface {
	ListApps(ctx context.Context, params *resiliencehub.ListAppsInput,
		optFns ...func(*resiliencehub.Options)) (*resiliencehub.ListAppsOutput, error)
	DeleteApp(ctx context.Context, params *resiliencehub.DeleteAppInput,
		optFns ...func(*resiliencehub.Options)) (*resiliencehub.DeleteAppOutput, error)
}

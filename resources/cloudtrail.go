package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

// CloudTrailClient is the interface for the CloudTrail SDK v2 client methods.
type CloudTrailClient interface {
	ListDashboards(ctx context.Context, params *cloudtrail.ListDashboardsInput,
		optFns ...func(*cloudtrail.Options)) (*cloudtrail.ListDashboardsOutput, error)
	GetDashboard(ctx context.Context, params *cloudtrail.GetDashboardInput,
		optFns ...func(*cloudtrail.Options)) (*cloudtrail.GetDashboardOutput, error)
	UpdateDashboard(ctx context.Context, params *cloudtrail.UpdateDashboardInput,
		optFns ...func(*cloudtrail.Options)) (*cloudtrail.UpdateDashboardOutput, error)
	DeleteDashboard(ctx context.Context, params *cloudtrail.DeleteDashboardInput,
		optFns ...func(*cloudtrail.Options)) (*cloudtrail.DeleteDashboardOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudWatchLogsIntegrationResource = "CloudWatchLogsIntegration"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudWatchLogsIntegrationResource,
		Scope:    nuke.Account,
		Resource: &CloudWatchLogsIntegration{},
		Lister:   &CloudWatchLogsIntegrationLister{},
	})
}

type CloudWatchLogsIntegrationLister struct {
	svc CloudWatchLogsV2Client
}

func (l *CloudWatchLogsIntegrationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudwatchlogs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	output, err := svc.ListIntegrations(ctx, &cloudwatchlogs.ListIntegrationsInput{})
	if err != nil {
		return nil, err
	}

	for _, integration := range output.IntegrationSummaries {
		resources = append(resources, &CloudWatchLogsIntegration{
			svc:             svc,
			IntegrationName: integration.IntegrationName,
			IntegrationType: integration.IntegrationType,
		})
	}

	return resources, nil
}

type CloudWatchLogsIntegration struct {
	svc             CloudWatchLogsV2Client
	IntegrationName *string
	IntegrationType cloudwatchlogstypes.IntegrationType
}

func (r *CloudWatchLogsIntegration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIntegration(ctx, &cloudwatchlogs.DeleteIntegrationInput{
		IntegrationName: r.IntegrationName,
		Force:           true,
	})
	return err
}

func (r *CloudWatchLogsIntegration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudWatchLogsIntegration) String() string {
	return *r.IntegrationName
}

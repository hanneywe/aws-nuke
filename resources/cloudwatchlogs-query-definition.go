package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudWatchLogsQueryDefinitionResource = "CloudWatchLogsQueryDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudWatchLogsQueryDefinitionResource,
		Scope:    nuke.Account,
		Resource: &CloudWatchLogsQueryDefinition{},
		Lister:   &CloudWatchLogsQueryDefinitionLister{},
	})
}

type CloudWatchLogsQueryDefinitionLister struct {
	svc CloudWatchLogsV2Client
}

func (l *CloudWatchLogsQueryDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudwatchlogs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &cloudwatchlogs.DescribeQueryDefinitionsInput{}
	for {
		output, err := svc.DescribeQueryDefinitions(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, queryDefinition := range output.QueryDefinitions {
			resources = append(resources, &CloudWatchLogsQueryDefinition{
				svc:               svc,
				QueryDefinitionID: queryDefinition.QueryDefinitionId,
				Name:              queryDefinition.Name,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type CloudWatchLogsQueryDefinition struct {
	svc               CloudWatchLogsV2Client
	QueryDefinitionID *string `property:"name=QueryDefinitionId"`
	Name              *string
}

func (r *CloudWatchLogsQueryDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteQueryDefinition(ctx, &cloudwatchlogs.DeleteQueryDefinitionInput{
		QueryDefinitionId: r.QueryDefinitionID,
	})
	return err
}

func (r *CloudWatchLogsQueryDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudWatchLogsQueryDefinition) String() string {
	return *r.Name
}

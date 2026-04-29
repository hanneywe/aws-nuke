package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudWatchLogsBearerTokenAuthenticationResource = "CloudWatchLogsBearerTokenAuthentication"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudWatchLogsBearerTokenAuthenticationResource,
		Scope:    nuke.Account,
		Resource: &CloudWatchLogsBearerTokenAuthentication{},
		Lister:   &CloudWatchLogsBearerTokenAuthenticationLister{},
	})
}

type CloudWatchLogsBearerTokenAuthenticationLister struct {
	svc CloudWatchLogsV2Client
}

func (l *CloudWatchLogsBearerTokenAuthenticationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudwatchlogs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(svc, &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int32(50),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range page.LogGroups {
			logGroup := &page.LogGroups[i]
			if logGroup.BearerTokenAuthenticationEnabled != nil && *logGroup.BearerTokenAuthenticationEnabled {
				resources = append(resources, &CloudWatchLogsBearerTokenAuthentication{
					svc:          svc,
					LogGroupName: logGroup.LogGroupName,
				})
			}
		}
	}

	return resources, nil
}

type CloudWatchLogsBearerTokenAuthentication struct {
	svc          CloudWatchLogsV2Client
	LogGroupName *string
}

func (r *CloudWatchLogsBearerTokenAuthentication) Remove(ctx context.Context) error {
	_, err := r.svc.PutBearerTokenAuthentication(ctx, &cloudwatchlogs.PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               r.LogGroupName,
		BearerTokenAuthenticationEnabled: aws.Bool(false),
	})
	return err
}

func (r *CloudWatchLogsBearerTokenAuthentication) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudWatchLogsBearerTokenAuthentication) String() string {
	return *r.LogGroupName
}

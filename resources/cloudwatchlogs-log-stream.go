package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudWatchLogsLogStreamResource = "CloudWatchLogsLogStream"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudWatchLogsLogStreamResource,
		Scope:    nuke.Account,
		Resource: &CloudWatchLogsLogStream{},
		Lister:   &CloudWatchLogsLogStreamLister{},
		DependsOn: []string{
			CloudWatchLogsLogGroupResource,
		},
	})
}

type CloudWatchLogsLogStreamLister struct {
	svc CloudWatchLogsV2Client
}

func (l *CloudWatchLogsLogStreamLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudwatchlogs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	groupPaginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(svc, &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int32(50),
	})

	for groupPaginator.HasMorePages() {
		groupPage, err := groupPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range groupPage.LogGroups {
			logGroup := &groupPage.LogGroups[i]
			streamPaginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(svc, &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: logGroup.LogGroupName,
				Limit:        aws.Int32(50),
			})

			for streamPaginator.HasMorePages() {
				streamPage, err := streamPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for j := range streamPage.LogStreams {
					stream := &streamPage.LogStreams[j]
					var lastEventTimestamp *time.Time
					if stream.LastEventTimestamp != nil {
						t := time.Unix(*stream.LastEventTimestamp/1000, 0).UTC()
						lastEventTimestamp = &t
					}

					resources = append(resources, &CloudWatchLogsLogStream{
						svc:                svc,
						LogGroupName:       logGroup.LogGroupName,
						LogStreamName:      stream.LogStreamName,
						LastEventTimestamp: lastEventTimestamp,
					})
				}
			}
		}
	}

	return resources, nil
}

type CloudWatchLogsLogStream struct {
	svc                CloudWatchLogsV2Client
	LogGroupName       *string
	LogStreamName      *string
	LastEventTimestamp *time.Time
}

func (r *CloudWatchLogsLogStream) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLogStream(ctx, &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  r.LogGroupName,
		LogStreamName: r.LogStreamName,
	})
	return err
}

func (r *CloudWatchLogsLogStream) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudWatchLogsLogStream) String() string {
	return fmt.Sprintf("%s/%s", *r.LogGroupName, *r.LogStreamName)
}

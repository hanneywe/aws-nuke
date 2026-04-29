package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaLiveCloudWatchAlarmTemplateGroupResource = "MediaLiveCloudWatchAlarmTemplateGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveCloudWatchAlarmTemplateGroupResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveCloudWatchAlarmTemplateGroup{},
		Lister:   &MediaLiveCloudWatchAlarmTemplateGroupLister{},
		DependsOn: []string{
			MediaLiveCloudWatchAlarmTemplateResource,
		},
	})
}

type MediaLiveCloudWatchAlarmTemplateGroupLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveCloudWatchAlarmTemplateGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := medialive.NewListCloudWatchAlarmTemplateGroupsPaginator(svc, &medialive.ListCloudWatchAlarmTemplateGroupsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.CloudWatchAlarmTemplateGroups {
			resources = append(resources, &MediaLiveCloudWatchAlarmTemplateGroup{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type MediaLiveCloudWatchAlarmTemplateGroup struct {
	svc  MediaLiveClient
	ID   *string
	Name *string
}

func (r *MediaLiveCloudWatchAlarmTemplateGroup) Filter() error {
	if r.Name != nil && strings.HasPrefix(*r.Name, "AWS-") {
		return fmt.Errorf("cannot delete AWS-managed template group")
	}
	return nil
}

func (r *MediaLiveCloudWatchAlarmTemplateGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCloudWatchAlarmTemplateGroup(ctx, &medialive.DeleteCloudWatchAlarmTemplateGroupInput{
		Identifier: r.ID,
	})
	return err
}

func (r *MediaLiveCloudWatchAlarmTemplateGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveCloudWatchAlarmTemplateGroup) String() string {
	return *r.ID
}

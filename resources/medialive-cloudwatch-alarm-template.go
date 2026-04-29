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

const MediaLiveCloudWatchAlarmTemplateResource = "MediaLiveCloudWatchAlarmTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveCloudWatchAlarmTemplateResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveCloudWatchAlarmTemplate{},
		Lister:   &MediaLiveCloudWatchAlarmTemplateLister{},
	})
}

type MediaLiveCloudWatchAlarmTemplateLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveCloudWatchAlarmTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := medialive.NewListCloudWatchAlarmTemplatesPaginator(svc, &medialive.ListCloudWatchAlarmTemplatesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.CloudWatchAlarmTemplates {
			item := &resp.CloudWatchAlarmTemplates[i]
			resources = append(resources, &MediaLiveCloudWatchAlarmTemplate{
				svc:     svc,
				ID:      item.Id,
				Name:    item.Name,
				GroupID: item.GroupId,
			})
		}
	}

	return resources, nil
}

type MediaLiveCloudWatchAlarmTemplate struct {
	svc     MediaLiveClient
	ID      *string
	Name    *string
	GroupID *string
}

func (r *MediaLiveCloudWatchAlarmTemplate) Filter() error {
	if r.Name != nil && strings.HasPrefix(*r.Name, "AWS-") {
		return fmt.Errorf("cannot delete AWS-managed alarm template")
	}
	return nil
}

func (r *MediaLiveCloudWatchAlarmTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCloudWatchAlarmTemplate(ctx, &medialive.DeleteCloudWatchAlarmTemplateInput{
		Identifier: r.ID,
	})
	return err
}

func (r *MediaLiveCloudWatchAlarmTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveCloudWatchAlarmTemplate) String() string {
	return *r.ID
}

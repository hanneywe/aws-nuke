package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaLiveEventBridgeRuleTemplateResource = "MediaLiveEventBridgeRuleTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveEventBridgeRuleTemplateResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveEventBridgeRuleTemplate{},
		Lister:   &MediaLiveEventBridgeRuleTemplateLister{},
	})
}

type MediaLiveEventBridgeRuleTemplateLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveEventBridgeRuleTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := medialive.NewListEventBridgeRuleTemplatesPaginator(svc, &medialive.ListEventBridgeRuleTemplatesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.EventBridgeRuleTemplates {
			resources = append(resources, &MediaLiveEventBridgeRuleTemplate{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type MediaLiveEventBridgeRuleTemplate struct {
	svc  MediaLiveClient
	ID   *string
	Name *string
}

func (r *MediaLiveEventBridgeRuleTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventBridgeRuleTemplate(ctx, &medialive.DeleteEventBridgeRuleTemplateInput{
		Identifier: r.ID,
	})
	return err
}

func (r *MediaLiveEventBridgeRuleTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveEventBridgeRuleTemplate) String() string {
	return *r.ID
}

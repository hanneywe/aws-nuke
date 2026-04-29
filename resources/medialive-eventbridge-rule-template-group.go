package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/medialive"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaLiveEventBridgeRuleTemplateGroupResource = "MediaLiveEventBridgeRuleTemplateGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaLiveEventBridgeRuleTemplateGroupResource,
		Scope:    nuke.Account,
		Resource: &MediaLiveEventBridgeRuleTemplateGroup{},
		Lister:   &MediaLiveEventBridgeRuleTemplateGroupLister{},
		DependsOn: []string{
			MediaLiveEventBridgeRuleTemplateResource,
		},
	})
}

type MediaLiveEventBridgeRuleTemplateGroupLister struct {
	svc MediaLiveClient
}

func (l *MediaLiveEventBridgeRuleTemplateGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = medialive.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &medialive.ListEventBridgeRuleTemplateGroupsInput{}
	for {
		resp, err := svc.ListEventBridgeRuleTemplateGroups(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.EventBridgeRuleTemplateGroups {
			resources = append(resources, &MediaLiveEventBridgeRuleTemplateGroup{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type MediaLiveEventBridgeRuleTemplateGroup struct {
	svc  MediaLiveClient
	ID   *string
	Name *string
}

func (r *MediaLiveEventBridgeRuleTemplateGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventBridgeRuleTemplateGroup(ctx, &medialive.DeleteEventBridgeRuleTemplateGroupInput{
		Identifier: r.ID,
	})
	return err
}

func (r *MediaLiveEventBridgeRuleTemplateGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaLiveEventBridgeRuleTemplateGroup) String() string {
	return *r.Name
}

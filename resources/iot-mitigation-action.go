package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTMitigationActionResource = "IoTMitigationAction"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTMitigationActionResource,
		Scope:    nuke.Account,
		Resource: &IoTMitigationAction{},
		Lister:   &IoTMitigationActionLister{},
	})
}

type IoTMitigationActionLister struct {
	svc IoTClient
}

func (l *IoTMitigationActionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iot.NewListMitigationActionsPaginator(svc, &iot.ListMitigationActionsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, action := range resp.ActionIdentifiers {
			resources = append(resources, &IoTMitigationAction{
				svc:        svc,
				ActionName: action.ActionName,
				ActionArn:  action.ActionArn,
			})
		}
	}

	return resources, nil
}

type IoTMitigationAction struct {
	svc        IoTClient
	ActionName *string
	ActionArn  *string
}

func (r *IoTMitigationAction) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMitigationAction(ctx, &iot.DeleteMitigationActionInput{
		ActionName: r.ActionName,
	})
	return err
}

func (r *IoTMitigationAction) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTMitigationAction) String() string {
	return *r.ActionName
}

package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const (
	PinpointSMSVoiceV2OptOutListResource    = "PinpointSMSVoiceV2OptOutList"
	PinpointSMSVoiceV2OptOutListNameDefault = "Default"
)

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2OptOutListResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2OptOutList{},
		Lister:   &PinpointSMSVoiceV2OptOutListLister{},
	})
}

type PinpointSMSVoiceV2OptOutListLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2OptOutListLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = pinpointsmsvoicev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := pinpointsmsvoicev2.NewDescribeOptOutListsPaginator(svc, &pinpointsmsvoicev2.DescribeOptOutListsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.OptOutLists {
			resources = append(resources, &PinpointSMSVoiceV2OptOutList{
				svc:            svc,
				OptOutListName: item.OptOutListName,
				OptOutListArn:  item.OptOutListArn,
			})
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2OptOutList struct {
	svc            PinpointSMSVoiceV2Client
	OptOutListName *string
	OptOutListArn  *string
}

func (r *PinpointSMSVoiceV2OptOutList) Filter() error {
	if r.OptOutListName != nil && *r.OptOutListName == PinpointSMSVoiceV2OptOutListNameDefault {
		return fmt.Errorf("cannot delete default opt-out list")
	}
	return nil
}

func (r *PinpointSMSVoiceV2OptOutList) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteOptOutList(ctx, &pinpointsmsvoicev2.DeleteOptOutListInput{
		OptOutListName: r.OptOutListName,
	})
	return err
}

func (r *PinpointSMSVoiceV2OptOutList) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2OptOutList) String() string {
	return *r.OptOutListName
}

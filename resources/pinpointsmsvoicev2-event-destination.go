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

const PinpointSMSVoiceV2EventDestinationResource = "PinpointSMSVoiceV2EventDestination"

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2EventDestinationResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2EventDestination{},
		Lister:   &PinpointSMSVoiceV2EventDestinationLister{},
	})
}

type PinpointSMSVoiceV2EventDestinationLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2EventDestinationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = pinpointsmsvoicev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := pinpointsmsvoicev2.NewDescribeConfigurationSetsPaginator(svc, &pinpointsmsvoicev2.DescribeConfigurationSetsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cs := range resp.ConfigurationSets {
			for _, ed := range cs.EventDestinations {
				resources = append(resources, &PinpointSMSVoiceV2EventDestination{
					svc:                  svc,
					EventDestinationName: ed.EventDestinationName,
					ConfigurationSetName: cs.ConfigurationSetName,
				})
			}
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2EventDestination struct {
	svc                  PinpointSMSVoiceV2Client
	EventDestinationName *string
	ConfigurationSetName *string
}

func (r *PinpointSMSVoiceV2EventDestination) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventDestination(ctx, &pinpointsmsvoicev2.DeleteEventDestinationInput{
		ConfigurationSetName: r.ConfigurationSetName,
		EventDestinationName: r.EventDestinationName,
	})
	return err
}

func (r *PinpointSMSVoiceV2EventDestination) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2EventDestination) String() string {
	return fmt.Sprintf("%s -> %s", *r.ConfigurationSetName, *r.EventDestinationName)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PinpointSMSVoiceV2ConfigurationSetResource = "PinpointSMSVoiceV2ConfigurationSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2ConfigurationSetResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2ConfigurationSet{},
		Lister:   &PinpointSMSVoiceV2ConfigurationSetLister{},
	})
}

type PinpointSMSVoiceV2ConfigurationSetLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2ConfigurationSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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

		for _, item := range resp.ConfigurationSets {
			resources = append(resources, &PinpointSMSVoiceV2ConfigurationSet{
				svc:                  svc,
				ConfigurationSetName: item.ConfigurationSetName,
				ConfigurationSetArn:  item.ConfigurationSetArn,
			})
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2ConfigurationSet struct {
	svc                  PinpointSMSVoiceV2Client
	ConfigurationSetName *string
	ConfigurationSetArn  *string
}

func (r *PinpointSMSVoiceV2ConfigurationSet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConfigurationSet(ctx, &pinpointsmsvoicev2.DeleteConfigurationSetInput{
		ConfigurationSetName: r.ConfigurationSetName,
	})
	return err
}

func (r *PinpointSMSVoiceV2ConfigurationSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2ConfigurationSet) String() string {
	return *r.ConfigurationSetName
}

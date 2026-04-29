package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PinpointSMSVoiceV2ProtectConfigurationResource = "PinpointSMSVoiceV2ProtectConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2ProtectConfigurationResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2ProtectConfiguration{},
		Lister:   &PinpointSMSVoiceV2ProtectConfigurationLister{},
		Settings: []string{
			"DisableDeletionProtection",
		},
	})
}

type PinpointSMSVoiceV2ProtectConfigurationLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2ProtectConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = pinpointsmsvoicev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := pinpointsmsvoicev2.NewDescribeProtectConfigurationsPaginator(svc, &pinpointsmsvoicev2.DescribeProtectConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ProtectConfigurations {
			resources = append(resources, &PinpointSMSVoiceV2ProtectConfiguration{
				svc:                       svc,
				ProtectConfigurationID:    item.ProtectConfigurationId,
				ProtectConfigurationArn:   item.ProtectConfigurationArn,
				DeletionProtectionEnabled: &item.DeletionProtectionEnabled,
			})
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2ProtectConfiguration struct {
	svc                       PinpointSMSVoiceV2Client
	settings                  *libsettings.Setting
	ProtectConfigurationID    *string `property:"name=ProtectConfigurationId"`
	ProtectConfigurationArn   *string
	DeletionProtectionEnabled *bool
}

func (r *PinpointSMSVoiceV2ProtectConfiguration) Settings(setting *libsettings.Setting) {
	r.settings = setting
}

func (r *PinpointSMSVoiceV2ProtectConfiguration) Remove(ctx context.Context) error {
	if r.DeletionProtectionEnabled != nil && *r.DeletionProtectionEnabled && r.settings.GetBool("DisableDeletionProtection") {
		_, err := r.svc.UpdateProtectConfiguration(ctx, &pinpointsmsvoicev2.UpdateProtectConfigurationInput{
			ProtectConfigurationId:    r.ProtectConfigurationID,
			DeletionProtectionEnabled: aws.Bool(false),
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteProtectConfiguration(ctx, &pinpointsmsvoicev2.DeleteProtectConfigurationInput{
		ProtectConfigurationId: r.ProtectConfigurationID,
	})
	return err
}

func (r *PinpointSMSVoiceV2ProtectConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2ProtectConfiguration) String() string {
	return *r.ProtectConfigurationID
}

package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PinpointSMSVoiceV2ProtectConfigRuleSetOverrideResource = "PinpointSMSVoiceV2ProtectConfigRuleSetOverride"

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2ProtectConfigRuleSetOverrideResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2ProtectConfigRuleSetOverride{},
		Lister:   &PinpointSMSVoiceV2ProtectConfigRuleSetOverrideLister{},
	})
}

type PinpointSMSVoiceV2ProtectConfigRuleSetOverrideLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2ProtectConfigRuleSetOverrideLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = pinpointsmsvoicev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	pcPaginator := pinpointsmsvoicev2.NewDescribeProtectConfigurationsPaginator(svc, &pinpointsmsvoicev2.DescribeProtectConfigurationsInput{})

	for pcPaginator.HasMorePages() {
		pcResp, err := pcPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, pc := range pcResp.ProtectConfigurations {
			overridePaginator := pinpointsmsvoicev2.NewListProtectConfigurationRuleSetNumberOverridesPaginator(svc,
				&pinpointsmsvoicev2.ListProtectConfigurationRuleSetNumberOverridesInput{
					ProtectConfigurationId: pc.ProtectConfigurationId,
				})

			for overridePaginator.HasMorePages() {
				overrideResp, err := overridePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, override := range overrideResp.RuleSetNumberOverrides {
					resources = append(resources, &PinpointSMSVoiceV2ProtectConfigRuleSetOverride{
						svc:                    svc,
						ProtectConfigurationID: pc.ProtectConfigurationId,
						DestinationPhoneNumber: override.DestinationPhoneNumber,
						Action:                 override.Action,
						IsoCountryCode:         override.IsoCountryCode,
						CreatedTimestamp:       override.CreatedTimestamp,
						ExpirationTimestamp:    override.ExpirationTimestamp,
					})
				}
			}
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2ProtectConfigRuleSetOverride struct {
	svc                    PinpointSMSVoiceV2Client
	ProtectConfigurationID *string `property:"name=ProtectConfigurationId"`
	DestinationPhoneNumber *string
	Action                 pinpointtypes.ProtectConfigurationRuleOverrideAction `property:"name=Action"`
	IsoCountryCode         *string
	CreatedTimestamp       *time.Time
	ExpirationTimestamp    *time.Time
}

func (r *PinpointSMSVoiceV2ProtectConfigRuleSetOverride) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProtectConfigurationRuleSetNumberOverride(ctx,
		&pinpointsmsvoicev2.DeleteProtectConfigurationRuleSetNumberOverrideInput{
			ProtectConfigurationId: r.ProtectConfigurationID,
			DestinationPhoneNumber: r.DestinationPhoneNumber,
		})
	return err
}

func (r *PinpointSMSVoiceV2ProtectConfigRuleSetOverride) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2ProtectConfigRuleSetOverride) String() string {
	return fmt.Sprintf("%s -> %s", *r.ProtectConfigurationID, *r.DestinationPhoneNumber)
}

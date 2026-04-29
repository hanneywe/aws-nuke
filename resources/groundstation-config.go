package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/groundstation"
	groundstationtypes "github.com/aws/aws-sdk-go-v2/service/groundstation/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GroundStationConfigResource = "GroundStationConfig"

func init() {
	registry.Register(&registry.Registration{
		Name:     GroundStationConfigResource,
		Scope:    nuke.Account,
		Resource: &GroundStationConfig{},
		Lister:   &GroundStationConfigLister{},
	})
}

type GroundStationConfigLister struct {
	svc GroundStationClient
}

func (l *GroundStationConfigLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = groundstation.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := groundstation.NewListConfigsPaginator(svc, &groundstation.ListConfigsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cfg := range resp.ConfigList {
			resources = append(resources, &GroundStationConfig{
				svc:        svc,
				ConfigID:   cfg.ConfigId,
				ConfigType: cfg.ConfigType,
				Name:       cfg.Name,
			})
		}
	}

	return resources, nil
}

type GroundStationConfig struct {
	svc        GroundStationClient
	ConfigID   *string `property:"name=ConfigId"`
	ConfigType groundstationtypes.ConfigCapabilityType
	Name       *string
}

func (r *GroundStationConfig) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConfig(ctx, &groundstation.DeleteConfigInput{
		ConfigId:   r.ConfigID,
		ConfigType: r.ConfigType,
	})
	return err
}

func (r *GroundStationConfig) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GroundStationConfig) String() string {
	return *r.ConfigID
}

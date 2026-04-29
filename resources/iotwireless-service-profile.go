package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessServiceProfileResource = "IoTWirelessServiceProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessServiceProfileResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessServiceProfile{},
		Lister:   &IoTWirelessServiceProfileLister{},
	})
}

type IoTWirelessServiceProfileLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessServiceProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListServiceProfilesPaginator(svc, &iotwireless.ListServiceProfilesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sp := range resp.ServiceProfileList {
			resources = append(resources, &IoTWirelessServiceProfile{
				svc:  svc,
				ID:   sp.Id,
				Name: sp.Name,
			})
		}
	}

	return resources, nil
}

type IoTWirelessServiceProfile struct {
	svc  IoTWirelessClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *IoTWirelessServiceProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServiceProfile(ctx, &iotwireless.DeleteServiceProfileInput{
		Id: r.ID,
	})
	return err
}

func (r *IoTWirelessServiceProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessServiceProfile) String() string {
	return *r.ID
}

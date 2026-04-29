package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessMulticastGroupResource = "IoTWirelessMulticastGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessMulticastGroupResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessMulticastGroup{},
		Lister:   &IoTWirelessMulticastGroupLister{},
	})
}

type IoTWirelessMulticastGroupLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessMulticastGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListMulticastGroupsPaginator(svc, &iotwireless.ListMulticastGroupsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, group := range resp.MulticastGroupList {
			resources = append(resources, &IoTWirelessMulticastGroup{
				svc:  svc,
				ID:   group.Id,
				Name: group.Name,
			})
		}
	}

	return resources, nil
}

type IoTWirelessMulticastGroup struct {
	svc  IoTWirelessClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *IoTWirelessMulticastGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMulticastGroup(ctx, &iotwireless.DeleteMulticastGroupInput{
		Id: r.ID,
	})
	return err
}

func (r *IoTWirelessMulticastGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessMulticastGroup) String() string {
	return *r.ID
}

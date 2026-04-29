package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessDeviceProfileResource = "IoTWirelessDeviceProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessDeviceProfileResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessDeviceProfile{},
		Lister:   &IoTWirelessDeviceProfileLister{},
	})
}

type IoTWirelessDeviceProfileLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessDeviceProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListDeviceProfilesPaginator(svc, &iotwireless.ListDeviceProfilesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, deviceProfile := range resp.DeviceProfileList {
			resources = append(resources, &IoTWirelessDeviceProfile{
				svc:  svc,
				ID:   deviceProfile.Id,
				Name: deviceProfile.Name,
			})
		}
	}

	return resources, nil
}

type IoTWirelessDeviceProfile struct {
	svc  IoTWirelessClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *IoTWirelessDeviceProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDeviceProfile(ctx, &iotwireless.DeleteDeviceProfileInput{
		Id: r.ID,
	})
	return err
}

func (r *IoTWirelessDeviceProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessDeviceProfile) String() string {
	return *r.ID
}

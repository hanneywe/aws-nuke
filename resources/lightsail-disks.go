package resources

import (
	"context"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go/service/lightsail" //nolint:staticcheck

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailDiskResource = "LightsailDisk"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailDiskResource,
		Scope:    nuke.Account,
		Resource: &LightsailDisk{},
		Lister:   &LightsailDiskLister{},
		Settings: []string{
			"ForceDeleteAddOns",
		},
	})
}

type LightsailDiskLister struct{}

func (l *LightsailDiskLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := lightsail.New(opts.Session)
	resources := make([]resource.Resource, 0)

	params := &lightsail.GetDisksInput{}

	for {
		output, err := svc.GetDisks(params)
		if err != nil {
			return nil, err
		}

		for _, disk := range output.Disks {
			resources = append(resources, &LightsailDisk{
				svc:      svc,
				diskName: disk.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

type LightsailDisk struct {
	svc      *lightsail.Lightsail
	diskName *string
	settings *libsettings.Setting
}

func (f *LightsailDisk) Settings(setting *libsettings.Setting) {
	f.settings = setting
}

func (f *LightsailDisk) Remove(_ context.Context) error {
	_, err := f.svc.DeleteDisk(&lightsail.DeleteDiskInput{
		DiskName:          f.diskName,
		ForceDeleteAddOns: ptr.Bool(f.settings.GetBool("ForceDeleteAddOns")),
	})

	return err
}

func (f *LightsailDisk) String() string {
	return *f.diskName
}

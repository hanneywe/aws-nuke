package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeviceFarmInstanceProfileResource = "DeviceFarmInstanceProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeviceFarmInstanceProfileResource,
		Scope:    nuke.Account,
		Resource: &DeviceFarmInstanceProfile{},
		Lister:   &DeviceFarmInstanceProfileLister{},
	})
}

type DeviceFarmInstanceProfileLister struct {
	svc DeviceFarmClient
}

func (l *DeviceFarmInstanceProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = devicefarm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &devicefarm.ListInstanceProfilesInput{}
	for {
		listOutput, err := svc.ListInstanceProfiles(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, instanceProfile := range listOutput.InstanceProfiles {
			resources = append(resources, &DeviceFarmInstanceProfile{
				svc:  svc,
				Arn:  instanceProfile.Arn,
				Name: instanceProfile.Name,
			})
		}

		if listOutput.NextToken == nil {
			break
		}
		params.NextToken = listOutput.NextToken
	}

	return resources, nil
}

type DeviceFarmInstanceProfile struct {
	svc  DeviceFarmClient
	Arn  *string
	Name *string
}

func (r *DeviceFarmInstanceProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteInstanceProfile(ctx, &devicefarm.DeleteInstanceProfileInput{
		Arn: r.Arn,
	})
	return err
}

func (r *DeviceFarmInstanceProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeviceFarmInstanceProfile) String() string {
	return *r.Name
}

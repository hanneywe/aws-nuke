package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailContainerServiceResource = "LightsailContainerService"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailContainerServiceResource,
		Scope:    nuke.Account,
		Resource: &LightsailContainerService{},
		Lister:   &LightsailContainerServiceLister{},
	})
}

type LightsailContainerServiceLister struct {
	svc LightsailClient
}

func (l *LightsailContainerServiceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	resp, err := svc.GetContainerServices(ctx, &lightsail.GetContainerServicesInput{})
	if err != nil {
		return nil, err
	}

	for i := range resp.ContainerServices {
		cs := &resp.ContainerServices[i]
		resources = append(resources, &LightsailContainerService{
			svc:                  svc,
			ContainerServiceName: cs.ContainerServiceName,
		})
	}

	return resources, nil
}

type LightsailContainerService struct {
	svc                  LightsailClient
	ContainerServiceName *string
}

func (r *LightsailContainerService) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteContainerService(ctx, &lightsail.DeleteContainerServiceInput{
		ServiceName: r.ContainerServiceName,
	})
	return err
}

func (r *LightsailContainerService) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailContainerService) String() string {
	return *r.ContainerServiceName
}

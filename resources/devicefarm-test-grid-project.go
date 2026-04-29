package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeviceFarmTestGridProjectResource = "DeviceFarmTestGridProject"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeviceFarmTestGridProjectResource,
		Scope:    nuke.Account,
		Resource: &DeviceFarmTestGridProject{},
		Lister:   &DeviceFarmTestGridProjectLister{},
	})
}

type DeviceFarmTestGridProjectLister struct {
	svc DeviceFarmClient
}

func (l *DeviceFarmTestGridProjectLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = devicefarm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := devicefarm.NewListTestGridProjectsPaginator(svc, &devicefarm.ListTestGridProjectsInput{})
	for paginator.HasMorePages() {
		listOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, project := range listOutput.TestGridProjects {
			resources = append(resources, &DeviceFarmTestGridProject{
				svc:  svc,
				Arn:  project.Arn,
				Name: project.Name,
			})
		}
	}

	return resources, nil
}

type DeviceFarmTestGridProject struct {
	svc  DeviceFarmClient
	Arn  *string
	Name *string
}

func (r *DeviceFarmTestGridProject) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTestGridProject(ctx, &devicefarm.DeleteTestGridProjectInput{
		ProjectArn: r.Arn,
	})
	return err
}

func (r *DeviceFarmTestGridProject) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeviceFarmTestGridProject) String() string {
	return *r.Name
}

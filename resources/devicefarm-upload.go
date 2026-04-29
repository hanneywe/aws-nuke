package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeviceFarmUploadResource = "DeviceFarmUpload"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeviceFarmUploadResource,
		Scope:    nuke.Account,
		Resource: &DeviceFarmUpload{},
		Lister:   &DeviceFarmUploadLister{},
	})
}

type DeviceFarmUploadLister struct {
	svc DeviceFarmClient
}

func (l *DeviceFarmUploadLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = devicefarm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, list all projects
	projectPaginator := devicefarm.NewListProjectsPaginator(svc, &devicefarm.ListProjectsInput{})
	for projectPaginator.HasMorePages() {
		projectsOutput, err := projectPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, project := range projectsOutput.Projects {
			// For each project, list all uploads
			uploadPaginator := devicefarm.NewListUploadsPaginator(svc, &devicefarm.ListUploadsInput{
				Arn: project.Arn,
			})
			for uploadPaginator.HasMorePages() {
				uploadsOutput, err := uploadPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, upload := range uploadsOutput.Uploads {
					resources = append(resources, &DeviceFarmUpload{
						svc:  svc,
						Arn:  upload.Arn,
						Name: upload.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type DeviceFarmUpload struct {
	svc  DeviceFarmClient
	Arn  *string
	Name *string
}

func (r *DeviceFarmUpload) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUpload(ctx, &devicefarm.DeleteUploadInput{
		Arn: r.Arn,
	})
	return err
}

func (r *DeviceFarmUpload) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeviceFarmUpload) String() string {
	return *r.Name
}

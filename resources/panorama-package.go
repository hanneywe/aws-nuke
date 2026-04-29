package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/panorama"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PanoramaPackageResource = "PanoramaPackage"

func init() {
	registry.Register(&registry.Registration{
		Name:     PanoramaPackageResource,
		Scope:    nuke.Account,
		Resource: &PanoramaPackage{},
		Lister:   &PanoramaPackageLister{},
	})
}

type PanoramaPackageLister struct {
	svc PanoramaClient
}

func (l *PanoramaPackageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = panorama.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := panorama.NewListPackagesPaginator(svc, &panorama.ListPackagesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, pkg := range resp.Packages {
			resources = append(resources, &PanoramaPackage{
				svc:         svc,
				PackageID:   pkg.PackageId,
				PackageName: pkg.PackageName,
			})
		}
	}
	return resources, nil
}

type PanoramaPackage struct {
	svc         PanoramaClient
	PackageID   *string `property:"name=PackageId"`
	PackageName *string
}

func (r *PanoramaPackage) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePackage(ctx, &panorama.DeletePackageInput{
		PackageId:   r.PackageID,
		ForceDelete: true,
	})
	return err
}

func (r *PanoramaPackage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PanoramaPackage) String() string {
	return *r.PackageName
}

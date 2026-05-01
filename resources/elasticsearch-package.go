package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	estypes "github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ElasticsearchPackageResource = "ElasticsearchPackage"

func init() {
	registry.Register(&registry.Registration{
		Name:     ElasticsearchPackageResource,
		Scope:    nuke.Account,
		Resource: &ElasticsearchPackage{},
		Lister:   &ElasticsearchPackageLister{},
	})
}

type ElasticsearchPackageLister struct {
	svc ElasticsearchserviceClient
}

func (l *ElasticsearchPackageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = elasticsearchservice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := elasticsearchservice.NewDescribePackagesPaginator(svc, &elasticsearchservice.DescribePackagesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.PackageDetailsList {
			item := &resp.PackageDetailsList[i]
			resources = append(resources, &ElasticsearchPackage{
				svc:           svc,
				PackageID:     item.PackageID,
				PackageName:   item.PackageName,
				PackageStatus: item.PackageStatus,
			})
		}
	}

	return resources, nil
}

type ElasticsearchPackage struct {
	svc           ElasticsearchserviceClient
	PackageID     *string
	PackageName   *string
	PackageStatus estypes.PackageStatus
}

func (r *ElasticsearchPackage) Filter() error {
	if r.PackageID != nil && !strings.HasPrefix(*r.PackageID, "pkg-") {
		return fmt.Errorf("AWS-managed package included with OpenSearch Service")
	}
	if r.PackageStatus == estypes.PackageStatusDeleted || r.PackageStatus == estypes.PackageStatusDeleting {
		return fmt.Errorf("package already %s", r.PackageStatus)
	}
	return nil
}

func (r *ElasticsearchPackage) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePackage(ctx, &elasticsearchservice.DeletePackageInput{
		PackageID: r.PackageID,
	})
	return err
}

func (r *ElasticsearchPackage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ElasticsearchPackage) String() string {
	return *r.PackageID
}

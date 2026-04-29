package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeServiceResource = "VPCLatticeService"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeServiceResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeService{},
		Lister:   &VPCLatticeServiceLister{},
		DependsOn: []string{
			VPCLatticeListenerResource,
			VPCLatticeAuthPolicyResource,
			VPCLatticeServiceNetworkServiceAssociationResource,
		},
	})
}

type VPCLatticeServiceLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeServiceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := vpclattice.NewListServicesPaginator(svc, &vpclattice.ListServicesInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			var tags map[string]string
			if item.Arn != nil {
				tagsResp, err := svc.ListTagsForResource(ctx, &vpclattice.ListTagsForResourceInput{
					ResourceArn: item.Arn,
				})
				if err != nil {
					opts.Logger.Warnf("unable to fetch tags for service: %s", *item.Arn)
				} else {
					tags = tagsResp.Tags
				}
			}

			resources = append(resources, &VPCLatticeService{
				svc:  svc,
				ID:   item.Id,
				ARN:  item.Arn,
				Name: item.Name,
				Tags: tags,
			})
		}
	}

	return resources, nil
}

type VPCLatticeService struct {
	svc  VPCLatticeClient
	ID   *string
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *VPCLatticeService) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteService(ctx, &vpclattice.DeleteServiceInput{
		ServiceIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeService) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeService) String() string {
	return *r.Name
}

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

const VPCLatticeServiceNetworkResource = "VPCLatticeServiceNetwork"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeServiceNetworkResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeServiceNetwork{},
		Lister:   &VPCLatticeServiceNetworkLister{},
		DependsOn: []string{
			VPCLatticeServiceNetworkServiceAssociationResource,
			VPCLatticeServiceNetworkVPCAssociationResource,
		},
	})
}

type VPCLatticeServiceNetworkLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeServiceNetworkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := vpclattice.NewListServiceNetworksPaginator(svc, &vpclattice.ListServiceNetworksInput{
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
					opts.Logger.Warnf("unable to fetch tags for service network: %s", *item.Arn)
				} else {
					tags = tagsResp.Tags
				}
			}

			resources = append(resources, &VPCLatticeServiceNetwork{
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

type VPCLatticeServiceNetwork struct {
	svc  VPCLatticeClient
	ID   *string
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *VPCLatticeServiceNetwork) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServiceNetwork(ctx, &vpclattice.DeleteServiceNetworkInput{
		ServiceNetworkIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeServiceNetwork) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeServiceNetwork) String() string {
	return *r.Name
}

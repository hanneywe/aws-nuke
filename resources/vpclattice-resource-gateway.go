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

const VPCLatticeResourceGatewayResource = "VPCLatticeResourceGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeResourceGatewayResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeResourceGateway{},
		Lister:   &VPCLatticeResourceGatewayLister{},
	})
}

type VPCLatticeResourceGatewayLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeResourceGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := vpclattice.NewListResourceGatewaysPaginator(svc, &vpclattice.ListResourceGatewaysInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Items {
			var tags map[string]string
			if resp.Items[i].Arn != nil {
				tagsResp, err := svc.ListTagsForResource(ctx, &vpclattice.ListTagsForResourceInput{
					ResourceArn: resp.Items[i].Arn,
				})
				if err != nil {
					opts.Logger.Warnf("unable to fetch tags for resource gateway: %s", *resp.Items[i].Arn)
				} else {
					tags = tagsResp.Tags
				}
			}

			resources = append(resources, &VPCLatticeResourceGateway{
				svc:  svc,
				ID:   resp.Items[i].Id,
				ARN:  resp.Items[i].Arn,
				Name: resp.Items[i].Name,
				Tags: tags,
			})
		}
	}

	return resources, nil
}

type VPCLatticeResourceGateway struct {
	svc  VPCLatticeClient
	ID   *string
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *VPCLatticeResourceGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteResourceGateway(ctx, &vpclattice.DeleteResourceGatewayInput{
		ResourceGatewayIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeResourceGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeResourceGateway) String() string {
	return *r.Name
}

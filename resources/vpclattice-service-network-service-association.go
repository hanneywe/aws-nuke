package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeServiceNetworkServiceAssociationResource = "VPCLatticeServiceNetworkServiceAssociation"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeServiceNetworkServiceAssociationResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeServiceNetworkServiceAssociation{},
		Lister:   &VPCLatticeServiceNetworkServiceAssociationLister{},
	})
}

type VPCLatticeServiceNetworkServiceAssociationLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeServiceNetworkServiceAssociationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	// The ListServiceNetworkServiceAssociations API requires at least one of
	// serviceNetworkIdentifier or serviceIdentifier, so we iterate all service networks first.
	snPaginator := vpclattice.NewListServiceNetworksPaginator(svc, &vpclattice.ListServiceNetworksInput{
		MaxResults: aws.Int32(100),
	})

	for snPaginator.HasMorePages() {
		snResp, err := snPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sn := range snResp.Items {
			paginator := vpclattice.NewListServiceNetworkServiceAssociationsPaginator(svc, &vpclattice.ListServiceNetworkServiceAssociationsInput{
				ServiceNetworkIdentifier: sn.Id,
				MaxResults:               aws.Int32(100),
			})

			for paginator.HasMorePages() {
				resp, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, item := range resp.Items {
					resources = append(resources, &VPCLatticeServiceNetworkServiceAssociation{
						svc:                svc,
						ID:                 item.Id,
						ARN:                item.Arn,
						ServiceNetworkName: item.ServiceNetworkName,
						ServiceName:        item.ServiceName,
						Status:             aws.String(string(item.Status)),
					})
				}
			}
		}
	}

	return resources, nil
}

type VPCLatticeServiceNetworkServiceAssociation struct {
	svc                VPCLatticeClient
	ID                 *string
	ARN                *string
	ServiceNetworkName *string
	ServiceName        *string
	Status             *string
}

func (r *VPCLatticeServiceNetworkServiceAssociation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServiceNetworkServiceAssociation(ctx, &vpclattice.DeleteServiceNetworkServiceAssociationInput{
		ServiceNetworkServiceAssociationIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeServiceNetworkServiceAssociation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeServiceNetworkServiceAssociation) String() string {
	return fmt.Sprintf("%s -> %s", *r.ServiceNetworkName, *r.ServiceName)
}

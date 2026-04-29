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

const VPCLatticeServiceNetworkVPCAssociationResource = "VPCLatticeServiceNetworkVPCAssociation"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeServiceNetworkVPCAssociationResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeServiceNetworkVPCAssociation{},
		Lister:   &VPCLatticeServiceNetworkVPCAssociationLister{},
	})
}

type VPCLatticeServiceNetworkVPCAssociationLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeServiceNetworkVPCAssociationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	// The ListServiceNetworkVpcAssociations API requires at least one of
	// serviceNetworkIdentifier or vpcIdentifier, so we iterate all service networks first.
	snPaginator := vpclattice.NewListServiceNetworksPaginator(svc, &vpclattice.ListServiceNetworksInput{
		MaxResults: aws.Int32(100),
	})

	for snPaginator.HasMorePages() {
		snResp, err := snPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sn := range snResp.Items {
			paginator := vpclattice.NewListServiceNetworkVpcAssociationsPaginator(svc, &vpclattice.ListServiceNetworkVpcAssociationsInput{
				ServiceNetworkIdentifier: sn.Id,
				MaxResults:               aws.Int32(100),
			})

			for paginator.HasMorePages() {
				resp, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, item := range resp.Items {
					resources = append(resources, &VPCLatticeServiceNetworkVPCAssociation{
						svc:                svc,
						ID:                 item.Id,
						ARN:                item.Arn,
						ServiceNetworkName: item.ServiceNetworkName,
						VPCID:              item.VpcId,
						Status:             aws.String(string(item.Status)),
					})
				}
			}
		}
	}

	return resources, nil
}

type VPCLatticeServiceNetworkVPCAssociation struct {
	svc                VPCLatticeClient
	ID                 *string
	ARN                *string
	ServiceNetworkName *string
	VPCID              *string
	Status             *string
}

func (r *VPCLatticeServiceNetworkVPCAssociation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteServiceNetworkVpcAssociation(ctx, &vpclattice.DeleteServiceNetworkVpcAssociationInput{
		ServiceNetworkVpcAssociationIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeServiceNetworkVPCAssociation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeServiceNetworkVPCAssociation) String() string {
	return fmt.Sprintf("%s -> %s", *r.ServiceNetworkName, *r.VPCID)
}

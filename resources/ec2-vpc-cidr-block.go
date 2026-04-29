package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2VPCCIDRBlockResource = "EC2VPCCIDRBlock"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2VPCCIDRBlockResource,
		Scope:    nuke.Account,
		Resource: &EC2VPCCIDRBlock{},
		Lister:   &EC2VPCCIDRBlockLister{},
	})
}

type EC2VPCCIDRBlockLister struct {
	svc EC2Client
}

func (l *EC2VPCCIDRBlockLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeVpcsPaginator(svc,
		&ec2.DescribeVpcsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Vpcs {
			// Add non-primary IPv4 CIDR associations
			for _, assoc := range resp.Vpcs[i].CidrBlockAssociationSet {
				// Skip the primary CIDR block — it is removed when the VPC itself is deleted
				if ptr.ToString(assoc.CidrBlock) == ptr.ToString(resp.Vpcs[i].CidrBlock) {
					continue
				}

				var state string
				if assoc.CidrBlockState != nil {
					state = string(assoc.CidrBlockState.State)
				}

				resources = append(resources, &EC2VPCCIDRBlock{
					svc:           svc,
					VpcID:         resp.Vpcs[i].VpcId,
					CidrBlock:     assoc.CidrBlock,
					AssociationID: assoc.AssociationId,
					State:         state,
					IsIPv6:        ptr.Bool(false),
				})
			}

			// Add IPv6 CIDR associations
			for _, assoc := range resp.Vpcs[i].Ipv6CidrBlockAssociationSet {
				var state string
				if assoc.Ipv6CidrBlockState != nil {
					state = string(assoc.Ipv6CidrBlockState.State)
				}

				resources = append(resources, &EC2VPCCIDRBlock{
					svc:           svc,
					VpcID:         resp.Vpcs[i].VpcId,
					CidrBlock:     assoc.Ipv6CidrBlock,
					AssociationID: assoc.AssociationId,
					State:         state,
					IsIPv6:        ptr.Bool(true),
				})
			}
		}
	}

	return resources, nil
}

type EC2VPCCIDRBlock struct {
	svc           EC2Client
	VpcID         *string `property:"name=VpcId"`
	CidrBlock     *string
	AssociationID *string `property:"name=AssociationId"`
	State         string
	IsIPv6        *bool
}

func (r *EC2VPCCIDRBlock) Remove(ctx context.Context) error {
	_, err := r.svc.DisassociateVpcCidrBlock(ctx, &ec2.DisassociateVpcCidrBlockInput{
		AssociationId: r.AssociationID,
	})
	return err
}

func (r *EC2VPCCIDRBlock) Filter() error {
	if r.State == string(ec2types.VpcCidrBlockStateCodeDisassociated) {
		return fmt.Errorf("already disassociated")
	}
	return nil
}

func (r *EC2VPCCIDRBlock) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2VPCCIDRBlock) String() string {
	return fmt.Sprintf("%s (%s)", ptr.ToString(r.CidrBlock), ptr.ToString(r.VpcID))
}

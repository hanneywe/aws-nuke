package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2IPAMPrefixListResolverResource = "EC2IPAMPrefixListResolver"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2IPAMPrefixListResolverResource,
		Scope:    nuke.Account,
		Resource: &EC2IPAMPrefixListResolver{},
		Lister:   &EC2IPAMPrefixListResolverLister{},
	})
}

type EC2IPAMPrefixListResolverLister struct {
	svc EC2Client
}

func (l *EC2IPAMPrefixListResolverLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeIpamPrefixListResolversPaginator(svc,
		&ec2.DescribeIpamPrefixListResolversInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.IpamPrefixListResolvers {
			r := &resp.IpamPrefixListResolvers[i]
			resources = append(resources, &EC2IPAMPrefixListResolver{
				svc:                      svc,
				IpamPrefixListResolverID: r.IpamPrefixListResolverId,
				IpamArn:                  r.IpamArn,
				OwnerID:                  r.OwnerId,
				State:                    string(r.State),
				Tags:                     r.Tags,
			})
		}
	}

	return resources, nil
}

type EC2IPAMPrefixListResolver struct {
	svc                      EC2Client
	IpamPrefixListResolverID *string `property:"name=IpamPrefixListResolverId"`
	IpamArn                  *string
	OwnerID                  *string `property:"name=OwnerId"`
	State                    string
	Tags                     []ec2types.Tag
}

func (r *EC2IPAMPrefixListResolver) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpamPrefixListResolver(ctx, &ec2.DeleteIpamPrefixListResolverInput{
		IpamPrefixListResolverId: r.IpamPrefixListResolverID,
	})
	return err
}

func (r *EC2IPAMPrefixListResolver) Filter() error {
	if r.State == string(ec2types.IpamPrefixListResolverStateDeleteComplete) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2IPAMPrefixListResolver) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2IPAMPrefixListResolver) String() string {
	return *r.IpamPrefixListResolverID
}

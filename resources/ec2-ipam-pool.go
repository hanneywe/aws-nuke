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

const EC2IPAMPoolResource = "EC2IPAMPool"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2IPAMPoolResource,
		Scope:    nuke.Account,
		Resource: &EC2IPAMPool{},
		Lister:   &EC2IPAMPoolLister{},
	})
}

type EC2IPAMPoolLister struct {
	svc EC2Client
}

func (l *EC2IPAMPoolLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeIpamPoolsPaginator(svc,
		&ec2.DescribeIpamPoolsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.IpamPools {
			resources = append(resources, &EC2IPAMPool{
				svc:           svc,
				IpamPoolID:    resp.IpamPools[i].IpamPoolId,
				IpamScopeID:   resp.IpamPools[i].IpamScopeArn,
				AddressFamily: string(resp.IpamPools[i].AddressFamily),
				State:         string(resp.IpamPools[i].State),
				Tags:          resp.IpamPools[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2IPAMPool struct {
	svc           EC2Client
	IpamPoolID    *string `property:"name=IpamPoolId"`
	IpamScopeID   *string `property:"name=IpamScopeId"`
	AddressFamily string
	State         string
	Tags          []ec2types.Tag
}

func (r *EC2IPAMPool) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpamPool(ctx, &ec2.DeleteIpamPoolInput{
		IpamPoolId: r.IpamPoolID,
	})
	return err
}

func (r *EC2IPAMPool) Filter() error {
	if r.State == string(ec2types.IpamPoolStateDeleteComplete) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2IPAMPool) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2IPAMPool) String() string {
	return *r.IpamPoolID
}

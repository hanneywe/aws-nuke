package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2PublicIpv4PoolResource = "EC2PublicIpv4Pool"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2PublicIpv4PoolResource,
		Scope:    nuke.Account,
		Resource: &EC2PublicIpv4Pool{},
		Lister:   &EC2PublicIpv4PoolLister{},
	})
}

type EC2PublicIpv4PoolLister struct {
	svc EC2Client
}

func (l *EC2PublicIpv4PoolLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribePublicIpv4PoolsPaginator(svc, &ec2.DescribePublicIpv4PoolsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, pool := range output.PublicIpv4Pools {
			tagMap := make(map[string]string)
			for _, tag := range pool.Tags {
				if tag.Key != nil && tag.Value != nil {
					tagMap[*tag.Key] = *tag.Value
				}
			}
			resources = append(resources, &EC2PublicIpv4Pool{
				svc:    svc,
				PoolID: pool.PoolId,
				Tags:   tagMap,
			})
		}
	}

	return resources, nil
}

type EC2PublicIpv4Pool struct {
	svc    EC2Client
	PoolID *string `property:"name=PoolId"`
	Tags   map[string]string
}

func (r *EC2PublicIpv4Pool) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePublicIpv4Pool(ctx, &ec2.DeletePublicIpv4PoolInput{
		PoolId: r.PoolID,
	})
	return err
}

func (r *EC2PublicIpv4Pool) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2PublicIpv4Pool) String() string {
	return *r.PoolID
}

package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2IPAMResource = "EC2IPAM"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2IPAMResource,
		Scope:    nuke.Account,
		Resource: &EC2IPAM{},
		Lister:   &EC2IPAMLister{},
		DependsOn: []string{
			EC2IPAMPoolResource,
			EC2IPAMScopeResource,
		},
	})
}

type EC2IPAMLister struct {
	svc EC2Client
}

func (l *EC2IPAMLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeIpamsPaginator(svc,
		&ec2.DescribeIpamsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Ipams {
			resources = append(resources, &EC2IPAM{
				svc:     svc,
				IpamID:  resp.Ipams[i].IpamId,
				OwnerID: resp.Ipams[i].OwnerId,
				State:   string(resp.Ipams[i].State),
				Tags:    resp.Ipams[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2IPAM struct {
	svc     EC2Client
	IpamID  *string `property:"name=IpamId"`
	OwnerID *string `property:"name=OwnerId"`
	State   string
	Tags    []ec2types.Tag
}

func (r *EC2IPAM) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpam(ctx, &ec2.DeleteIpamInput{
		IpamId:  r.IpamID,
		Cascade: aws.Bool(true),
	})
	return err
}

func (r *EC2IPAM) Filter() error {
	if r.State == string(ec2types.IpamStateDeleteComplete) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2IPAM) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2IPAM) String() string {
	return *r.IpamID
}

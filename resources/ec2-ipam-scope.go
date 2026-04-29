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

const EC2IPAMScopeResource = "EC2IPAMScope"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2IPAMScopeResource,
		Scope:    nuke.Account,
		Resource: &EC2IPAMScope{},
		Lister:   &EC2IPAMScopeLister{},
		DependsOn: []string{
			EC2IPAMPoolResource,
		},
	})
}

type EC2IPAMScopeLister struct {
	svc EC2Client
}

func (l *EC2IPAMScopeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeIpamScopesPaginator(svc,
		&ec2.DescribeIpamScopesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.IpamScopes {
			resources = append(resources, &EC2IPAMScope{
				svc:         svc,
				IpamScopeID: resp.IpamScopes[i].IpamScopeId,
				IpamID:      resp.IpamScopes[i].IpamArn,
				IsDefault:   resp.IpamScopes[i].IsDefault,
				State:       string(resp.IpamScopes[i].State),
				Tags:        resp.IpamScopes[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2IPAMScope struct {
	svc         EC2Client
	IpamScopeID *string `property:"name=IpamScopeId"`
	IpamID      *string `property:"name=IpamId"`
	IsDefault   *bool
	State       string
	Tags        []ec2types.Tag
}

func (r *EC2IPAMScope) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpamScope(ctx, &ec2.DeleteIpamScopeInput{
		IpamScopeId: r.IpamScopeID,
	})
	return err
}

func (r *EC2IPAMScope) Filter() error {
	if ptr.ToBool(r.IsDefault) {
		return fmt.Errorf("cannot delete default IPAM scope")
	}
	if r.State == string(ec2types.IpamScopeStateDeleteComplete) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2IPAMScope) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2IPAMScope) String() string {
	return *r.IpamScopeID
}

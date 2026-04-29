package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const (
	EC2RouteServerResource      = "EC2RouteServer"
	EC2RouteServerStateDeleted  = "deleted"
	EC2RouteServerStateDeleting = "deleting"
)

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2RouteServerResource,
		Scope:    nuke.Account,
		Resource: &EC2RouteServer{},
		Lister:   &EC2RouteServerLister{},
	})
}

type EC2RouteServerLister struct {
	svc EC2Client
}

func (l *EC2RouteServerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ec2.DescribeRouteServersInput{}
	for {
		resp, err := svc.DescribeRouteServers(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, rs := range resp.RouteServers {
			resources = append(resources, &EC2RouteServer{
				svc:           svc,
				RouteServerID: rs.RouteServerId,
				State:         ptr.String(string(rs.State)),
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type EC2RouteServer struct {
	svc           EC2Client
	RouteServerID *string `property:"name=RouteServerId"`
	State         *string
}

func (r *EC2RouteServer) Filter() error {
	if ptr.ToString(r.State) == EC2RouteServerStateDeleting || ptr.ToString(r.State) == EC2RouteServerStateDeleted {
		return fmt.Errorf("already %s", ptr.ToString(r.State))
	}
	return nil
}

func (r *EC2RouteServer) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRouteServer(ctx, &ec2.DeleteRouteServerInput{
		RouteServerId: r.RouteServerID,
	})
	return err
}

func (r *EC2RouteServer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2RouteServer) String() string {
	return *r.RouteServerID
}

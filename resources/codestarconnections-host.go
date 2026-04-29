package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codestarconnections"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeStarConnectionsHostResource = "CodeStarConnectionsHost"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeStarConnectionsHostResource,
		Scope:    nuke.Account,
		Resource: &CodeStarConnectionsHost{},
		Lister:   &CodeStarConnectionsHostLister{},
	})
}

type CodeStarConnectionsHostLister struct {
	svc CodeStarConnectionsClient
}

func (l *CodeStarConnectionsHostLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = codestarconnections.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &codestarconnections.ListHostsInput{}
	for {
		resp, err := svc.ListHosts(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, h := range resp.Hosts {
			resources = append(resources, &CodeStarConnectionsHost{
				svc:     svc,
				HostArn: h.HostArn,
				Name:    h.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type CodeStarConnectionsHost struct {
	svc     CodeStarConnectionsClient
	HostArn *string
	Name    *string
}

func (r *CodeStarConnectionsHost) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHost(ctx, &codestarconnections.DeleteHostInput{
		HostArn: r.HostArn,
	})
	return err
}

func (r *CodeStarConnectionsHost) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeStarConnectionsHost) String() string {
	return *r.Name
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/synthetics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SyntheticsGroupResource = "SyntheticsGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     SyntheticsGroupResource,
		Scope:    nuke.Account,
		Resource: &SyntheticsGroup{},
		Lister:   &SyntheticsGroupLister{},
	})
}

type SyntheticsGroupLister struct {
	svc SyntheticsClient
}

func (l *SyntheticsGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = synthetics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := synthetics.NewListGroupsPaginator(svc, &synthetics.ListGroupsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, g := range resp.Groups {
			resources = append(resources, &SyntheticsGroup{
				svc:     svc,
				GroupID: g.Id,
				Name:    g.Name,
			})
		}
	}
	return resources, nil
}

type SyntheticsGroup struct {
	svc     SyntheticsClient
	GroupID *string `property:"name=GroupId"`
	Name    *string
}

func (r *SyntheticsGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGroup(ctx, &synthetics.DeleteGroupInput{
		GroupIdentifier: r.GroupID,
	})
	return err
}

func (r *SyntheticsGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SyntheticsGroup) String() string {
	return *r.Name
}

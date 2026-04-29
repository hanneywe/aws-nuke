package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspaces"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const WorkSpacesIPGroupResource = "WorkSpacesIpGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     WorkSpacesIPGroupResource,
		Scope:    nuke.Account,
		Resource: &WorkSpacesIPGroup{},
		Lister:   &WorkSpacesIPGroupLister{},
	})
}

type WorkSpacesIPGroupLister struct {
	svc WorkSpacesV2Client
}

func (l *WorkSpacesIPGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = workspaces.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &workspaces.DescribeIpGroupsInput{}
	for {
		resp, err := svc.DescribeIpGroups(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Result {
			resources = append(resources, &WorkSpacesIPGroup{
				svc:       svc,
				GroupID:   resp.Result[i].GroupId,
				GroupName: resp.Result[i].GroupName,
				GroupDesc: resp.Result[i].GroupDesc,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type WorkSpacesIPGroup struct {
	svc       WorkSpacesV2Client
	GroupID   *string `property:"name=GroupId"`
	GroupName *string
	GroupDesc *string
}

func (r *WorkSpacesIPGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpGroup(ctx, &workspaces.DeleteIpGroupInput{
		GroupId: r.GroupID,
	})
	return err
}

func (r *WorkSpacesIPGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *WorkSpacesIPGroup) String() string {
	return *r.GroupID
}

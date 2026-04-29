package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/neptune"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NeptuneDBClusterParameterGroupResource = "NeptuneDBClusterParameterGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneDBClusterParameterGroupResource,
		Scope:    nuke.Account,
		Resource: &NeptuneDBClusterParameterGroup{},
		Lister:   &NeptuneDBClusterParameterGroupLister{},
	})
}

type NeptuneDBClusterParameterGroupLister struct {
	svc NeptuneV2Client
}

func (l *NeptuneDBClusterParameterGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptune.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := neptune.NewDescribeDBClusterParameterGroupsPaginator(svc, &neptune.DescribeDBClusterParameterGroupsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, clusterParameterGroup := range output.DBClusterParameterGroups {
			resources = append(resources, &NeptuneDBClusterParameterGroup{
				svc:                         svc,
				DBClusterParameterGroupName: clusterParameterGroup.DBClusterParameterGroupName,
				DBClusterParameterGroupArn:  clusterParameterGroup.DBClusterParameterGroupArn,
			})
		}
	}

	return resources, nil
}

type NeptuneDBClusterParameterGroup struct {
	svc NeptuneV2Client

	DBClusterParameterGroupName *string
	DBClusterParameterGroupArn  *string
}

func (r *NeptuneDBClusterParameterGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDBClusterParameterGroup(ctx, &neptune.DeleteDBClusterParameterGroupInput{
		DBClusterParameterGroupName: r.DBClusterParameterGroupName,
	})
	return err
}

func (r *NeptuneDBClusterParameterGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneDBClusterParameterGroup) String() string {
	return *r.DBClusterParameterGroupName
}

func (r *NeptuneDBClusterParameterGroup) Filter() error {
	if r.DBClusterParameterGroupName != nil && strings.HasPrefix(*r.DBClusterParameterGroupName, "default.") {
		return fmt.Errorf("cannot delete default cluster parameter group")
	}
	return nil
}

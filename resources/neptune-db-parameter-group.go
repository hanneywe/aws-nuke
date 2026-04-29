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

const NeptuneDBParameterGroupResource = "NeptuneDBParameterGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneDBParameterGroupResource,
		Scope:    nuke.Account,
		Resource: &NeptuneDBParameterGroup{},
		Lister:   &NeptuneDBParameterGroupLister{},
	})
}

type NeptuneDBParameterGroupLister struct {
	svc NeptuneV2Client
}

func (l *NeptuneDBParameterGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptune.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := neptune.NewDescribeDBParameterGroupsPaginator(svc, &neptune.DescribeDBParameterGroupsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, parameterGroup := range output.DBParameterGroups {
			resources = append(resources, &NeptuneDBParameterGroup{
				svc:                  svc,
				DBParameterGroupName: parameterGroup.DBParameterGroupName,
				DBParameterGroupArn:  parameterGroup.DBParameterGroupArn,
			})
		}
	}

	return resources, nil
}

type NeptuneDBParameterGroup struct {
	svc NeptuneV2Client

	DBParameterGroupName *string
	DBParameterGroupArn  *string
}

func (r *NeptuneDBParameterGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDBParameterGroup(ctx, &neptune.DeleteDBParameterGroupInput{
		DBParameterGroupName: r.DBParameterGroupName,
	})
	return err
}

func (r *NeptuneDBParameterGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneDBParameterGroup) String() string {
	return *r.DBParameterGroupName
}

func (r *NeptuneDBParameterGroup) Filter() error {
	if r.DBParameterGroupName != nil && strings.HasPrefix(*r.DBParameterGroupName, "default.") {
		return fmt.Errorf("cannot delete default parameter group")
	}
	return nil
}

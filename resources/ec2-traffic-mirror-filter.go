package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2TrafficMirrorFilterResource = "EC2TrafficMirrorFilter"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TrafficMirrorFilterResource,
		Scope:    nuke.Account,
		Resource: &EC2TrafficMirrorFilter{},
		Lister:   &EC2TrafficMirrorFilterLister{},
		DependsOn: []string{
			EC2TrafficMirrorSessionResource,
		},
	})
}

type EC2TrafficMirrorFilterLister struct {
	svc EC2Client
}

func (l *EC2TrafficMirrorFilterLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTrafficMirrorFiltersPaginator(svc,
		&ec2.DescribeTrafficMirrorFiltersInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, filter := range resp.TrafficMirrorFilters {
			resources = append(resources, &EC2TrafficMirrorFilter{
				svc:                   svc,
				TrafficMirrorFilterID: filter.TrafficMirrorFilterId,
				Description:           filter.Description,
				Tags:                  filter.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TrafficMirrorFilter struct {
	svc                   EC2Client
	TrafficMirrorFilterID *string `property:"name=TrafficMirrorFilterId"`
	Description           *string
	Tags                  []ec2types.Tag
}

func (r *EC2TrafficMirrorFilter) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTrafficMirrorFilter(ctx, &ec2.DeleteTrafficMirrorFilterInput{
		TrafficMirrorFilterId: r.TrafficMirrorFilterID,
	})
	return err
}

func (r *EC2TrafficMirrorFilter) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TrafficMirrorFilter) String() string {
	return *r.TrafficMirrorFilterID
}

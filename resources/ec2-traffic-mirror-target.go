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

const EC2TrafficMirrorTargetResource = "EC2TrafficMirrorTarget"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TrafficMirrorTargetResource,
		Scope:    nuke.Account,
		Resource: &EC2TrafficMirrorTarget{},
		Lister:   &EC2TrafficMirrorTargetLister{},
		DependsOn: []string{
			EC2TrafficMirrorSessionResource,
		},
	})
}

type EC2TrafficMirrorTargetLister struct {
	svc EC2Client
}

func (l *EC2TrafficMirrorTargetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTrafficMirrorTargetsPaginator(svc,
		&ec2.DescribeTrafficMirrorTargetsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, target := range resp.TrafficMirrorTargets {
			resources = append(resources, &EC2TrafficMirrorTarget{
				svc:                    svc,
				TrafficMirrorTargetID:  target.TrafficMirrorTargetId,
				Type:                   string(target.Type),
				NetworkInterfaceID:     target.NetworkInterfaceId,
				NetworkLoadBalancerArn: target.NetworkLoadBalancerArn,
				Tags:                   target.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TrafficMirrorTarget struct {
	svc                    EC2Client
	TrafficMirrorTargetID  *string `property:"name=TrafficMirrorTargetId"`
	Type                   string
	NetworkInterfaceID     *string `property:"name=NetworkInterfaceId"`
	NetworkLoadBalancerArn *string
	Tags                   []ec2types.Tag
}

func (r *EC2TrafficMirrorTarget) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTrafficMirrorTarget(ctx, &ec2.DeleteTrafficMirrorTargetInput{
		TrafficMirrorTargetId: r.TrafficMirrorTargetID,
	})
	return err
}

func (r *EC2TrafficMirrorTarget) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TrafficMirrorTarget) String() string {
	return *r.TrafficMirrorTargetID
}

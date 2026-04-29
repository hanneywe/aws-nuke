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

const EC2TrafficMirrorSessionResource = "EC2TrafficMirrorSession"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TrafficMirrorSessionResource,
		Scope:    nuke.Account,
		Resource: &EC2TrafficMirrorSession{},
		Lister:   &EC2TrafficMirrorSessionLister{},
	})
}

type EC2TrafficMirrorSessionLister struct {
	svc EC2Client
}

func (l *EC2TrafficMirrorSessionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTrafficMirrorSessionsPaginator(svc,
		&ec2.DescribeTrafficMirrorSessionsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, session := range resp.TrafficMirrorSessions {
			resources = append(resources, &EC2TrafficMirrorSession{
				svc:                    svc,
				TrafficMirrorSessionID: session.TrafficMirrorSessionId,
				TrafficMirrorTargetID:  session.TrafficMirrorTargetId,
				TrafficMirrorFilterID:  session.TrafficMirrorFilterId,
				NetworkInterfaceID:     session.NetworkInterfaceId,
				Tags:                   session.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TrafficMirrorSession struct {
	svc                    EC2Client
	TrafficMirrorSessionID *string `property:"name=TrafficMirrorSessionId"`
	TrafficMirrorTargetID  *string `property:"name=TrafficMirrorTargetId"`
	TrafficMirrorFilterID  *string `property:"name=TrafficMirrorFilterId"`
	NetworkInterfaceID     *string `property:"name=NetworkInterfaceId"`
	Tags                   []ec2types.Tag
}

func (r *EC2TrafficMirrorSession) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTrafficMirrorSession(ctx, &ec2.DeleteTrafficMirrorSessionInput{
		TrafficMirrorSessionId: r.TrafficMirrorSessionID,
	})
	return err
}

func (r *EC2TrafficMirrorSession) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TrafficMirrorSession) String() string {
	return *r.TrafficMirrorSessionID
}

package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2FleetResource = "EC2Fleet"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2FleetResource,
		Scope:    nuke.Account,
		Resource: &EC2Fleet{},
		Lister:   &EC2FleetLister{},
	})
}

type EC2FleetLister struct {
	svc EC2Client
}

func (l *EC2FleetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeFleetsPaginator(svc, &ec2.DescribeFleetsInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for index := range resp.Fleets {
			fleet := resp.Fleets[index]
			resources = append(resources, &EC2Fleet{
				svc:        svc,
				FleetID:    fleet.FleetId,
				FleetState: fleet.FleetState,
			})
		}
	}

	return resources, nil
}

type EC2Fleet struct {
	svc        EC2Client
	FleetID    *string `property:"name=FleetId"`
	FleetState ec2types.FleetStateCode
}

func (r *EC2Fleet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFleets(ctx, &ec2.DeleteFleetsInput{
		FleetIds:           []string{*r.FleetID},
		TerminateInstances: aws.Bool(false),
	})
	return err
}

func (r *EC2Fleet) Filter() error {
	if r.FleetState == ec2types.FleetStateCodeDeleted ||
		r.FleetState == ec2types.FleetStateCodeDeletedRunning ||
		r.FleetState == ec2types.FleetStateCodeDeletedTerminatingInstances {
		return fmt.Errorf("fleet is in %s state", r.FleetState)
	}
	return nil
}

func (r *EC2Fleet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2Fleet) String() string {
	return *r.FleetID
}

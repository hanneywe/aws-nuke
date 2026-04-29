package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EFSAccessPointResource = "EFSAccessPoint"

func init() {
	registry.Register(&registry.Registration{
		Name:     EFSAccessPointResource,
		Scope:    nuke.Account,
		Resource: &EFSAccessPoint{},
		Lister:   &EFSAccessPointLister{},
	})
}

type EFSAccessPointLister struct {
	svc EFSV2Client
}

func (l *EFSAccessPointLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = efs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := efs.NewDescribeAccessPointsPaginator(svc, &efs.DescribeAccessPointsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, accessPoint := range output.AccessPoints {
			lifeCycleState := string(accessPoint.LifeCycleState)
			tagMap := make(map[string]string)
			for _, tag := range accessPoint.Tags {
				if tag.Key != nil && tag.Value != nil {
					tagMap[*tag.Key] = *tag.Value
				}
			}
			resources = append(resources, &EFSAccessPoint{
				svc:            svc,
				AccessPointID:  accessPoint.AccessPointId,
				AccessPointArn: accessPoint.AccessPointArn,
				Name:           accessPoint.Name,
				LifeCycleState: &lifeCycleState,
				Tags:           tagMap,
			})
		}
	}

	return resources, nil
}

type EFSAccessPoint struct {
	svc            EFSV2Client
	AccessPointID  *string `property:"name=AccessPointId"`
	AccessPointArn *string
	Name           *string
	LifeCycleState *string
	Tags           map[string]string
}

func (r *EFSAccessPoint) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAccessPoint(ctx, &efs.DeleteAccessPointInput{
		AccessPointId: r.AccessPointID,
	})
	return err
}

func (r *EFSAccessPoint) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EFSAccessPoint) String() string {
	return *r.AccessPointID
}

func (r *EFSAccessPoint) Filter() error {
	if r.LifeCycleState != nil {
		state := efstypes.LifeCycleState(*r.LifeCycleState)
		if state == efstypes.LifeCycleStateDeleting || state == efstypes.LifeCycleStateDeleted {
			return fmt.Errorf("already %s", *r.LifeCycleState)
		}
	}
	return nil
}

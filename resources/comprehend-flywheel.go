package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	comprehendtypes "github.com/aws/aws-sdk-go-v2/service/comprehend/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ComprehendFlywheelResource = "ComprehendFlywheel"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComprehendFlywheelResource,
		Scope:    nuke.Account,
		Resource: &ComprehendFlywheel{},
		Lister:   &ComprehendFlywheelLister{},
	})
}

type ComprehendFlywheelLister struct {
	svc ComprehendClient
}

func (l *ComprehendFlywheelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = comprehend.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &comprehend.ListFlywheelsInput{}

	for {
		resp, err := svc.ListFlywheels(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, fw := range resp.FlywheelSummaryList {
			resources = append(resources, &ComprehendFlywheel{
				svc:         svc,
				FlywheelArn: fw.FlywheelArn,
				ModelType:   fw.ModelType,
				Status:      fw.Status,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendFlywheel struct {
	svc         ComprehendClient
	FlywheelArn *string
	ModelType   comprehendtypes.ModelType
	Status      comprehendtypes.FlywheelStatus
}

func (r *ComprehendFlywheel) Filter() error {
	if r.Status == comprehendtypes.FlywheelStatusDeleting {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *ComprehendFlywheel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFlywheel(ctx, &comprehend.DeleteFlywheelInput{
		FlywheelArn: r.FlywheelArn,
	})
	return err
}

func (r *ComprehendFlywheel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ComprehendFlywheel) String() string {
	return *r.FlywheelArn
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTBillingGroupResource = "IoTBillingGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTBillingGroupResource,
		Scope:    nuke.Account,
		Resource: &IoTBillingGroup{},
		Lister:   &IoTBillingGroupLister{},
	})
}

type IoTBillingGroupLister struct {
	svc IoTClient
}

func (l *IoTBillingGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iot.NewListBillingGroupsPaginator(svc, &iot.ListBillingGroupsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, group := range resp.BillingGroups {
			resources = append(resources, &IoTBillingGroup{
				svc:              svc,
				BillingGroupName: group.GroupName,
				BillingGroupArn:  group.GroupArn,
			})
		}
	}

	return resources, nil
}

type IoTBillingGroup struct {
	svc              IoTClient
	BillingGroupName *string
	BillingGroupArn  *string
}

func (r *IoTBillingGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBillingGroup(ctx, &iot.DeleteBillingGroupInput{
		BillingGroupName: r.BillingGroupName,
	})
	return err
}

func (r *IoTBillingGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTBillingGroup) String() string {
	return *r.BillingGroupName
}

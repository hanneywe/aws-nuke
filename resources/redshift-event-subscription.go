package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RedshiftEventSubscriptionResource = "RedshiftEventSubscription"

func init() {
	registry.Register(&registry.Registration{
		Name:     RedshiftEventSubscriptionResource,
		Scope:    nuke.Account,
		Resource: &RedshiftEventSubscription{},
		Lister:   &RedshiftEventSubscriptionLister{},
	})
}

type RedshiftEventSubscriptionLister struct {
	svc RedshiftClient
}

func (l *RedshiftEventSubscriptionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = redshift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := redshift.NewDescribeEventSubscriptionsPaginator(svc, &redshift.DescribeEventSubscriptionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range output.EventSubscriptionsList {
			eventSubscription := &output.EventSubscriptionsList[i]
			resources = append(resources, &RedshiftEventSubscription{
				svc:                svc,
				CustSubscriptionID: eventSubscription.CustSubscriptionId,
			})
		}
	}

	return resources, nil
}

type RedshiftEventSubscription struct {
	svc                RedshiftClient
	CustSubscriptionID *string `property:"name=CustSubscriptionId"`
}

func (r *RedshiftEventSubscription) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventSubscription(ctx, &redshift.DeleteEventSubscriptionInput{
		SubscriptionName: r.CustSubscriptionID,
	})
	return err
}

func (r *RedshiftEventSubscription) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RedshiftEventSubscription) String() string {
	return *r.CustSubscriptionID
}

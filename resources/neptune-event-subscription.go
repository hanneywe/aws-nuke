package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/neptune"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NeptuneEventSubscriptionResource = "NeptuneEventSubscription"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneEventSubscriptionResource,
		Scope:    nuke.Account,
		Resource: &NeptuneEventSubscription{},
		Lister:   &NeptuneEventSubscriptionLister{},
	})
}

type NeptuneEventSubscriptionLister struct {
	svc NeptuneV2Client
}

func (l *NeptuneEventSubscriptionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptune.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := neptune.NewDescribeEventSubscriptionsPaginator(svc, &neptune.DescribeEventSubscriptionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, eventSubscription := range output.EventSubscriptionsList {
			resources = append(resources, &NeptuneEventSubscription{
				svc:                  svc,
				CustSubscriptionID:   eventSubscription.CustSubscriptionId,
				EventSubscriptionArn: eventSubscription.EventSubscriptionArn,
			})
		}
	}

	return resources, nil
}

type NeptuneEventSubscription struct {
	svc NeptuneV2Client

	CustSubscriptionID   *string `property:"name=CustSubscriptionId"`
	EventSubscriptionArn *string
}

func (r *NeptuneEventSubscription) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventSubscription(ctx, &neptune.DeleteEventSubscriptionInput{
		SubscriptionName: r.CustSubscriptionID,
	})
	return err
}

func (r *NeptuneEventSubscription) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneEventSubscription) String() string {
	return *r.CustSubscriptionID
}

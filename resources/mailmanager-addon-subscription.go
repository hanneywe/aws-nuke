package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerAddonSubscriptionResource = "MailManagerAddonSubscription"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerAddonSubscriptionResource,
		Scope:    nuke.Account,
		Resource: &MailManagerAddonSubscription{},
		Lister:   &MailManagerAddonSubscriptionLister{},
	})
}

type MailManagerAddonSubscriptionLister struct {
	svc MailManagerClient
}

func (l *MailManagerAddonSubscriptionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mailmanager.NewListAddonSubscriptionsPaginator(svc, &mailmanager.ListAddonSubscriptionsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, as := range resp.AddonSubscriptions {
			resources = append(resources, &MailManagerAddonSubscription{
				svc:                 svc,
				AddonSubscriptionID: as.AddonSubscriptionId,
			})
		}
	}
	return resources, nil
}

type MailManagerAddonSubscription struct {
	svc                 MailManagerClient
	AddonSubscriptionID *string `property:"name=AddonSubscriptionId"`
}

func (r *MailManagerAddonSubscription) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAddonSubscription(ctx, &mailmanager.DeleteAddonSubscriptionInput{
		AddonSubscriptionId: r.AddonSubscriptionID,
	})
	return err
}

func (r *MailManagerAddonSubscription) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerAddonSubscription) String() string {
	return *r.AddonSubscriptionID
}

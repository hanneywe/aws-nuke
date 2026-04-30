package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EKSAnywhereSubscriptionResource = "EKSAnywhereSubscription"

func init() {
	registry.Register(&registry.Registration{
		Name:     EKSAnywhereSubscriptionResource,
		Scope:    nuke.Account,
		Resource: &EKSAnywhereSubscription{},
		Lister:   &EKSAnywhereSubscriptionLister{},
	})
}

type EKSAnywhereSubscriptionLister struct {
	svc EKSv2Client
}

func (l *EKSAnywhereSubscriptionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = eks.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &eks.ListEksAnywhereSubscriptionsInput{}
	for {
		resp, err := svc.ListEksAnywhereSubscriptions(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.Subscriptions {
			resources = append(resources, &EKSAnywhereSubscription{
				svc:    svc,
				ID:     resp.Subscriptions[i].Id,
				Arn:    resp.Subscriptions[i].Arn,
				Status: resp.Subscriptions[i].Status,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type EKSAnywhereSubscription struct {
	svc    EKSv2Client
	ID     *string `property:"name=Id"`
	Arn    *string
	Status *string
}

func (r *EKSAnywhereSubscription) Filter() error {
	if r.Status != nil && *r.Status != "EXPIRED" && *r.Status != "INACTIVE" {
		return fmt.Errorf("cannot delete subscription in %s status, must be expired or inactive", *r.Status)
	}
	return nil
}

func (r *EKSAnywhereSubscription) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEksAnywhereSubscription(ctx, &eks.DeleteEksAnywhereSubscriptionInput{
		Id: r.ID,
	})
	return err
}

func (r *EKSAnywhereSubscription) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EKSAnywhereSubscription) String() string {
	return *r.ID
}

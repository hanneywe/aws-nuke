package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/billingconductor"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BillingConductorPricingPlanResource = "BillingConductorPricingPlan"

func init() {
	registry.Register(&registry.Registration{
		Name:     BillingConductorPricingPlanResource,
		Scope:    nuke.Account,
		Resource: &BillingConductorPricingPlan{},
		Lister:   &BillingConductorPricingPlanLister{},
	})
}

type BillingConductorPricingPlanLister struct {
	svc BillingConductorClient
}

func (l *BillingConductorPricingPlanLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = billingconductor.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := billingconductor.NewListPricingPlansPaginator(svc, &billingconductor.ListPricingPlansInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.PricingPlans {
			resources = append(resources, &BillingConductorPricingPlan{
				svc:  svc,
				Arn:  item.Arn,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type BillingConductorPricingPlan struct {
	svc  BillingConductorClient
	Arn  *string
	Name *string
}

func (r *BillingConductorPricingPlan) Filter() error {
	if r.Arn != nil && strings.Contains(*r.Arn, "::aws:pricingplan/") {
		return fmt.Errorf("cannot delete AWS-managed pricing plan")
	}
	return nil
}

func (r *BillingConductorPricingPlan) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePricingPlan(ctx, &billingconductor.DeletePricingPlanInput{
		Arn: r.Arn,
	})
	return err
}

func (r *BillingConductorPricingPlan) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BillingConductorPricingPlan) String() string {
	return *r.Arn
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/billingconductor"
)

// BillingConductorClient is the interface for the BillingConductor SDK client methods.
type BillingConductorClient interface {
	ListPricingPlans(ctx context.Context, params *billingconductor.ListPricingPlansInput,
		optFns ...func(*billingconductor.Options)) (*billingconductor.ListPricingPlansOutput, error)
	DeletePricingPlan(ctx context.Context, params *billingconductor.DeletePricingPlanInput,
		optFns ...func(*billingconductor.Options)) (*billingconductor.DeletePricingPlanOutput, error)
}

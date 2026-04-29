package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/billingconductor"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockBillingConductorClient struct {
	mock.Mock
}

func (m *mockBillingConductorClient) ListPricingPlans(ctx context.Context,
	params *billingconductor.ListPricingPlansInput,
	_ ...func(*billingconductor.Options)) (*billingconductor.ListPricingPlansOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*billingconductor.ListPricingPlansOutput), args.Error(1)
}

func (m *mockBillingConductorClient) DeletePricingPlan(ctx context.Context,
	params *billingconductor.DeletePricingPlanInput,
	_ ...func(*billingconductor.Options)) (*billingconductor.DeletePricingPlanOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*billingconductor.DeletePricingPlanOutput), args.Error(1)
}

var testBillingConductorListerOpts = &nuke.ListerOpts{}

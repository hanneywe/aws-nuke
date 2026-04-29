package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/billingconductor"
	bctypes "github.com/aws/aws-sdk-go-v2/service/billingconductor/types"
)

func Test_Mock_BillingConductorPricingPlan_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBillingConductorClient)

	mockClient.
		On("ListPricingPlans", mock.Anything, mock.Anything).
		Return(&billingconductor.ListPricingPlansOutput{
			PricingPlans: []bctypes.PricingPlanListElement{
				{
					Arn:  ptr.String("arn:aws:billingconductor::123456789012:pricingplan/plan-1"),
					Name: ptr.String("my-plan"),
				},
			},
		}, nil)

	lister := &BillingConductorPricingPlanLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBillingConductorListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	plan := resources[0].(*BillingConductorPricingPlan)
	a.Equal("my-plan", *plan.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BillingConductorPricingPlan_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBillingConductorClient)

	mockClient.
		On("ListPricingPlans", mock.Anything, mock.Anything).
		Return(&billingconductor.ListPricingPlansOutput{
			PricingPlans: []bctypes.PricingPlanListElement{},
		}, nil)

	lister := &BillingConductorPricingPlanLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBillingConductorListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BillingConductorPricingPlan_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBillingConductorClient)

	plan := &BillingConductorPricingPlan{
		svc: mockClient,
		Arn: ptr.String("arn:aws:billingconductor::123456789012:pricingplan/plan-1"),
	}

	mockClient.
		On("DeletePricingPlan", mock.Anything, &billingconductor.DeletePricingPlanInput{
			Arn: plan.Arn,
		}).
		Return(&billingconductor.DeletePricingPlanOutput{}, nil)

	err := plan.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BillingConductorPricingPlan_Properties(t *testing.T) {
	a := assert.New(t)

	plan := BillingConductorPricingPlan{
		Arn:  ptr.String("arn:aws:billingconductor::123456789012:pricingplan/plan-1"),
		Name: ptr.String("my-plan"),
	}

	props := plan.Properties()
	a.Equal("arn:aws:billingconductor::123456789012:pricingplan/plan-1", props.Get("Arn"))
	a.Equal("my-plan", props.Get("Name"))
}

func Test_Mock_BillingConductorPricingPlan_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)

	plan := BillingConductorPricingPlan{
		Arn:  ptr.String("arn:aws:billingconductor::aws:pricingplan/BasicPricingPlan"),
		Name: ptr.String("BasicPricingPlan"),
	}

	err := plan.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed pricing plan")
}

func Test_Mock_BillingConductorPricingPlan_Filter_UserOwned(t *testing.T) {
	a := assert.New(t)

	plan := BillingConductorPricingPlan{
		Arn:  ptr.String("arn:aws:billingconductor::123456789012:pricingplan/plan-1"),
		Name: ptr.String("my-plan"),
	}

	err := plan.Filter()
	a.NoError(err)
}

func Test_Mock_BillingConductorPricingPlan_String(t *testing.T) {
	a := assert.New(t)

	plan := BillingConductorPricingPlan{
		Arn: ptr.String("arn:aws:billingconductor::123456789012:pricingplan/plan-1"),
	}

	a.Equal("arn:aws:billingconductor::123456789012:pricingplan/plan-1", plan.String())
}

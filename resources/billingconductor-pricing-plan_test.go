//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/billingconductor"
)

type TestBillingConductorPricingPlanSuite struct {
	suite.Suite
	svc *billingconductor.Client
	arn *string
}

func (suite *TestBillingConductorPricingPlanSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = billingconductor.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := suite.svc.CreatePricingPlan(ctx, &billingconductor.CreatePricingPlanInput{
		Name: &name,
	})
	if err != nil {
		suite.T().Fatalf("failed to create test pricing plan: %v", err)
	}
	suite.arn = resp.Arn
}

func (suite *TestBillingConductorPricingPlanSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.arn != nil {
		_, _ = suite.svc.DeletePricingPlan(ctx, &billingconductor.DeletePricingPlanInput{
			Arn: suite.arn,
		})
	}
}

func (suite *TestBillingConductorPricingPlanSuite) TestList() {
	a := assert.New(suite.T())

	lister := BillingConductorPricingPlanLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testBillingConductorListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestBillingConductorPricingPlanSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := BillingConductorPricingPlan{
		svc: suite.svc,
		Arn: suite.arn,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.arn = nil
}

func TestBillingConductorPricingPlanIntegration(t *testing.T) {
	suite.Run(t, new(TestBillingConductorPricingPlanSuite))
}

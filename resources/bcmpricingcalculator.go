package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"
)

// BCMPricingCalculatorClient is an interface for the AWS BCM Pricing Calculator SDK client methods
// used by all BCMPricingCalculator resources. It enables mock testing of List and Remove operations.
type BCMPricingCalculatorClient interface {
	ListBillEstimates(ctx context.Context, params *bcmpricingcalculator.ListBillEstimatesInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListBillEstimatesOutput, error)
	DeleteBillEstimate(ctx context.Context, params *bcmpricingcalculator.DeleteBillEstimateInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteBillEstimateOutput, error)
	ListBillScenarios(ctx context.Context, params *bcmpricingcalculator.ListBillScenariosInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListBillScenariosOutput, error)
	DeleteBillScenario(ctx context.Context, params *bcmpricingcalculator.DeleteBillScenarioInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteBillScenarioOutput, error)
	ListWorkloadEstimates(ctx context.Context, params *bcmpricingcalculator.ListWorkloadEstimatesInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListWorkloadEstimatesOutput, error)
	DeleteWorkloadEstimate(ctx context.Context, params *bcmpricingcalculator.DeleteWorkloadEstimateInput,
		optFns ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteWorkloadEstimateOutput, error)
}

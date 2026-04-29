package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"
)

type mockBCMPricingCalculatorClient struct {
	mock.Mock
}

func (m *mockBCMPricingCalculatorClient) ListBillEstimates(ctx context.Context, params *bcmpricingcalculator.ListBillEstimatesInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListBillEstimatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.ListBillEstimatesOutput), args.Error(1)
}

func (m *mockBCMPricingCalculatorClient) DeleteBillEstimate(ctx context.Context, params *bcmpricingcalculator.DeleteBillEstimateInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteBillEstimateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.DeleteBillEstimateOutput), args.Error(1)
}

func (m *mockBCMPricingCalculatorClient) ListBillScenarios(ctx context.Context, params *bcmpricingcalculator.ListBillScenariosInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListBillScenariosOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.ListBillScenariosOutput), args.Error(1)
}

func (m *mockBCMPricingCalculatorClient) DeleteBillScenario(ctx context.Context, params *bcmpricingcalculator.DeleteBillScenarioInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteBillScenarioOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.DeleteBillScenarioOutput), args.Error(1)
}

func (m *mockBCMPricingCalculatorClient) ListWorkloadEstimates(ctx context.Context, params *bcmpricingcalculator.ListWorkloadEstimatesInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.ListWorkloadEstimatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.ListWorkloadEstimatesOutput), args.Error(1)
}

func (m *mockBCMPricingCalculatorClient) DeleteWorkloadEstimate(
	ctx context.Context, params *bcmpricingcalculator.DeleteWorkloadEstimateInput,
	_ ...func(*bcmpricingcalculator.Options)) (*bcmpricingcalculator.DeleteWorkloadEstimateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bcmpricingcalculator.DeleteWorkloadEstimateOutput), args.Error(1)
}

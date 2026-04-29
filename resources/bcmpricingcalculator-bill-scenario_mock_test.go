package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"
	bcmpricingcalculatortypes "github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator/types"
)

func Test_Mock_BCMPricingCalculatorBillScenario_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListBillScenarios", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListBillScenariosOutput{
				Items: []bcmpricingcalculatortypes.BillScenarioSummary{
					{
						Id:   ptr.String("bs-12345"),
						Name: ptr.String("test-bill-scenario"),
					},
				},
			}, nil,
		)

	lister := &BCMPricingCalculatorBillScenarioLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	billScenario := resources[0].(*BCMPricingCalculatorBillScenario)
	assertions.Equal("bs-12345", *billScenario.BillScenarioID)
	assertions.Equal("test-bill-scenario", *billScenario.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillScenario_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListBillScenarios", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListBillScenariosOutput{
				Items: []bcmpricingcalculatortypes.BillScenarioSummary{},
			}, nil,
		)

	lister := &BCMPricingCalculatorBillScenarioLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillScenario_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	billScenario := &BCMPricingCalculatorBillScenario{
		svc:            mockClient,
		BillScenarioID: ptr.String("bs-12345"),
		Name:           ptr.String("test-bill-scenario"),
	}

	mockClient.
		On("DeleteBillScenario", mock.Anything, &bcmpricingcalculator.DeleteBillScenarioInput{
			Identifier: billScenario.BillScenarioID,
		}).
		Return(&bcmpricingcalculator.DeleteBillScenarioOutput{}, nil)

	err := billScenario.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillScenario_Properties(t *testing.T) {
	assertions := assert.New(t)

	billScenario := BCMPricingCalculatorBillScenario{
		BillScenarioID: ptr.String("bs-12345"),
		Name:           ptr.String("test-bill-scenario"),
	}

	properties := billScenario.Properties()
	assertions.Equal("bs-12345", properties.Get("BillScenarioId"))
	assertions.Equal("test-bill-scenario", properties.Get("Name"))
}

func Test_Mock_BCMPricingCalculatorBillScenario_String(t *testing.T) {
	assertions := assert.New(t)

	billScenario := BCMPricingCalculatorBillScenario{
		BillScenarioID: ptr.String("bs-12345"),
	}

	assertions.Equal("bs-12345", billScenario.String())
}

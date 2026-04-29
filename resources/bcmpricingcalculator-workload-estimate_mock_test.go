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

func Test_Mock_BCMPricingCalculatorWorkloadEstimate_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListWorkloadEstimates", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListWorkloadEstimatesOutput{
				Items: []bcmpricingcalculatortypes.WorkloadEstimateSummary{
					{
						Id:   ptr.String("we-12345"),
						Name: ptr.String("test-workload-estimate"),
					},
				},
			}, nil,
		)

	lister := &BCMPricingCalculatorWorkloadEstimateLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	workloadEstimate := resources[0].(*BCMPricingCalculatorWorkloadEstimate)
	assertions.Equal("we-12345", *workloadEstimate.Identifier)
	assertions.Equal("test-workload-estimate", *workloadEstimate.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorWorkloadEstimate_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListWorkloadEstimates", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListWorkloadEstimatesOutput{
				Items: []bcmpricingcalculatortypes.WorkloadEstimateSummary{},
			}, nil,
		)

	lister := &BCMPricingCalculatorWorkloadEstimateLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorWorkloadEstimate_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	workloadEstimate := &BCMPricingCalculatorWorkloadEstimate{
		svc:        mockClient,
		Identifier: ptr.String("we-12345"),
		Name:       ptr.String("test-workload-estimate"),
	}

	mockClient.
		On("DeleteWorkloadEstimate", mock.Anything, &bcmpricingcalculator.DeleteWorkloadEstimateInput{
			Identifier: workloadEstimate.Identifier,
		}).
		Return(&bcmpricingcalculator.DeleteWorkloadEstimateOutput{}, nil)

	err := workloadEstimate.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorWorkloadEstimate_Properties(t *testing.T) {
	assertions := assert.New(t)

	workloadEstimate := BCMPricingCalculatorWorkloadEstimate{
		Identifier: ptr.String("we-12345"),
		Name:       ptr.String("test-workload-estimate"),
	}

	properties := workloadEstimate.Properties()
	assertions.Equal("we-12345", properties.Get("Identifier"))
	assertions.Equal("test-workload-estimate", properties.Get("Name"))
}

func Test_Mock_BCMPricingCalculatorWorkloadEstimate_String(t *testing.T) {
	assertions := assert.New(t)

	workloadEstimate := BCMPricingCalculatorWorkloadEstimate{
		Identifier: ptr.String("we-12345"),
	}

	assertions.Equal("we-12345", workloadEstimate.String())
}

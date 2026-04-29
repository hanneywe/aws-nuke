package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"
	bcmpricingcalculatortypes "github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testBCMPricingCalculatorListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_BCMPricingCalculatorBillEstimate_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListBillEstimates", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListBillEstimatesOutput{
				Items: []bcmpricingcalculatortypes.BillEstimateSummary{
					{
						Id:     ptr.String("be-12345"),
						Name:   ptr.String("test-bill-estimate"),
						Status: bcmpricingcalculatortypes.BillEstimateStatusComplete,
					},
				},
			}, nil,
		)

	lister := &BCMPricingCalculatorBillEstimateLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	billEstimate := resources[0].(*BCMPricingCalculatorBillEstimate)
	assertions.Equal("be-12345", *billEstimate.Identifier)
	assertions.Equal("test-bill-estimate", *billEstimate.Name)
	assertions.Equal("COMPLETE", billEstimate.Status)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillEstimate_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	mockClient.
		On("ListBillEstimates", mock.Anything, mock.Anything).
		Return(
			&bcmpricingcalculator.ListBillEstimatesOutput{
				Items: []bcmpricingcalculatortypes.BillEstimateSummary{},
			}, nil,
		)

	lister := &BCMPricingCalculatorBillEstimateLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBCMPricingCalculatorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillEstimate_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBCMPricingCalculatorClient)

	billEstimate := &BCMPricingCalculatorBillEstimate{
		svc:        mockClient,
		Identifier: ptr.String("be-12345"),
		Name:       ptr.String("test-bill-estimate"),
	}

	mockClient.
		On("DeleteBillEstimate", mock.Anything, &bcmpricingcalculator.DeleteBillEstimateInput{
			Identifier: billEstimate.Identifier,
		}).
		Return(&bcmpricingcalculator.DeleteBillEstimateOutput{}, nil)

	err := billEstimate.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BCMPricingCalculatorBillEstimate_Properties(t *testing.T) {
	assertions := assert.New(t)

	billEstimate := BCMPricingCalculatorBillEstimate{
		Identifier: ptr.String("be-12345"),
		Name:       ptr.String("test-bill-estimate"),
		Status:     "COMPLETE",
	}

	properties := billEstimate.Properties()
	assertions.Equal("be-12345", properties.Get("Identifier"))
	assertions.Equal("test-bill-estimate", properties.Get("Name"))
	assertions.Equal("COMPLETE", properties.Get("Status"))
}

func Test_Mock_BCMPricingCalculatorBillEstimate_String(t *testing.T) {
	assertions := assert.New(t)

	billEstimate := BCMPricingCalculatorBillEstimate{
		Identifier: ptr.String("be-12345"),
	}

	assertions.Equal("be-12345", billEstimate.String())
}

func Test_Mock_BCMPricingCalculatorBillEstimate_Filter_InProgress(t *testing.T) {
	a := assert.New(t)

	be := &BCMPricingCalculatorBillEstimate{
		Identifier: ptr.String("be-12345"),
		Status:     string(bcmpricingcalculatortypes.BillEstimateStatusInProgress),
	}

	err := be.Filter()
	a.Error(err)
	a.Contains(err.Error(), "in progress")
}

func Test_Mock_BCMPricingCalculatorBillEstimate_Filter_Complete(t *testing.T) {
	a := assert.New(t)

	be := &BCMPricingCalculatorBillEstimate{
		Identifier: ptr.String("be-12345"),
		Status:     string(bcmpricingcalculatortypes.BillEstimateStatusComplete),
	}

	err := be.Filter()
	a.NoError(err)
}

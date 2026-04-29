package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"
	bcmtypes "github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BCMPricingCalculatorBillEstimateResource = "BCMPricingCalculatorBillEstimate"

func init() {
	registry.Register(&registry.Registration{
		Name:     BCMPricingCalculatorBillEstimateResource,
		Scope:    nuke.Account,
		Resource: &BCMPricingCalculatorBillEstimate{},
		Lister:   &BCMPricingCalculatorBillEstimateLister{},
	})
}

type BCMPricingCalculatorBillEstimateLister struct {
	svc BCMPricingCalculatorClient
}

func (l *BCMPricingCalculatorBillEstimateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = bcmpricingcalculator.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := bcmpricingcalculator.NewListBillEstimatesPaginator(svc, &bcmpricingcalculator.ListBillEstimatesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, billEstimate := range output.Items {
			resources = append(resources, &BCMPricingCalculatorBillEstimate{
				svc:        svc,
				Identifier: billEstimate.Id,
				Name:       billEstimate.Name,
				Status:     string(billEstimate.Status),
			})
		}
	}

	return resources, nil
}

type BCMPricingCalculatorBillEstimate struct {
	svc        BCMPricingCalculatorClient
	Identifier *string
	Name       *string
	Status     string
}

func (r *BCMPricingCalculatorBillEstimate) Filter() error {
	if r.Status == string(bcmtypes.BillEstimateStatusInProgress) {
		return fmt.Errorf("bill estimate is still in progress")
	}
	return nil
}

func (r *BCMPricingCalculatorBillEstimate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBillEstimate(ctx, &bcmpricingcalculator.DeleteBillEstimateInput{
		Identifier: r.Identifier,
	})
	return err
}

func (r *BCMPricingCalculatorBillEstimate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BCMPricingCalculatorBillEstimate) String() string {
	return *r.Identifier
}

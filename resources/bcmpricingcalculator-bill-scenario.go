package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BCMPricingCalculatorBillScenarioResource = "BCMPricingCalculatorBillScenario"

func init() {
	registry.Register(&registry.Registration{
		Name:     BCMPricingCalculatorBillScenarioResource,
		Scope:    nuke.Account,
		Resource: &BCMPricingCalculatorBillScenario{},
		Lister:   &BCMPricingCalculatorBillScenarioLister{},
	})
}

type BCMPricingCalculatorBillScenarioLister struct {
	svc BCMPricingCalculatorClient
}

func (l *BCMPricingCalculatorBillScenarioLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = bcmpricingcalculator.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := bcmpricingcalculator.NewListBillScenariosPaginator(svc, &bcmpricingcalculator.ListBillScenariosInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, billScenario := range output.Items {
			resources = append(resources, &BCMPricingCalculatorBillScenario{
				svc:            svc,
				BillScenarioID: billScenario.Id,
				Name:           billScenario.Name,
			})
		}
	}

	return resources, nil
}

type BCMPricingCalculatorBillScenario struct {
	svc            BCMPricingCalculatorClient
	BillScenarioID *string `property:"name=BillScenarioId"`
	Name           *string
}

func (r *BCMPricingCalculatorBillScenario) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBillScenario(ctx, &bcmpricingcalculator.DeleteBillScenarioInput{
		Identifier: r.BillScenarioID,
	})
	return err
}

func (r *BCMPricingCalculatorBillScenario) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BCMPricingCalculatorBillScenario) String() string {
	return *r.BillScenarioID
}

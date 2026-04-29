package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bcmpricingcalculator"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BCMPricingCalculatorWorkloadEstimateResource = "BCMPricingCalculatorWorkloadEstimate"

func init() {
	registry.Register(&registry.Registration{
		Name:     BCMPricingCalculatorWorkloadEstimateResource,
		Scope:    nuke.Account,
		Resource: &BCMPricingCalculatorWorkloadEstimate{},
		Lister:   &BCMPricingCalculatorWorkloadEstimateLister{},
	})
}

type BCMPricingCalculatorWorkloadEstimateLister struct {
	svc BCMPricingCalculatorClient
}

func (l *BCMPricingCalculatorWorkloadEstimateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = bcmpricingcalculator.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := bcmpricingcalculator.NewListWorkloadEstimatesPaginator(svc, &bcmpricingcalculator.ListWorkloadEstimatesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, workloadEstimate := range output.Items {
			resources = append(resources, &BCMPricingCalculatorWorkloadEstimate{
				svc:        svc,
				Identifier: workloadEstimate.Id,
				Name:       workloadEstimate.Name,
			})
		}
	}

	return resources, nil
}

type BCMPricingCalculatorWorkloadEstimate struct {
	svc        BCMPricingCalculatorClient
	Identifier *string
	Name       *string
}

func (r *BCMPricingCalculatorWorkloadEstimate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkloadEstimate(ctx, &bcmpricingcalculator.DeleteWorkloadEstimateInput{
		Identifier: r.Identifier,
	})
	return err
}

func (r *BCMPricingCalculatorWorkloadEstimate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BCMPricingCalculatorWorkloadEstimate) String() string {
	return *r.Identifier
}

package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/snowball"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SnowballLongTermPricingResource = "SnowballLongTermPricing"

func init() {
	registry.Register(&registry.Registration{
		Name:     SnowballLongTermPricingResource,
		Scope:    nuke.Account,
		Resource: &SnowballLongTermPricing{},
		Lister:   &SnowballLongTermPricingLister{},
	})
}

type SnowballLongTermPricingLister struct {
	svc SnowballClient
}

func (l *SnowballLongTermPricingLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = snowball.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := snowball.NewListLongTermPricingPaginator(svc, &snowball.ListLongTermPricingInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range page.LongTermPricingEntries {
			entry := &page.LongTermPricingEntries[i]
			if entry.IsLongTermPricingAutoRenew != nil && *entry.IsLongTermPricingAutoRenew {
				resources = append(resources, &SnowballLongTermPricing{
					svc:               svc,
					LongTermPricingID: entry.LongTermPricingId,
					Status:            entry.LongTermPricingStatus,
					StartDate:         entry.LongTermPricingStartDate,
					EndDate:           entry.LongTermPricingEndDate,
				})
			}
		}
	}

	return resources, nil
}

type SnowballLongTermPricing struct {
	svc               SnowballClient
	LongTermPricingID *string
	Status            *string
	StartDate         *time.Time
	EndDate           *time.Time
}

func (r *SnowballLongTermPricing) Remove(ctx context.Context) error {
	_, err := r.svc.UpdateLongTermPricing(ctx, &snowball.UpdateLongTermPricingInput{
		LongTermPricingId:          r.LongTermPricingID,
		IsLongTermPricingAutoRenew: aws.Bool(false),
	})
	return err
}

func (r *SnowballLongTermPricing) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SnowballLongTermPricing) String() string {
	return *r.LongTermPricingID
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/snowball"
)

type SnowballClient interface {
	ListLongTermPricing(ctx context.Context, params *snowball.ListLongTermPricingInput,
		optFns ...func(*snowball.Options)) (*snowball.ListLongTermPricingOutput, error)
	UpdateLongTermPricing(ctx context.Context, params *snowball.UpdateLongTermPricingInput,
		optFns ...func(*snowball.Options)) (*snowball.UpdateLongTermPricingOutput, error)
}

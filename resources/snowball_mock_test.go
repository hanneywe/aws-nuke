package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/snowball"
)

type mockSnowballClient struct {
	mock.Mock
}

func (m *mockSnowballClient) ListLongTermPricing(ctx context.Context, params *snowball.ListLongTermPricingInput,
	_ ...func(*snowball.Options)) (*snowball.ListLongTermPricingOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*snowball.ListLongTermPricingOutput), args.Error(1)
}

func (m *mockSnowballClient) UpdateLongTermPricing(ctx context.Context, params *snowball.UpdateLongTermPricingInput,
	_ ...func(*snowball.Options)) (*snowball.UpdateLongTermPricingOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*snowball.UpdateLongTermPricingOutput), args.Error(1)
}

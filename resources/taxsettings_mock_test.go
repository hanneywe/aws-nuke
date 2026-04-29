package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/taxsettings"
)

type mockTaxSettingsClient struct {
	mock.Mock
}

func (m *mockTaxSettingsClient) GetTaxInheritance(
	ctx context.Context, params *taxsettings.GetTaxInheritanceInput,
	_ ...func(*taxsettings.Options),
) (*taxsettings.GetTaxInheritanceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*taxsettings.GetTaxInheritanceOutput), args.Error(1)
}

func (m *mockTaxSettingsClient) PutTaxInheritance(
	ctx context.Context, params *taxsettings.PutTaxInheritanceInput,
	_ ...func(*taxsettings.Options),
) (*taxsettings.PutTaxInheritanceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*taxsettings.PutTaxInheritanceOutput), args.Error(1)
}

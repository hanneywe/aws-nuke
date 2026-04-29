package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMigrationHubStrategyClient struct {
	mock.Mock
}

func (m *mockMigrationHubStrategyClient) GetPortfolioPreferences(
	ctx context.Context, params *migrationhubstrategy.GetPortfolioPreferencesInput,
	_ ...func(*migrationhubstrategy.Options),
) (*migrationhubstrategy.GetPortfolioPreferencesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*migrationhubstrategy.GetPortfolioPreferencesOutput), args.Error(1)
}

func (m *mockMigrationHubStrategyClient) PutPortfolioPreferences(
	ctx context.Context, params *migrationhubstrategy.PutPortfolioPreferencesInput,
	_ ...func(*migrationhubstrategy.Options),
) (*migrationhubstrategy.PutPortfolioPreferencesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*migrationhubstrategy.PutPortfolioPreferencesOutput), args.Error(1)
}

var testMigrationHubStrategyListerOpts = &nuke.ListerOpts{}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy"
)

type MigrationHubStrategyClient interface {
	GetPortfolioPreferences(ctx context.Context, params *migrationhubstrategy.GetPortfolioPreferencesInput,
		optFns ...func(*migrationhubstrategy.Options)) (*migrationhubstrategy.GetPortfolioPreferencesOutput, error)
	PutPortfolioPreferences(ctx context.Context, params *migrationhubstrategy.PutPortfolioPreferencesInput,
		optFns ...func(*migrationhubstrategy.Options)) (*migrationhubstrategy.PutPortfolioPreferencesOutput, error)
}

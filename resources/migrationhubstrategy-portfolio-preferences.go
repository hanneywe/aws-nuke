package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy"
	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libTypes "github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MigrationHubStrategyPortfolioPreferencesResource = "MigrationHubStrategyPortfolioPreferences"

func init() {
	registry.Register(&registry.Registration{
		Name:     MigrationHubStrategyPortfolioPreferencesResource,
		Scope:    nuke.Account,
		Resource: &MigrationHubStrategyPortfolioPreferences{},
		Lister:   &MigrationHubStrategyPortfolioPreferencesLister{},
	})
}

type MigrationHubStrategyPortfolioPreferencesLister struct {
	svc MigrationHubStrategyClient
}

func (l *MigrationHubStrategyPortfolioPreferencesLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = migrationhubstrategy.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetPortfolioPreferences(ctx, &migrationhubstrategy.GetPortfolioPreferencesInput{})
	if err != nil {
		return nil, err
	}

	// Portfolio preferences cannot be deleted via the API, only updated.
	// Once the service has been used, GetPortfolioPreferences always returns
	// non-nil sub-preferences. The Remove method resets ApplicationMode to ALL,
	// so skip listing when the mode is already ALL (the default/reset state)
	// to avoid an infinite removal loop.
	if resp.ApplicationMode == "" || resp.ApplicationMode == types.ApplicationModeAll {
		return nil, nil
	}

	return []resource.Resource{
		&MigrationHubStrategyPortfolioPreferences{
			svc:             svc,
			ApplicationMode: resp.ApplicationMode,
		},
	}, nil
}

type MigrationHubStrategyPortfolioPreferences struct {
	svc             MigrationHubStrategyClient
	ApplicationMode types.ApplicationMode `property:"name=ApplicationMode"`
}

func (r *MigrationHubStrategyPortfolioPreferences) Remove(ctx context.Context) error {
	_, err := r.svc.PutPortfolioPreferences(ctx, &migrationhubstrategy.PutPortfolioPreferencesInput{
		ApplicationMode: types.ApplicationModeAll,
	})
	return err
}

func (r *MigrationHubStrategyPortfolioPreferences) Properties() libTypes.Properties {
	return libTypes.NewPropertiesFromStruct(r)
}

func (r *MigrationHubStrategyPortfolioPreferences) String() string {
	return "MigrationHubStrategyPortfolioPreferences"
}

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy"
	"github.com/aws/aws-sdk-go-v2/service/migrationhubstrategy/types"
)

func Test_Mock_MigrationHubStrategyPortfolioPreferences_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubStrategyClient)

	mockClient.On("GetPortfolioPreferences", mock.Anything, mock.Anything).
		Return(&migrationhubstrategy.GetPortfolioPreferencesOutput{
			ApplicationMode:        types.ApplicationModeKnown,
			ApplicationPreferences: &types.ApplicationPreferences{},
		}, nil)

	lister := &MigrationHubStrategyPortfolioPreferencesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMigrationHubStrategyListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MigrationHubStrategyPortfolioPreferences)
	a.Equal(types.ApplicationModeKnown, r.ApplicationMode)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubStrategyPortfolioPreferences_List_DefaultState(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubStrategyClient)

	mockClient.On("GetPortfolioPreferences", mock.Anything, mock.Anything).
		Return(&migrationhubstrategy.GetPortfolioPreferencesOutput{
			ApplicationMode:        types.ApplicationModeAll,
			ApplicationPreferences: &types.ApplicationPreferences{},
		}, nil)

	lister := &MigrationHubStrategyPortfolioPreferencesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMigrationHubStrategyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubStrategyPortfolioPreferences_List_NoConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubStrategyClient)

	mockClient.On("GetPortfolioPreferences", mock.Anything, mock.Anything).
		Return(&migrationhubstrategy.GetPortfolioPreferencesOutput{}, nil)

	lister := &MigrationHubStrategyPortfolioPreferencesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMigrationHubStrategyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubStrategyPortfolioPreferences_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMigrationHubStrategyClient)

	r := &MigrationHubStrategyPortfolioPreferences{
		svc:             mockClient,
		ApplicationMode: types.ApplicationModeAll,
	}

	mockClient.On("PutPortfolioPreferences", mock.Anything,
		&migrationhubstrategy.PutPortfolioPreferencesInput{
			ApplicationMode: types.ApplicationModeAll,
		}).
		Return(&migrationhubstrategy.PutPortfolioPreferencesOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MigrationHubStrategyPortfolioPreferences_Properties(t *testing.T) {
	a := assert.New(t)
	r := &MigrationHubStrategyPortfolioPreferences{
		ApplicationMode: types.ApplicationModeAll,
	}
	props := r.Properties()
	a.Equal("ALL", props.Get("ApplicationMode"))
}

func Test_Mock_MigrationHubStrategyPortfolioPreferences_String(t *testing.T) {
	a := assert.New(t)
	r := &MigrationHubStrategyPortfolioPreferences{
		ApplicationMode: types.ApplicationModeAll,
	}
	a.Equal("MigrationHubStrategyPortfolioPreferences", r.String())
}

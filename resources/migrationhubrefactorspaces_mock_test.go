package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMigrationHubRefactorSpacesClient struct {
	mock.Mock
}

func (m *mockMigrationHubRefactorSpacesClient) ListEnvironments(
	ctx context.Context, params *migrationhubrefactorspaces.ListEnvironmentsInput,
	_ ...func(*migrationhubrefactorspaces.Options),
) (*migrationhubrefactorspaces.ListEnvironmentsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*migrationhubrefactorspaces.ListEnvironmentsOutput), args.Error(1)
}

func (m *mockMigrationHubRefactorSpacesClient) DeleteEnvironment(
	ctx context.Context, params *migrationhubrefactorspaces.DeleteEnvironmentInput,
	_ ...func(*migrationhubrefactorspaces.Options),
) (*migrationhubrefactorspaces.DeleteEnvironmentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*migrationhubrefactorspaces.DeleteEnvironmentOutput), args.Error(1)
}

var testMigrationHubRefactorSpacesListerOpts = &nuke.ListerOpts{}

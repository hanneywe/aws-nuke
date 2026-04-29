package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces"
)

// MigrationHubRefactorSpacesClient is the interface for the Migration Hub Refactor Spaces SDK client methods.
type MigrationHubRefactorSpacesClient interface {
	ListEnvironments(ctx context.Context, params *migrationhubrefactorspaces.ListEnvironmentsInput,
		optFns ...func(*migrationhubrefactorspaces.Options)) (*migrationhubrefactorspaces.ListEnvironmentsOutput, error)
	DeleteEnvironment(ctx context.Context, params *migrationhubrefactorspaces.DeleteEnvironmentInput,
		optFns ...func(*migrationhubrefactorspaces.Options)) (*migrationhubrefactorspaces.DeleteEnvironmentOutput, error)
}

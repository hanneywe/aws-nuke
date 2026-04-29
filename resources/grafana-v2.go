package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/grafana"
)

type GrafanaV2Client interface {
	ListWorkspaces(ctx context.Context, params *grafana.ListWorkspacesInput,
		optFns ...func(*grafana.Options)) (*grafana.ListWorkspacesOutput, error)
	ListWorkspaceServiceAccounts(ctx context.Context, params *grafana.ListWorkspaceServiceAccountsInput,
		optFns ...func(*grafana.Options)) (*grafana.ListWorkspaceServiceAccountsOutput, error)
	ListWorkspaceServiceAccountTokens(ctx context.Context, params *grafana.ListWorkspaceServiceAccountTokensInput,
		optFns ...func(*grafana.Options)) (*grafana.ListWorkspaceServiceAccountTokensOutput, error)
	DeleteWorkspaceServiceAccountToken(ctx context.Context, params *grafana.DeleteWorkspaceServiceAccountTokenInput,
		optFns ...func(*grafana.Options)) (*grafana.DeleteWorkspaceServiceAccountTokenOutput, error)
	DeleteWorkspaceServiceAccount(ctx context.Context, params *grafana.DeleteWorkspaceServiceAccountInput,
		optFns ...func(*grafana.Options)) (*grafana.DeleteWorkspaceServiceAccountOutput, error)
}

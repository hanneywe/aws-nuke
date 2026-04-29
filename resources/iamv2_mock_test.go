package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIAMClient struct {
	mock.Mock
}

func (m *mockIAMClient) ListAccountAliases(
	ctx context.Context, params *iam.ListAccountAliasesInput,
	_ ...func(*iam.Options),
) (*iam.ListAccountAliasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iam.ListAccountAliasesOutput), args.Error(1)
}

func (m *mockIAMClient) DeleteAccountAlias(
	ctx context.Context, params *iam.DeleteAccountAliasInput,
	_ ...func(*iam.Options),
) (*iam.DeleteAccountAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iam.DeleteAccountAliasOutput), args.Error(1)
}

var testIamListerOpts = &nuke.ListerOpts{}

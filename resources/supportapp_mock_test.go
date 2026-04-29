package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/supportapp"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSupportAppClient struct {
	mock.Mock
}

func (m *mockSupportAppClient) GetAccountAlias(
	ctx context.Context, params *supportapp.GetAccountAliasInput,
	_ ...func(*supportapp.Options),
) (*supportapp.GetAccountAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*supportapp.GetAccountAliasOutput), args.Error(1)
}

func (m *mockSupportAppClient) DeleteAccountAlias(
	ctx context.Context, params *supportapp.DeleteAccountAliasInput,
	_ ...func(*supportapp.Options),
) (*supportapp.DeleteAccountAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*supportapp.DeleteAccountAliasOutput), args.Error(1)
}

var testSupportAppListerOpts = &nuke.ListerOpts{}

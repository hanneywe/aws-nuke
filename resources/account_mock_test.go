package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/account"
)

type mockAccountClient struct {
	mock.Mock
}

func (m *mockAccountClient) GetAlternateContact(ctx context.Context, params *account.GetAlternateContactInput,
	_ ...func(*account.Options)) (*account.GetAlternateContactOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*account.GetAlternateContactOutput), args.Error(1)
}

func (m *mockAccountClient) DeleteAlternateContact(ctx context.Context, params *account.DeleteAlternateContactInput,
	_ ...func(*account.Options)) (*account.DeleteAlternateContactOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*account.DeleteAlternateContactOutput), args.Error(1)
}

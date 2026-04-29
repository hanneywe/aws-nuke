package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testPaymentCryptographyListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockPaymentCryptographyClient struct {
	mock.Mock
}

func (m *mockPaymentCryptographyClient) ListKeys(ctx context.Context, params *paymentcryptography.ListKeysInput,
	_ ...func(*paymentcryptography.Options)) (*paymentcryptography.ListKeysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*paymentcryptography.ListKeysOutput), args.Error(1)
}

func (m *mockPaymentCryptographyClient) DeleteKey(ctx context.Context, params *paymentcryptography.DeleteKeyInput,
	_ ...func(*paymentcryptography.Options)) (*paymentcryptography.DeleteKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*paymentcryptography.DeleteKeyOutput), args.Error(1)
}

func (m *mockPaymentCryptographyClient) ListAliases(ctx context.Context, params *paymentcryptography.ListAliasesInput,
	_ ...func(*paymentcryptography.Options)) (*paymentcryptography.ListAliasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*paymentcryptography.ListAliasesOutput), args.Error(1)
}

func (m *mockPaymentCryptographyClient) DeleteAlias(ctx context.Context, params *paymentcryptography.DeleteAliasInput,
	_ ...func(*paymentcryptography.Options)) (*paymentcryptography.DeleteAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*paymentcryptography.DeleteAliasOutput), args.Error(1)
}

package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/xray"
)

type mockXRayClient struct {
	mock.Mock
}

func (m *mockXRayClient) GetEncryptionConfig(
	ctx context.Context, params *xray.GetEncryptionConfigInput,
	_ ...func(*xray.Options),
) (*xray.GetEncryptionConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*xray.GetEncryptionConfigOutput), args.Error(1)
}

func (m *mockXRayClient) PutEncryptionConfig(
	ctx context.Context, params *xray.PutEncryptionConfigInput,
	_ ...func(*xray.Options),
) (*xray.PutEncryptionConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*xray.PutEncryptionConfigOutput), args.Error(1)
}

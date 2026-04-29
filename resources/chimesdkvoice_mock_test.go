package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkvoice"
)

type mockChimeSDKVoiceClient struct {
	mock.Mock
}

func (m *mockChimeSDKVoiceClient) ListPhoneNumbers(ctx context.Context, params *chimesdkvoice.ListPhoneNumbersInput,
	_ ...func(*chimesdkvoice.Options)) (*chimesdkvoice.ListPhoneNumbersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chimesdkvoice.ListPhoneNumbersOutput), args.Error(1)
}

func (m *mockChimeSDKVoiceClient) DeletePhoneNumber(ctx context.Context, params *chimesdkvoice.DeletePhoneNumberInput,
	_ ...func(*chimesdkvoice.Options)) (*chimesdkvoice.DeletePhoneNumberOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chimesdkvoice.DeletePhoneNumberOutput), args.Error(1)
}

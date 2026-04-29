package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/signer"
)

type mockSignerClient struct {
	mock.Mock
}

func (m *mockSignerClient) ListSigningProfiles(ctx context.Context, params *signer.ListSigningProfilesInput,
	_ ...func(*signer.Options)) (*signer.ListSigningProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*signer.ListSigningProfilesOutput), args.Error(1)
}

func (m *mockSignerClient) CancelSigningProfile(ctx context.Context, params *signer.CancelSigningProfileInput,
	_ ...func(*signer.Options)) (*signer.CancelSigningProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*signer.CancelSigningProfileOutput), args.Error(1)
}

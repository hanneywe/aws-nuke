package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/b2bi"
)

type mockB2BIClient struct {
	mock.Mock
}

func (m *mockB2BIClient) ListProfiles(ctx context.Context, params *b2bi.ListProfilesInput,
	_ ...func(*b2bi.Options)) (*b2bi.ListProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*b2bi.ListProfilesOutput), args.Error(1)
}

func (m *mockB2BIClient) DeleteProfile(ctx context.Context, params *b2bi.DeleteProfileInput,
	_ ...func(*b2bi.Options)) (*b2bi.DeleteProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*b2bi.DeleteProfileOutput), args.Error(1)
}

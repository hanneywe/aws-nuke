package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/healthlake"
)

type mockHealthLakeClient struct {
	mock.Mock
}

func (m *mockHealthLakeClient) ListFHIRDatastores(ctx context.Context, params *healthlake.ListFHIRDatastoresInput,
	_ ...func(*healthlake.Options)) (*healthlake.ListFHIRDatastoresOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*healthlake.ListFHIRDatastoresOutput), args.Error(1)
}

func (m *mockHealthLakeClient) DeleteFHIRDatastore(ctx context.Context, params *healthlake.DeleteFHIRDatastoreInput,
	_ ...func(*healthlake.Options)) (*healthlake.DeleteFHIRDatastoreOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*healthlake.DeleteFHIRDatastoreOutput), args.Error(1)
}

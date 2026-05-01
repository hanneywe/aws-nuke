package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3control"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockS3AccountClient struct {
	mock.Mock
}

func (m *mockS3AccountClient) GetPublicAccessBlock(
	ctx context.Context, params *s3control.GetPublicAccessBlockInput,
	_ ...func(*s3control.Options),
) (*s3control.GetPublicAccessBlockOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3control.GetPublicAccessBlockOutput), args.Error(1)
}

func (m *mockS3AccountClient) DeletePublicAccessBlock(
	ctx context.Context, params *s3control.DeletePublicAccessBlockInput,
	_ ...func(*s3control.Options),
) (*s3control.DeletePublicAccessBlockOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3control.DeletePublicAccessBlockOutput), args.Error(1)
}

var testS3AccountListerOpts = &nuke.ListerOpts{
	AccountID: func() *string { s := "123456789012"; return &s }(),
}

package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/emr"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockEMRV2Client struct {
	mock.Mock
}

func (m *mockEMRV2Client) GetBlockPublicAccessConfiguration(
	ctx context.Context, params *emr.GetBlockPublicAccessConfigurationInput,
	_ ...func(*emr.Options),
) (*emr.GetBlockPublicAccessConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*emr.GetBlockPublicAccessConfigurationOutput), args.Error(1)
}

func (m *mockEMRV2Client) PutBlockPublicAccessConfiguration(
	ctx context.Context, params *emr.PutBlockPublicAccessConfigurationInput,
	_ ...func(*emr.Options),
) (*emr.PutBlockPublicAccessConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*emr.PutBlockPublicAccessConfigurationOutput), args.Error(1)
}

var testEMRV2ListerOpts = &nuke.ListerOpts{}

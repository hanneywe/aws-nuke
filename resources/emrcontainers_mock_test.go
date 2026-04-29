package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/emrcontainers"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockEMRContainersClient struct {
	mock.Mock
}

func (m *mockEMRContainersClient) ListJobTemplates(
	ctx context.Context, params *emrcontainers.ListJobTemplatesInput,
	_ ...func(*emrcontainers.Options),
) (*emrcontainers.ListJobTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*emrcontainers.ListJobTemplatesOutput), args.Error(1)
}

func (m *mockEMRContainersClient) DeleteJobTemplate(
	ctx context.Context, params *emrcontainers.DeleteJobTemplateInput,
	_ ...func(*emrcontainers.Options),
) (*emrcontainers.DeleteJobTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*emrcontainers.DeleteJobTemplateOutput), args.Error(1)
}

var testEMRContainersListerOpts = &nuke.ListerOpts{}

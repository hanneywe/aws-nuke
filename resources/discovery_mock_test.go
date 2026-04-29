package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockApplicationdiscoveryserviceClient struct {
	mock.Mock
}

func (m *mockApplicationdiscoveryserviceClient) ListConfigurations(
	ctx context.Context, params *applicationdiscoveryservice.ListConfigurationsInput,
	_ ...func(*applicationdiscoveryservice.Options),
) (*applicationdiscoveryservice.ListConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*applicationdiscoveryservice.ListConfigurationsOutput), args.Error(1)
}

func (m *mockApplicationdiscoveryserviceClient) DeleteApplications(
	ctx context.Context, params *applicationdiscoveryservice.DeleteApplicationsInput,
	_ ...func(*applicationdiscoveryservice.Options),
) (*applicationdiscoveryservice.DeleteApplicationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*applicationdiscoveryservice.DeleteApplicationsOutput), args.Error(1)
}

var testApplicationdiscoveryserviceListerOpts = &nuke.ListerOpts{}

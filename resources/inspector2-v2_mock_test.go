package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockInspector2V2Client struct {
	mock.Mock
}

func (m *mockInspector2V2Client) ListCodeSecurityScanConfigurations(
	ctx context.Context, params *inspector2.ListCodeSecurityScanConfigurationsInput,
	_ ...func(*inspector2.Options),
) (*inspector2.ListCodeSecurityScanConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*inspector2.ListCodeSecurityScanConfigurationsOutput), args.Error(1)
}

func (m *mockInspector2V2Client) DeleteCodeSecurityScanConfiguration(
	ctx context.Context, params *inspector2.DeleteCodeSecurityScanConfigurationInput,
	_ ...func(*inspector2.Options),
) (*inspector2.DeleteCodeSecurityScanConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*inspector2.DeleteCodeSecurityScanConfigurationOutput), args.Error(1)
}

var testInspector2V2ListerOpts = &nuke.ListerOpts{}

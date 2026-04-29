package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testWorkSpacesWebListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockWorkSpacesWebClient struct {
	mock.Mock
}

func (m *mockWorkSpacesWebClient) ListPortals(ctx context.Context, params *workspacesweb.ListPortalsInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.ListPortalsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.ListPortalsOutput), args.Error(1)
}

func (m *mockWorkSpacesWebClient) ListIpAccessSettings(ctx context.Context, params *workspacesweb.ListIpAccessSettingsInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.ListIpAccessSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.ListIpAccessSettingsOutput), args.Error(1)
}

func (m *mockWorkSpacesWebClient) DeletePortal(ctx context.Context, params *workspacesweb.DeletePortalInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.DeletePortalOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.DeletePortalOutput), args.Error(1)
}

func (m *mockWorkSpacesWebClient) DeleteIpAccessSettings(ctx context.Context, params *workspacesweb.DeleteIpAccessSettingsInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.DeleteIpAccessSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.DeleteIpAccessSettingsOutput), args.Error(1)
}

func (m *mockWorkSpacesWebClient) ListDataProtectionSettings(ctx context.Context, params *workspacesweb.ListDataProtectionSettingsInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.ListDataProtectionSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.ListDataProtectionSettingsOutput), args.Error(1)
}

func (m *mockWorkSpacesWebClient) DeleteDataProtectionSettings(ctx context.Context, params *workspacesweb.DeleteDataProtectionSettingsInput,
	_ ...func(*workspacesweb.Options)) (*workspacesweb.DeleteDataProtectionSettingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*workspacesweb.DeleteDataProtectionSettingsOutput), args.Error(1)
}

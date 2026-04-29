package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"
)

// WorkSpacesWebClient is an interface for the WorkSpaces Web SDK client methods used by all WorkSpaces Web resources.
// It enables mock testing of List and Remove operations.
type WorkSpacesWebClient interface {
	// Listing
	ListPortals(ctx context.Context, params *workspacesweb.ListPortalsInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.ListPortalsOutput, error)
	ListIpAccessSettings(ctx context.Context, params *workspacesweb.ListIpAccessSettingsInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.ListIpAccessSettingsOutput, error)

	// Deletion
	DeletePortal(ctx context.Context, params *workspacesweb.DeletePortalInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.DeletePortalOutput, error)
	DeleteIpAccessSettings(ctx context.Context, params *workspacesweb.DeleteIpAccessSettingsInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.DeleteIpAccessSettingsOutput, error)
	ListDataProtectionSettings(ctx context.Context, params *workspacesweb.ListDataProtectionSettingsInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.ListDataProtectionSettingsOutput, error)
	DeleteDataProtectionSettings(ctx context.Context, params *workspacesweb.DeleteDataProtectionSettingsInput,
		optFns ...func(*workspacesweb.Options)) (*workspacesweb.DeleteDataProtectionSettingsOutput, error)
}

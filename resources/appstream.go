package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
)

// AppStreamClient is the interface for the AppStream SDK v2 client methods.
type AppStreamClient interface {
	DescribeUsers(ctx context.Context, params *appstream.DescribeUsersInput,
		optFns ...func(*appstream.Options)) (*appstream.DescribeUsersOutput, error)
	DeleteUser(ctx context.Context, params *appstream.DeleteUserInput,
		optFns ...func(*appstream.Options)) (*appstream.DeleteUserOutput, error)
	DescribeStacks(ctx context.Context, params *appstream.DescribeStacksInput,
		optFns ...func(*appstream.Options)) (*appstream.DescribeStacksOutput, error)
	DescribeEntitlements(ctx context.Context, params *appstream.DescribeEntitlementsInput,
		optFns ...func(*appstream.Options)) (*appstream.DescribeEntitlementsOutput, error)
	DeleteEntitlement(ctx context.Context, params *appstream.DeleteEntitlementInput,
		optFns ...func(*appstream.Options)) (*appstream.DeleteEntitlementOutput, error)
}

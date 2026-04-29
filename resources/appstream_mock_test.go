package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appstream"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAppStreamClient struct {
	mock.Mock
}

func (m *mockAppStreamClient) DescribeUsers(
	ctx context.Context, params *appstream.DescribeUsersInput,
	_ ...func(*appstream.Options),
) (*appstream.DescribeUsersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appstream.DescribeUsersOutput), args.Error(1)
}

func (m *mockAppStreamClient) DeleteUser(
	ctx context.Context, params *appstream.DeleteUserInput,
	_ ...func(*appstream.Options),
) (*appstream.DeleteUserOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appstream.DeleteUserOutput), args.Error(1)
}

func (m *mockAppStreamClient) DescribeStacks(
	ctx context.Context, params *appstream.DescribeStacksInput,
	_ ...func(*appstream.Options),
) (*appstream.DescribeStacksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appstream.DescribeStacksOutput), args.Error(1)
}

func (m *mockAppStreamClient) DescribeEntitlements(
	ctx context.Context, params *appstream.DescribeEntitlementsInput,
	_ ...func(*appstream.Options),
) (*appstream.DescribeEntitlementsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appstream.DescribeEntitlementsOutput), args.Error(1)
}

func (m *mockAppStreamClient) DeleteEntitlement(
	ctx context.Context, params *appstream.DeleteEntitlementInput,
	_ ...func(*appstream.Options),
) (*appstream.DeleteEntitlementOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appstream.DeleteEntitlementOutput), args.Error(1)
}

var testAppStreamListerOpts = &nuke.ListerOpts{}

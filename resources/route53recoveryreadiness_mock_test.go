package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockRoute53RecoveryReadinessClient struct {
	mock.Mock
}

func (m *mockRoute53RecoveryReadinessClient) ListCells(
	ctx context.Context, params *route53recoveryreadiness.ListCellsInput,
	_ ...func(*route53recoveryreadiness.Options),
) (*route53recoveryreadiness.ListCellsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoveryreadiness.ListCellsOutput), args.Error(1)
}

func (m *mockRoute53RecoveryReadinessClient) DeleteCell(
	ctx context.Context, params *route53recoveryreadiness.DeleteCellInput,
	_ ...func(*route53recoveryreadiness.Options),
) (*route53recoveryreadiness.DeleteCellOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoveryreadiness.DeleteCellOutput), args.Error(1)
}

func (m *mockRoute53RecoveryReadinessClient) ListRecoveryGroups(
	ctx context.Context, params *route53recoveryreadiness.ListRecoveryGroupsInput,
	_ ...func(*route53recoveryreadiness.Options),
) (*route53recoveryreadiness.ListRecoveryGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoveryreadiness.ListRecoveryGroupsOutput), args.Error(1)
}

func (m *mockRoute53RecoveryReadinessClient) DeleteRecoveryGroup(
	ctx context.Context, params *route53recoveryreadiness.DeleteRecoveryGroupInput,
	_ ...func(*route53recoveryreadiness.Options),
) (*route53recoveryreadiness.DeleteRecoveryGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53recoveryreadiness.DeleteRecoveryGroupOutput), args.Error(1)
}

var testRoute53RecoveryReadinessListerOpts = &nuke.ListerOpts{}

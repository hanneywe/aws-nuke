package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIVSRealtimeClient struct {
	mock.Mock
}

func (m *mockIVSRealtimeClient) ListEncoderConfigurations(ctx context.Context,
	params *ivsrealtime.ListEncoderConfigurationsInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.ListEncoderConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.ListEncoderConfigurationsOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) DeleteEncoderConfiguration(ctx context.Context,
	params *ivsrealtime.DeleteEncoderConfigurationInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteEncoderConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.DeleteEncoderConfigurationOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) ListStages(ctx context.Context,
	params *ivsrealtime.ListStagesInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.ListStagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.ListStagesOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) DeleteStage(ctx context.Context,
	params *ivsrealtime.DeleteStageInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteStageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.DeleteStageOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) ListIngestConfigurations(ctx context.Context,
	params *ivsrealtime.ListIngestConfigurationsInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.ListIngestConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.ListIngestConfigurationsOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) DeleteIngestConfiguration(ctx context.Context,
	params *ivsrealtime.DeleteIngestConfigurationInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteIngestConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.DeleteIngestConfigurationOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) ListStorageConfigurations(ctx context.Context,
	params *ivsrealtime.ListStorageConfigurationsInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.ListStorageConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.ListStorageConfigurationsOutput), args.Error(1)
}

func (m *mockIVSRealtimeClient) DeleteStorageConfiguration(ctx context.Context,
	params *ivsrealtime.DeleteStorageConfigurationInput,
	_ ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteStorageConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ivsrealtime.DeleteStorageConfigurationOutput), args.Error(1)
}

var testIVSRealtimeListerOpts = &nuke.ListerOpts{}

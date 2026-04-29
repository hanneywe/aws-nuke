package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"
)

// IVSRealtimeClient is the interface for the IVS Realtime SDK client methods.
type IVSRealtimeClient interface {
	ListEncoderConfigurations(ctx context.Context, params *ivsrealtime.ListEncoderConfigurationsInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.ListEncoderConfigurationsOutput, error)
	DeleteEncoderConfiguration(ctx context.Context, params *ivsrealtime.DeleteEncoderConfigurationInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteEncoderConfigurationOutput, error)
	ListStages(ctx context.Context, params *ivsrealtime.ListStagesInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.ListStagesOutput, error)
	DeleteStage(ctx context.Context, params *ivsrealtime.DeleteStageInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteStageOutput, error)
	ListIngestConfigurations(ctx context.Context, params *ivsrealtime.ListIngestConfigurationsInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.ListIngestConfigurationsOutput, error)
	DeleteIngestConfiguration(ctx context.Context, params *ivsrealtime.DeleteIngestConfigurationInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteIngestConfigurationOutput, error)

	ListStorageConfigurations(ctx context.Context, params *ivsrealtime.ListStorageConfigurationsInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.ListStorageConfigurationsOutput, error)
	DeleteStorageConfiguration(ctx context.Context, params *ivsrealtime.DeleteStorageConfigurationInput,
		optFns ...func(*ivsrealtime.Options)) (*ivsrealtime.DeleteStorageConfigurationOutput, error)
}

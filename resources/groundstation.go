package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/groundstation"
)

// GroundStationClient is the interface for the Ground Station SDK client methods.
type GroundStationClient interface {
	ListConfigs(ctx context.Context, params *groundstation.ListConfigsInput,
		optFns ...func(*groundstation.Options)) (*groundstation.ListConfigsOutput, error)
	DeleteConfig(ctx context.Context, params *groundstation.DeleteConfigInput,
		optFns ...func(*groundstation.Options)) (*groundstation.DeleteConfigOutput, error)
}

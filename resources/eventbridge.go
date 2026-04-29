package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

// EventBridgeClient is an interface for the AWS EventBridge SDK v2 client methods used by all EventBridge resources.
// It enables mock testing of List and Remove operations.
// Note: This is separate from the existing cloudwatchevents resources which use SDK v1.
type EventBridgeClient interface {
	ListArchives(ctx context.Context, params *eventbridge.ListArchivesInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.ListArchivesOutput, error)
	DeleteArchive(ctx context.Context, params *eventbridge.DeleteArchiveInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.DeleteArchiveOutput, error)
	ListConnections(ctx context.Context, params *eventbridge.ListConnectionsInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.ListConnectionsOutput, error)
	DeleteConnection(ctx context.Context, params *eventbridge.DeleteConnectionInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.DeleteConnectionOutput, error)
	ListEventBuses(ctx context.Context, params *eventbridge.ListEventBusesInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.ListEventBusesOutput, error)
	DeleteEventBus(ctx context.Context, params *eventbridge.DeleteEventBusInput,
		optFns ...func(*eventbridge.Options)) (*eventbridge.DeleteEventBusOutput, error)
}

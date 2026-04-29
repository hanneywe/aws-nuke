package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
)

// NeptuneGraphClient is the interface for the Neptune Graph SDK client methods.
type NeptuneGraphClient interface {
	ListGraphSnapshots(ctx context.Context, params *neptunegraph.ListGraphSnapshotsInput,
		optFns ...func(*neptunegraph.Options)) (*neptunegraph.ListGraphSnapshotsOutput, error)
	DeleteGraphSnapshot(ctx context.Context, params *neptunegraph.DeleteGraphSnapshotInput,
		optFns ...func(*neptunegraph.Options)) (*neptunegraph.DeleteGraphSnapshotOutput, error)
}

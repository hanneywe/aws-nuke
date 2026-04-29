package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockNeptuneGraphClient struct {
	mock.Mock
}

func (m *mockNeptuneGraphClient) ListGraphSnapshots(ctx context.Context,
	params *neptunegraph.ListGraphSnapshotsInput,
	_ ...func(*neptunegraph.Options)) (*neptunegraph.ListGraphSnapshotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptunegraph.ListGraphSnapshotsOutput), args.Error(1)
}

func (m *mockNeptuneGraphClient) DeleteGraphSnapshot(ctx context.Context,
	params *neptunegraph.DeleteGraphSnapshotInput,
	_ ...func(*neptunegraph.Options)) (*neptunegraph.DeleteGraphSnapshotOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptunegraph.DeleteGraphSnapshotOutput), args.Error(1)
}

var testNeptuneGraphListerOpts = &nuke.ListerOpts{}

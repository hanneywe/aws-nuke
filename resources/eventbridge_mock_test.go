package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

type mockEventBridgeClient struct {
	mock.Mock
}

func (m *mockEventBridgeClient) ListArchives(ctx context.Context, params *eventbridge.ListArchivesInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.ListArchivesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.ListArchivesOutput), args.Error(1)
}

func (m *mockEventBridgeClient) DeleteArchive(ctx context.Context, params *eventbridge.DeleteArchiveInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.DeleteArchiveOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.DeleteArchiveOutput), args.Error(1)
}

func (m *mockEventBridgeClient) ListConnections(ctx context.Context, params *eventbridge.ListConnectionsInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.ListConnectionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.ListConnectionsOutput), args.Error(1)
}

func (m *mockEventBridgeClient) DeleteConnection(ctx context.Context, params *eventbridge.DeleteConnectionInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.DeleteConnectionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.DeleteConnectionOutput), args.Error(1)
}

func (m *mockEventBridgeClient) ListEventBuses(ctx context.Context, params *eventbridge.ListEventBusesInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.ListEventBusesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.ListEventBusesOutput), args.Error(1)
}

func (m *mockEventBridgeClient) DeleteEventBus(ctx context.Context, params *eventbridge.DeleteEventBusInput,
	_ ...func(*eventbridge.Options)) (*eventbridge.DeleteEventBusOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eventbridge.DeleteEventBusOutput), args.Error(1)
}

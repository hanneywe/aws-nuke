package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/datasync"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDataSyncClient struct {
	mock.Mock
}

func (m *mockDataSyncClient) ListLocations(ctx context.Context,
	params *datasync.ListLocationsInput,
	_ ...func(*datasync.Options)) (*datasync.ListLocationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*datasync.ListLocationsOutput), args.Error(1)
}

func (m *mockDataSyncClient) DeleteLocation(ctx context.Context,
	params *datasync.DeleteLocationInput,
	_ ...func(*datasync.Options)) (*datasync.DeleteLocationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*datasync.DeleteLocationOutput), args.Error(1)
}

var testDataSyncListerOpts = &nuke.ListerOpts{}

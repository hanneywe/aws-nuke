package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/groundstation"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockGroundStationClient struct {
	mock.Mock
}

func (m *mockGroundStationClient) ListConfigs(ctx context.Context,
	params *groundstation.ListConfigsInput,
	_ ...func(*groundstation.Options)) (*groundstation.ListConfigsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*groundstation.ListConfigsOutput), args.Error(1)
}

func (m *mockGroundStationClient) DeleteConfig(ctx context.Context,
	params *groundstation.DeleteConfigInput,
	_ ...func(*groundstation.Options)) (*groundstation.DeleteConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*groundstation.DeleteConfigOutput), args.Error(1)
}

var testGroundStationListerOpts = &nuke.ListerOpts{}

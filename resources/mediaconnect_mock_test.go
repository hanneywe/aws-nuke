package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testMediaConnectListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockMediaConnectClient struct {
	mock.Mock
}

func (m *mockMediaConnectClient) ListGateways(ctx context.Context, params *mediaconnect.ListGatewaysInput,
	_ ...func(*mediaconnect.Options)) (*mediaconnect.ListGatewaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediaconnect.ListGatewaysOutput), args.Error(1)
}

func (m *mockMediaConnectClient) DeleteGateway(ctx context.Context, params *mediaconnect.DeleteGatewayInput,
	_ ...func(*mediaconnect.Options)) (*mediaconnect.DeleteGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediaconnect.DeleteGatewayOutput), args.Error(1)
}

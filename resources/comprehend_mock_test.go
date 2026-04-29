package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/comprehend"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockComprehendClient struct {
	mock.Mock
}

func (m *mockComprehendClient) ListFlywheels(
	ctx context.Context, params *comprehend.ListFlywheelsInput,
	_ ...func(*comprehend.Options),
) (*comprehend.ListFlywheelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*comprehend.ListFlywheelsOutput), args.Error(1)
}

func (m *mockComprehendClient) DeleteFlywheel(
	ctx context.Context, params *comprehend.DeleteFlywheelInput,
	_ ...func(*comprehend.Options),
) (*comprehend.DeleteFlywheelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*comprehend.DeleteFlywheelOutput), args.Error(1)
}

var testComprehendListerOpts = &nuke.ListerOpts{}

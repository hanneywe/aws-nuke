package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSFNv2Client struct {
	mock.Mock
}

func (m *mockSFNv2Client) ListActivities(
	ctx context.Context, params *sfn.ListActivitiesInput,
	_ ...func(*sfn.Options),
) (*sfn.ListActivitiesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sfn.ListActivitiesOutput), args.Error(1)
}

func (m *mockSFNv2Client) DeleteActivity(
	ctx context.Context, params *sfn.DeleteActivityInput,
	_ ...func(*sfn.Options),
) (*sfn.DeleteActivityOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sfn.DeleteActivityOutput), args.Error(1)
}

var testSFNv2ListerOpts = &nuke.ListerOpts{}

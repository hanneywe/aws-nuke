package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCodeBuildClient struct {
	mock.Mock
}

func (m *mockCodeBuildClient) ListFleets(
	ctx context.Context, params *codebuild.ListFleetsInput,
	_ ...func(*codebuild.Options),
) (*codebuild.ListFleetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codebuild.ListFleetsOutput), args.Error(1)
}

func (m *mockCodeBuildClient) BatchGetFleets(
	ctx context.Context, params *codebuild.BatchGetFleetsInput,
	_ ...func(*codebuild.Options),
) (*codebuild.BatchGetFleetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codebuild.BatchGetFleetsOutput), args.Error(1)
}

func (m *mockCodeBuildClient) DeleteFleet(
	ctx context.Context, params *codebuild.DeleteFleetInput,
	_ ...func(*codebuild.Options),
) (*codebuild.DeleteFleetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codebuild.DeleteFleetOutput), args.Error(1)
}

var testCodeBuildListerOpts = &nuke.ListerOpts{}

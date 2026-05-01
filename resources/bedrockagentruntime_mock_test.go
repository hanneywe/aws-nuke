package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockBedrockagentruntimeClient struct {
	mock.Mock
}

func (m *mockBedrockagentruntimeClient) ListSessions(
	ctx context.Context, params *bedrockagentruntime.ListSessionsInput,
	_ ...func(*bedrockagentruntime.Options),
) (*bedrockagentruntime.ListSessionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bedrockagentruntime.ListSessionsOutput), args.Error(1)
}

func (m *mockBedrockagentruntimeClient) EndSession(
	ctx context.Context, params *bedrockagentruntime.EndSessionInput,
	_ ...func(*bedrockagentruntime.Options),
) (*bedrockagentruntime.EndSessionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bedrockagentruntime.EndSessionOutput), args.Error(1)
}

func (m *mockBedrockagentruntimeClient) DeleteSession(
	ctx context.Context, params *bedrockagentruntime.DeleteSessionInput,
	_ ...func(*bedrockagentruntime.Options),
) (*bedrockagentruntime.DeleteSessionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bedrockagentruntime.DeleteSessionOutput), args.Error(1)
}

var testBedrockagentruntimeListerOpts = &nuke.ListerOpts{}

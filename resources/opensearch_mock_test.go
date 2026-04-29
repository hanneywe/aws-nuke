package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/opensearch"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockOpenSearchClient struct {
	mock.Mock
}

func (m *mockOpenSearchClient) ListApplications(ctx context.Context,
	params *opensearch.ListApplicationsInput,
	_ ...func(*opensearch.Options)) (*opensearch.ListApplicationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*opensearch.ListApplicationsOutput), args.Error(1)
}

func (m *mockOpenSearchClient) DeleteApplication(ctx context.Context,
	params *opensearch.DeleteApplicationInput,
	_ ...func(*opensearch.Options)) (*opensearch.DeleteApplicationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*opensearch.DeleteApplicationOutput), args.Error(1)
}

var testOpenSearchListerOpts = &nuke.ListerOpts{}

package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/dataexchange"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDataexchangeClient struct {
	mock.Mock
}

func (m *mockDataexchangeClient) ListDataSets(
	ctx context.Context, params *dataexchange.ListDataSetsInput,
	_ ...func(*dataexchange.Options),
) (*dataexchange.ListDataSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dataexchange.ListDataSetsOutput), args.Error(1)
}

func (m *mockDataexchangeClient) DeleteDataSet(
	ctx context.Context, params *dataexchange.DeleteDataSetInput,
	_ ...func(*dataexchange.Options),
) (*dataexchange.DeleteDataSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dataexchange.DeleteDataSetOutput), args.Error(1)
}

func (m *mockDataexchangeClient) ListDataSetRevisions(
	ctx context.Context, params *dataexchange.ListDataSetRevisionsInput,
	_ ...func(*dataexchange.Options),
) (*dataexchange.ListDataSetRevisionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dataexchange.ListDataSetRevisionsOutput), args.Error(1)
}

func (m *mockDataexchangeClient) DeleteRevision(
	ctx context.Context, params *dataexchange.DeleteRevisionInput,
	_ ...func(*dataexchange.Options),
) (*dataexchange.DeleteRevisionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dataexchange.DeleteRevisionOutput), args.Error(1)
}

var testDataexchangeListerOpts = &nuke.ListerOpts{}

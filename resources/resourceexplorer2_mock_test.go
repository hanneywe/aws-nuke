package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockResourceExplorer2Client struct {
	mock.Mock
}

func (m *mockResourceExplorer2Client) ListIndexes(
	ctx context.Context, params *resourceexplorer2.ListIndexesInput,
	_ ...func(*resourceexplorer2.Options),
) (*resourceexplorer2.ListIndexesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*resourceexplorer2.ListIndexesOutput), args.Error(1)
}

func (m *mockResourceExplorer2Client) ListTagsForResource(
	ctx context.Context, params *resourceexplorer2.ListTagsForResourceInput,
	_ ...func(*resourceexplorer2.Options),
) (*resourceexplorer2.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*resourceexplorer2.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockResourceExplorer2Client) DeleteIndex(
	ctx context.Context, params *resourceexplorer2.DeleteIndexInput,
	_ ...func(*resourceexplorer2.Options),
) (*resourceexplorer2.DeleteIndexOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*resourceexplorer2.DeleteIndexOutput), args.Error(1)
}

var testResourceExplorer2ListerOpts = &nuke.ListerOpts{}

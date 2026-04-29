package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagevod"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMediaPackageVODClient struct {
	mock.Mock
}

func (m *mockMediaPackageVODClient) ListPackagingGroups(
	ctx context.Context, params *mediapackagevod.ListPackagingGroupsInput,
	_ ...func(*mediapackagevod.Options),
) (*mediapackagevod.ListPackagingGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagevod.ListPackagingGroupsOutput), args.Error(1)
}

func (m *mockMediaPackageVODClient) DeletePackagingGroup(
	ctx context.Context, params *mediapackagevod.DeletePackagingGroupInput,
	_ ...func(*mediapackagevod.Options),
) (*mediapackagevod.DeletePackagingGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mediapackagevod.DeletePackagingGroupOutput), args.Error(1)
}

var testMediaPackageVODListerOpts = &nuke.ListerOpts{}

package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/panorama"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockPanoramaClient struct {
	mock.Mock
}

func (m *mockPanoramaClient) ListPackages(
	ctx context.Context, params *panorama.ListPackagesInput,
	_ ...func(*panorama.Options),
) (*panorama.ListPackagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*panorama.ListPackagesOutput), args.Error(1)
}

func (m *mockPanoramaClient) DeletePackage(
	ctx context.Context, params *panorama.DeletePackageInput,
	_ ...func(*panorama.Options),
) (*panorama.DeletePackageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*panorama.DeletePackageOutput), args.Error(1)
}

var testPanoramaListerOpts = &nuke.ListerOpts{}

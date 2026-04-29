package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codeartifact"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCodeartifactClient struct {
	mock.Mock
}

func (m *mockCodeartifactClient) ListDomains(
	ctx context.Context, params *codeartifact.ListDomainsInput,
	_ ...func(*codeartifact.Options),
) (*codeartifact.ListDomainsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeartifact.ListDomainsOutput), args.Error(1)
}

func (m *mockCodeartifactClient) ListRepositoriesInDomain(
	ctx context.Context, params *codeartifact.ListRepositoriesInDomainInput,
	_ ...func(*codeartifact.Options),
) (*codeartifact.ListRepositoriesInDomainOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeartifact.ListRepositoriesInDomainOutput), args.Error(1)
}

func (m *mockCodeartifactClient) ListPackages(
	ctx context.Context, params *codeartifact.ListPackagesInput,
	_ ...func(*codeartifact.Options),
) (*codeartifact.ListPackagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeartifact.ListPackagesOutput), args.Error(1)
}

func (m *mockCodeartifactClient) DeletePackageVersions(
	ctx context.Context, params *codeartifact.DeletePackageVersionsInput,
	_ ...func(*codeartifact.Options),
) (*codeartifact.DeletePackageVersionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeartifact.DeletePackageVersionsOutput), args.Error(1)
}

func (m *mockCodeartifactClient) ListPackageVersions(
	ctx context.Context, params *codeartifact.ListPackageVersionsInput,
	_ ...func(*codeartifact.Options),
) (*codeartifact.ListPackageVersionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeartifact.ListPackageVersionsOutput), args.Error(1)
}

var testCodeartifactListerOpts = &nuke.ListerOpts{}

package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockECRv2Client struct {
	mock.Mock
}

func (m *mockECRv2Client) DescribeRepositoryCreationTemplates(
	ctx context.Context, params *ecr.DescribeRepositoryCreationTemplatesInput,
	_ ...func(*ecr.Options),
) (*ecr.DescribeRepositoryCreationTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DescribeRepositoryCreationTemplatesOutput), args.Error(1)
}

func (m *mockECRv2Client) DeleteRepositoryCreationTemplate(
	ctx context.Context, params *ecr.DeleteRepositoryCreationTemplateInput,
	_ ...func(*ecr.Options),
) (*ecr.DeleteRepositoryCreationTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DeleteRepositoryCreationTemplateOutput), args.Error(1)
}

func (m *mockECRv2Client) DescribePullThroughCacheRules(
	ctx context.Context, params *ecr.DescribePullThroughCacheRulesInput,
	_ ...func(*ecr.Options),
) (*ecr.DescribePullThroughCacheRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DescribePullThroughCacheRulesOutput), args.Error(1)
}

func (m *mockECRv2Client) DeletePullThroughCacheRule(
	ctx context.Context, params *ecr.DeletePullThroughCacheRuleInput,
	_ ...func(*ecr.Options),
) (*ecr.DeletePullThroughCacheRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DeletePullThroughCacheRuleOutput), args.Error(1)
}

func (m *mockECRv2Client) DescribeRepositories(
	ctx context.Context, params *ecr.DescribeRepositoriesInput,
	_ ...func(*ecr.Options),
) (*ecr.DescribeRepositoriesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DescribeRepositoriesOutput), args.Error(1)
}

func (m *mockECRv2Client) PutImageScanningConfiguration(
	ctx context.Context, params *ecr.PutImageScanningConfigurationInput,
	_ ...func(*ecr.Options),
) (*ecr.PutImageScanningConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.PutImageScanningConfigurationOutput), args.Error(1)
}

func (m *mockECRv2Client) PutImageTagMutability(
	ctx context.Context, params *ecr.PutImageTagMutabilityInput,
	_ ...func(*ecr.Options),
) (*ecr.PutImageTagMutabilityOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.PutImageTagMutabilityOutput), args.Error(1)
}

func (m *mockECRv2Client) GetRegistryScanningConfiguration(
	ctx context.Context, params *ecr.GetRegistryScanningConfigurationInput,
	_ ...func(*ecr.Options),
) (*ecr.GetRegistryScanningConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.GetRegistryScanningConfigurationOutput), args.Error(1)
}

func (m *mockECRv2Client) PutRegistryScanningConfiguration(
	ctx context.Context, params *ecr.PutRegistryScanningConfigurationInput,
	_ ...func(*ecr.Options),
) (*ecr.PutRegistryScanningConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.PutRegistryScanningConfigurationOutput), args.Error(1)
}

func (m *mockECRv2Client) ListPullTimeUpdateExclusions(
	ctx context.Context, params *ecr.ListPullTimeUpdateExclusionsInput,
	_ ...func(*ecr.Options),
) (*ecr.ListPullTimeUpdateExclusionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.ListPullTimeUpdateExclusionsOutput), args.Error(1)
}

func (m *mockECRv2Client) DeregisterPullTimeUpdateExclusion(
	ctx context.Context, params *ecr.DeregisterPullTimeUpdateExclusionInput,
	_ ...func(*ecr.Options),
) (*ecr.DeregisterPullTimeUpdateExclusionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ecr.DeregisterPullTimeUpdateExclusionOutput), args.Error(1)
}

var testECRv2ListerOpts = &nuke.ListerOpts{}

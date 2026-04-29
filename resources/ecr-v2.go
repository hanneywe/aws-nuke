package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// ECRv2Client is the interface for the ECR SDK v2 client methods.
type ECRv2Client interface {
	DescribeRepositoryCreationTemplates(ctx context.Context, params *ecr.DescribeRepositoryCreationTemplatesInput,
		optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoryCreationTemplatesOutput, error)
	DeleteRepositoryCreationTemplate(ctx context.Context, params *ecr.DeleteRepositoryCreationTemplateInput,
		optFns ...func(*ecr.Options)) (*ecr.DeleteRepositoryCreationTemplateOutput, error)
	DescribePullThroughCacheRules(ctx context.Context, params *ecr.DescribePullThroughCacheRulesInput,
		optFns ...func(*ecr.Options)) (*ecr.DescribePullThroughCacheRulesOutput, error)
	DeletePullThroughCacheRule(ctx context.Context, params *ecr.DeletePullThroughCacheRuleInput,
		optFns ...func(*ecr.Options)) (*ecr.DeletePullThroughCacheRuleOutput, error)
	DescribeRepositories(ctx context.Context, params *ecr.DescribeRepositoriesInput,
		optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
	PutImageScanningConfiguration(ctx context.Context, params *ecr.PutImageScanningConfigurationInput,
		optFns ...func(*ecr.Options)) (*ecr.PutImageScanningConfigurationOutput, error)
	PutImageTagMutability(ctx context.Context, params *ecr.PutImageTagMutabilityInput,
		optFns ...func(*ecr.Options)) (*ecr.PutImageTagMutabilityOutput, error)
	GetRegistryScanningConfiguration(ctx context.Context, params *ecr.GetRegistryScanningConfigurationInput,
		optFns ...func(*ecr.Options)) (*ecr.GetRegistryScanningConfigurationOutput, error)
	PutRegistryScanningConfiguration(ctx context.Context, params *ecr.PutRegistryScanningConfigurationInput,
		optFns ...func(*ecr.Options)) (*ecr.PutRegistryScanningConfigurationOutput, error)
	ListPullTimeUpdateExclusions(ctx context.Context, params *ecr.ListPullTimeUpdateExclusionsInput,
		optFns ...func(*ecr.Options)) (*ecr.ListPullTimeUpdateExclusionsOutput, error)
	DeregisterPullTimeUpdateExclusion(ctx context.Context, params *ecr.DeregisterPullTimeUpdateExclusionInput,
		optFns ...func(*ecr.Options)) (*ecr.DeregisterPullTimeUpdateExclusionOutput, error)
}

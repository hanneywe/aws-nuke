package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

// SageMakerV2Client is the interface for the SageMaker SDK v2 client methods
// used by new SageMaker sub-resources. Existing SageMaker resources use SDK v1.
type SageMakerV2Client interface {
	ListArtifacts(ctx context.Context, params *sagemaker.ListArtifactsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListArtifactsOutput, error)
	DeleteArtifact(ctx context.Context, params *sagemaker.DeleteArtifactInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteArtifactOutput, error)

	ListCodeRepositories(ctx context.Context, params *sagemaker.ListCodeRepositoriesInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListCodeRepositoriesOutput, error)
	DeleteCodeRepository(ctx context.Context, params *sagemaker.DeleteCodeRepositoryInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteCodeRepositoryOutput, error)

	ListExperiments(ctx context.Context, params *sagemaker.ListExperimentsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListExperimentsOutput, error)
	DeleteExperiment(ctx context.Context, params *sagemaker.DeleteExperimentInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteExperimentOutput, error)

	ListHubs(ctx context.Context, params *sagemaker.ListHubsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListHubsOutput, error)
	DeleteHub(ctx context.Context, params *sagemaker.DeleteHubInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteHubOutput, error)

	ListHumanTaskUis(ctx context.Context, params *sagemaker.ListHumanTaskUisInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListHumanTaskUisOutput, error)
	DeleteHumanTaskUi(ctx context.Context, params *sagemaker.DeleteHumanTaskUiInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteHumanTaskUiOutput, error)

	ListModelCards(ctx context.Context, params *sagemaker.ListModelCardsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListModelCardsOutput, error)
	DeleteModelCard(ctx context.Context, params *sagemaker.DeleteModelCardInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteModelCardOutput, error)

	ListModelPackageGroups(ctx context.Context, params *sagemaker.ListModelPackageGroupsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListModelPackageGroupsOutput, error)
	DeleteModelPackageGroup(ctx context.Context, params *sagemaker.DeleteModelPackageGroupInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteModelPackageGroupOutput, error)

	ListStudioLifecycleConfigs(ctx context.Context, params *sagemaker.ListStudioLifecycleConfigsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListStudioLifecycleConfigsOutput, error)
	DeleteStudioLifecycleConfig(ctx context.Context, params *sagemaker.DeleteStudioLifecycleConfigInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteStudioLifecycleConfigOutput, error)

	ListTrialComponents(ctx context.Context, params *sagemaker.ListTrialComponentsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListTrialComponentsOutput, error)
	DeleteTrialComponent(ctx context.Context, params *sagemaker.DeleteTrialComponentInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteTrialComponentOutput, error)

	ListActions(ctx context.Context, params *sagemaker.ListActionsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListActionsOutput, error)
	DeleteAction(ctx context.Context, params *sagemaker.DeleteActionInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteActionOutput, error)

	ListContexts(ctx context.Context, params *sagemaker.ListContextsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListContextsOutput, error)
	DeleteContext(ctx context.Context, params *sagemaker.DeleteContextInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteContextOutput, error)

	ListAssociations(ctx context.Context, params *sagemaker.ListAssociationsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListAssociationsOutput, error)
	DeleteAssociation(ctx context.Context, params *sagemaker.DeleteAssociationInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteAssociationOutput, error)

	ListImages(ctx context.Context, params *sagemaker.ListImagesInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListImagesOutput, error)
	DeleteImage(ctx context.Context, params *sagemaker.DeleteImageInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteImageOutput, error)

	ListMlflowTrackingServers(ctx context.Context, params *sagemaker.ListMlflowTrackingServersInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListMlflowTrackingServersOutput, error)
	DeleteMlflowTrackingServer(ctx context.Context, params *sagemaker.DeleteMlflowTrackingServerInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteMlflowTrackingServerOutput, error)

	ListMlflowApps(ctx context.Context, params *sagemaker.ListMlflowAppsInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.ListMlflowAppsOutput, error)
	DeleteMlflowApp(ctx context.Context, params *sagemaker.DeleteMlflowAppInput,
		optFns ...func(*sagemaker.Options)) (*sagemaker.DeleteMlflowAppOutput, error)
}

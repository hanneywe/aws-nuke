package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

type mockSageMakerV2Client struct {
	mock.Mock
}

func (m *mockSageMakerV2Client) ListArtifacts(ctx context.Context, params *sagemaker.ListArtifactsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListArtifactsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListArtifactsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteArtifact(ctx context.Context, params *sagemaker.DeleteArtifactInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteArtifactOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteArtifactOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListCodeRepositories(ctx context.Context, params *sagemaker.ListCodeRepositoriesInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListCodeRepositoriesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListCodeRepositoriesOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteCodeRepository(ctx context.Context, params *sagemaker.DeleteCodeRepositoryInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteCodeRepositoryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteCodeRepositoryOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListExperiments(ctx context.Context, params *sagemaker.ListExperimentsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListExperimentsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListExperimentsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteExperiment(ctx context.Context, params *sagemaker.DeleteExperimentInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteExperimentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteExperimentOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListHubs(ctx context.Context, params *sagemaker.ListHubsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListHubsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListHubsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteHub(ctx context.Context, params *sagemaker.DeleteHubInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteHubOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteHubOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListHumanTaskUis(ctx context.Context, params *sagemaker.ListHumanTaskUisInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListHumanTaskUisOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListHumanTaskUisOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteHumanTaskUi(
	ctx context.Context, params *sagemaker.DeleteHumanTaskUiInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteHumanTaskUiOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteHumanTaskUiOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListModelCards(ctx context.Context, params *sagemaker.ListModelCardsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListModelCardsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListModelCardsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteModelCard(ctx context.Context, params *sagemaker.DeleteModelCardInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteModelCardOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteModelCardOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListModelPackageGroups(ctx context.Context, params *sagemaker.ListModelPackageGroupsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListModelPackageGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListModelPackageGroupsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteModelPackageGroup(ctx context.Context, params *sagemaker.DeleteModelPackageGroupInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteModelPackageGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteModelPackageGroupOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListStudioLifecycleConfigs(ctx context.Context, params *sagemaker.ListStudioLifecycleConfigsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListStudioLifecycleConfigsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListStudioLifecycleConfigsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteStudioLifecycleConfig(ctx context.Context, params *sagemaker.DeleteStudioLifecycleConfigInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteStudioLifecycleConfigOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteStudioLifecycleConfigOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListTrialComponents(ctx context.Context, params *sagemaker.ListTrialComponentsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListTrialComponentsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListTrialComponentsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteTrialComponent(ctx context.Context, params *sagemaker.DeleteTrialComponentInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteTrialComponentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteTrialComponentOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListActions(ctx context.Context, params *sagemaker.ListActionsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListActionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListActionsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteAction(ctx context.Context, params *sagemaker.DeleteActionInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteActionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteActionOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListContexts(ctx context.Context, params *sagemaker.ListContextsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListContextsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListContextsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteContext(ctx context.Context, params *sagemaker.DeleteContextInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteContextOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteContextOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListAssociations(ctx context.Context, params *sagemaker.ListAssociationsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListAssociationsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteAssociation(ctx context.Context, params *sagemaker.DeleteAssociationInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteAssociationOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListImages(ctx context.Context, params *sagemaker.ListImagesInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListImagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListImagesOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteImage(ctx context.Context, params *sagemaker.DeleteImageInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteImageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteImageOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListMlflowTrackingServers(ctx context.Context, params *sagemaker.ListMlflowTrackingServersInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListMlflowTrackingServersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListMlflowTrackingServersOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteMlflowTrackingServer(ctx context.Context, params *sagemaker.DeleteMlflowTrackingServerInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteMlflowTrackingServerOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteMlflowTrackingServerOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) ListMlflowApps(ctx context.Context, params *sagemaker.ListMlflowAppsInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.ListMlflowAppsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.ListMlflowAppsOutput), args.Error(1)
}

func (m *mockSageMakerV2Client) DeleteMlflowApp(ctx context.Context, params *sagemaker.DeleteMlflowAppInput,
	_ ...func(*sagemaker.Options)) (*sagemaker.DeleteMlflowAppOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sagemaker.DeleteMlflowAppOutput), args.Error(1)
}

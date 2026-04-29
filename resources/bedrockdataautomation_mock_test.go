package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockBedrockDataAutomationClient struct {
	mock.Mock
}

func (m *mockBedrockDataAutomationClient) ListDataAutomationProjects(ctx context.Context,
	params *bedrockdataautomation.ListDataAutomationProjectsInput,
	_ ...func(*bedrockdataautomation.Options)) (*bedrockdataautomation.ListDataAutomationProjectsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bedrockdataautomation.ListDataAutomationProjectsOutput), args.Error(1)
}

func (m *mockBedrockDataAutomationClient) DeleteDataAutomationProject(ctx context.Context,
	params *bedrockdataautomation.DeleteDataAutomationProjectInput,
	_ ...func(*bedrockdataautomation.Options)) (*bedrockdataautomation.DeleteDataAutomationProjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*bedrockdataautomation.DeleteDataAutomationProjectOutput), args.Error(1)
}

var testBedrockDataAutomationListerOpts = &nuke.ListerOpts{}

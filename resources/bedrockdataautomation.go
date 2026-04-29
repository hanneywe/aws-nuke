package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation"
)

// BedrockDataAutomationClient is the interface for the Bedrock Data Automation SDK client methods.
type BedrockDataAutomationClient interface {
	ListDataAutomationProjects(ctx context.Context, params *bedrockdataautomation.ListDataAutomationProjectsInput,
		optFns ...func(*bedrockdataautomation.Options)) (*bedrockdataautomation.ListDataAutomationProjectsOutput, error)
	DeleteDataAutomationProject(ctx context.Context, params *bedrockdataautomation.DeleteDataAutomationProjectInput,
		optFns ...func(*bedrockdataautomation.Options)) (*bedrockdataautomation.DeleteDataAutomationProjectOutput, error)
}

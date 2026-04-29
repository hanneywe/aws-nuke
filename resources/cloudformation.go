package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// CloudFormationClient is the interface for the CloudFormation SDK v2 client methods.
type CloudFormationClient interface {
	ListGeneratedTemplates(ctx context.Context, params *cloudformation.ListGeneratedTemplatesInput,
		optFns ...func(*cloudformation.Options)) (*cloudformation.ListGeneratedTemplatesOutput, error)
	DeleteGeneratedTemplate(ctx context.Context, params *cloudformation.DeleteGeneratedTemplateInput,
		optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteGeneratedTemplateOutput, error)
}

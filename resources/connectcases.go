package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"
)

// ConnectCasesClient is the interface for the ConnectCases SDK client methods.
type ConnectCasesClient interface {
	ListDomains(ctx context.Context, params *connectcases.ListDomainsInput,
		optFns ...func(*connectcases.Options)) (*connectcases.ListDomainsOutput, error)
	DeleteDomain(ctx context.Context, params *connectcases.DeleteDomainInput,
		optFns ...func(*connectcases.Options)) (*connectcases.DeleteDomainOutput, error)
	ListFields(ctx context.Context, params *connectcases.ListFieldsInput,
		optFns ...func(*connectcases.Options)) (*connectcases.ListFieldsOutput, error)
	DeleteField(ctx context.Context, params *connectcases.DeleteFieldInput,
		optFns ...func(*connectcases.Options)) (*connectcases.DeleteFieldOutput, error)
	ListTemplates(ctx context.Context, params *connectcases.ListTemplatesInput,
		optFns ...func(*connectcases.Options)) (*connectcases.ListTemplatesOutput, error)
	DeleteTemplate(ctx context.Context, params *connectcases.DeleteTemplateInput,
		optFns ...func(*connectcases.Options)) (*connectcases.DeleteTemplateOutput, error)
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

// SESv2Client is the interface for the SESv2 SDK client methods.
type SESv2Client interface {
	ListConfigurationSets(ctx context.Context, params *sesv2.ListConfigurationSetsInput,
		optFns ...func(*sesv2.Options)) (*sesv2.ListConfigurationSetsOutput, error)
	DeleteConfigurationSet(ctx context.Context, params *sesv2.DeleteConfigurationSetInput,
		optFns ...func(*sesv2.Options)) (*sesv2.DeleteConfigurationSetOutput, error)
	ListDedicatedIpPools(ctx context.Context, params *sesv2.ListDedicatedIpPoolsInput,
		optFns ...func(*sesv2.Options)) (*sesv2.ListDedicatedIpPoolsOutput, error)
	DeleteDedicatedIpPool(ctx context.Context, params *sesv2.DeleteDedicatedIpPoolInput,
		optFns ...func(*sesv2.Options)) (*sesv2.DeleteDedicatedIpPoolOutput, error)
	ListEmailIdentities(ctx context.Context, params *sesv2.ListEmailIdentitiesInput,
		optFns ...func(*sesv2.Options)) (*sesv2.ListEmailIdentitiesOutput, error)
	DeleteEmailIdentity(ctx context.Context, params *sesv2.DeleteEmailIdentityInput,
		optFns ...func(*sesv2.Options)) (*sesv2.DeleteEmailIdentityOutput, error)
	ListEmailTemplates(ctx context.Context, params *sesv2.ListEmailTemplatesInput,
		optFns ...func(*sesv2.Options)) (*sesv2.ListEmailTemplatesOutput, error)
	DeleteEmailTemplate(ctx context.Context, params *sesv2.DeleteEmailTemplateInput,
		optFns ...func(*sesv2.Options)) (*sesv2.DeleteEmailTemplateOutput, error)
}

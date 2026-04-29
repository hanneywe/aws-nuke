package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dlm"
)

// DlmClient is the interface for the dlm SDK client methods.
type DlmClient interface {
	GetLifecyclePolicies(ctx context.Context, params *dlm.GetLifecyclePoliciesInput,
		optFns ...func(*dlm.Options)) (*dlm.GetLifecyclePoliciesOutput, error)
	DeleteLifecyclePolicy(ctx context.Context, params *dlm.DeleteLifecyclePolicyInput,
		optFns ...func(*dlm.Options)) (*dlm.DeleteLifecyclePolicyOutput, error)
}

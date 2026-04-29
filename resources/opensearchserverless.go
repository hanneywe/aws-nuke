package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
)

// OpenSearchServerlessClient is the interface for the OpenSearch Serverless SDK v2 client methods
// used by all OpenSearch Serverless resources. It enables mock testing of List and Remove operations.
type OpenSearchServerlessClient interface {
	ListLifecyclePolicies(ctx context.Context, params *opensearchserverless.ListLifecyclePoliciesInput,
		optFns ...func(*opensearchserverless.Options)) (*opensearchserverless.ListLifecyclePoliciesOutput, error)
	DeleteLifecyclePolicy(ctx context.Context, params *opensearchserverless.DeleteLifecyclePolicyInput,
		optFns ...func(*opensearchserverless.Options)) (*opensearchserverless.DeleteLifecyclePolicyOutput, error)
}

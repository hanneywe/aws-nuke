package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockOpenSearchServerlessClient struct {
	mock.Mock
}

func (m *mockOpenSearchServerlessClient) ListLifecyclePolicies(ctx context.Context,
	params *opensearchserverless.ListLifecyclePoliciesInput,
	_ ...func(*opensearchserverless.Options)) (*opensearchserverless.ListLifecyclePoliciesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*opensearchserverless.ListLifecyclePoliciesOutput), args.Error(1)
}

func (m *mockOpenSearchServerlessClient) DeleteLifecyclePolicy(ctx context.Context,
	params *opensearchserverless.DeleteLifecyclePolicyInput,
	_ ...func(*opensearchserverless.Options)) (*opensearchserverless.DeleteLifecyclePolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*opensearchserverless.DeleteLifecyclePolicyOutput), args.Error(1)
}

var testOpenSearchServerlessListerOpts = &nuke.ListerOpts{}

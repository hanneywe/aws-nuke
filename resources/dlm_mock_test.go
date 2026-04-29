package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/dlm"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDlmClient struct {
	mock.Mock
}

func (m *mockDlmClient) GetLifecyclePolicies(
	ctx context.Context, params *dlm.GetLifecyclePoliciesInput,
	_ ...func(*dlm.Options),
) (*dlm.GetLifecyclePoliciesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dlm.GetLifecyclePoliciesOutput), args.Error(1)
}

func (m *mockDlmClient) DeleteLifecyclePolicy(
	ctx context.Context, params *dlm.DeleteLifecyclePolicyInput,
	_ ...func(*dlm.Options),
) (*dlm.DeleteLifecyclePolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dlm.DeleteLifecyclePolicyOutput), args.Error(1)
}

var testDlmListerOpts = &nuke.ListerOpts{}

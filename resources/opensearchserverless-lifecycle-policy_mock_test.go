package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	osstypes "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
)

func Test_Mock_OpenSearchServerlessLifecyclePolicy_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchServerlessClient)

	mockClient.On("ListLifecyclePolicies", mock.Anything, mock.Anything).
		Return(&opensearchserverless.ListLifecyclePoliciesOutput{
			LifecyclePolicySummaries: []osstypes.LifecyclePolicySummary{
				{
					Name: ptr.String("my-retention-policy"),
					Type: osstypes.LifecyclePolicyTypeRetention,
				},
			},
		}, nil)

	lister := &OpenSearchServerlessLifecyclePolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOpenSearchServerlessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	policy := resources[0].(*OpenSearchServerlessLifecyclePolicy)
	a.Equal("my-retention-policy", *policy.Name)
	a.Equal(osstypes.LifecyclePolicyTypeRetention, policy.Type)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchServerlessLifecyclePolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchServerlessClient)

	mockClient.On("ListLifecyclePolicies", mock.Anything, mock.Anything).
		Return(&opensearchserverless.ListLifecyclePoliciesOutput{
			LifecyclePolicySummaries: []osstypes.LifecyclePolicySummary{},
		}, nil)

	lister := &OpenSearchServerlessLifecyclePolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOpenSearchServerlessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchServerlessLifecyclePolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchServerlessClient)

	policy := &OpenSearchServerlessLifecyclePolicy{
		svc:  mockClient,
		Name: ptr.String("my-retention-policy"),
		Type: osstypes.LifecyclePolicyTypeRetention,
	}

	mockClient.On("DeleteLifecyclePolicy", mock.Anything, &opensearchserverless.DeleteLifecyclePolicyInput{
		Name: policy.Name,
		Type: osstypes.LifecyclePolicyTypeRetention,
	}).Return(&opensearchserverless.DeleteLifecyclePolicyOutput{}, nil)

	a.NoError(policy.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchServerlessLifecyclePolicy_Properties(t *testing.T) {
	a := assert.New(t)

	policy := OpenSearchServerlessLifecyclePolicy{
		Name: ptr.String("my-retention-policy"),
		Type: osstypes.LifecyclePolicyTypeRetention,
	}

	props := policy.Properties()
	a.Equal("my-retention-policy", props.Get("Name"))
	a.Equal(string(osstypes.LifecyclePolicyTypeRetention), props.Get("Type"))
}

func Test_Mock_OpenSearchServerlessLifecyclePolicy_String(t *testing.T) {
	a := assert.New(t)
	policy := OpenSearchServerlessLifecyclePolicy{Name: ptr.String("my-retention-policy")}
	a.Equal("my-retention-policy", policy.String())
}

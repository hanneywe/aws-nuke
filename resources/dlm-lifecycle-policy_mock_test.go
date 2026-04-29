package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/dlm"
	dlmtypes "github.com/aws/aws-sdk-go-v2/service/dlm/types"
)

func Test_Mock_DLMLifecyclePolicy_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDlmClient)

	mockClient.On("GetLifecyclePolicies", mock.Anything, mock.Anything).
		Return(&dlm.GetLifecyclePoliciesOutput{
			Policies: []dlmtypes.LifecyclePolicySummary{
				{
					PolicyId:    ptr.String("policy-abc123"),
					Description: ptr.String("test-description"),
					Tags:        map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &DLMLifecyclePolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDlmListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*DLMLifecyclePolicy)
	a.Equal("policy-abc123", *r.PolicyID)
	a.Equal("test-description", *r.Description)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DLMLifecyclePolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDlmClient)

	mockClient.On("GetLifecyclePolicies", mock.Anything, mock.Anything).
		Return(&dlm.GetLifecyclePoliciesOutput{
			Policies: []dlmtypes.LifecyclePolicySummary{},
		}, nil)

	lister := &DLMLifecyclePolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDlmListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DLMLifecyclePolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDlmClient)

	r := &DLMLifecyclePolicy{
		svc:      mockClient,
		PolicyID: ptr.String("policy-abc123"),
	}

	mockClient.On("DeleteLifecyclePolicy", mock.Anything,
		&dlm.DeleteLifecyclePolicyInput{
			PolicyId: r.PolicyID,
		}).Return(&dlm.DeleteLifecyclePolicyOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_DLMLifecyclePolicy_Properties(t *testing.T) {
	a := assert.New(t)
	r := &DLMLifecyclePolicy{
		PolicyID:    ptr.String("policy-abc123"),
		Description: ptr.String("test-description"),
	}
	props := r.Properties()
	a.Equal("policy-abc123", props.Get("PolicyID"))
	a.Equal("test-description", props.Get("Description"))
}

func Test_Mock_DLMLifecyclePolicy_String(t *testing.T) {
	a := assert.New(t)
	r := &DLMLifecyclePolicy{
		PolicyID: ptr.String("policy-abc123"),
	}
	a.Equal("policy-abc123", r.String())
}

package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
)

func Test_Mock_BatchSchedulingPolicy_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	mockClient.On("ListSchedulingPolicies", mock.Anything, mock.Anything).
		Return(&batch.ListSchedulingPoliciesOutput{
			SchedulingPolicies: []batchtypes.SchedulingPolicyListingDetail{
				{Arn: ptr.String("arn:aws:batch:us-east-1:123456789012:scheduling-policy/my-policy")},
			},
		}, nil)
	lister := &BatchSchedulingPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchSchedulingPolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	mockClient.On("ListSchedulingPolicies", mock.Anything, mock.Anything).
		Return(&batch.ListSchedulingPoliciesOutput{SchedulingPolicies: []batchtypes.SchedulingPolicyListingDetail{}}, nil)
	lister := &BatchSchedulingPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchSchedulingPolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	r := &BatchSchedulingPolicy{svc: mockClient, Arn: ptr.String("arn:aws:batch:us-east-1:123456789012:scheduling-policy/my-policy")}
	mockClient.On("DeleteSchedulingPolicy", mock.Anything, &batch.DeleteSchedulingPolicyInput{Arn: r.Arn}).
		Return(&batch.DeleteSchedulingPolicyOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchSchedulingPolicy_Properties(t *testing.T) {
	a := assert.New(t)
	policyArn := "arn:aws:batch:us-east-1:123456789012:scheduling-policy/my-policy"
	r := BatchSchedulingPolicy{Arn: ptr.String(policyArn)}
	a.Equal(policyArn, r.Properties().Get("Arn"))
}

func Test_Mock_BatchSchedulingPolicy_String(t *testing.T) {
	a := assert.New(t)
	policyArn := ptr.String("arn:aws:batch:us-east-1:123456789012:scheduling-policy/my-policy")
	a.Equal(*policyArn, (&BatchSchedulingPolicy{Arn: policyArn}).String())
}

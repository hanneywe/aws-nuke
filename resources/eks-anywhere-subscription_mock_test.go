package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func Test_Mock_EKSAnywhereSubscription_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEKSv2Client)
	mockClient.On("ListEksAnywhereSubscriptions", mock.Anything, mock.Anything).
		Return(&eks.ListEksAnywhereSubscriptionsOutput{
			Subscriptions: []ekstypes.EksAnywhereSubscription{
				{
					Id:     ptr.String("sub-12345"),
					Arn:    ptr.String("arn:aws:eks:us-east-1:123456789012:subscription/sub-12345"),
					Status: ptr.String("EXPIRED"),
				},
			},
		}, nil)
	lister := &EKSAnywhereSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEKSv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("sub-12345", resources[0].(*EKSAnywhereSubscription).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_EKSAnywhereSubscription_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEKSv2Client)
	mockClient.On("ListEksAnywhereSubscriptions", mock.Anything, mock.Anything).
		Return(&eks.ListEksAnywhereSubscriptionsOutput{Subscriptions: []ekstypes.EksAnywhereSubscription{}}, nil)
	lister := &EKSAnywhereSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEKSv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EKSAnywhereSubscription_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEKSv2Client)
	r := &EKSAnywhereSubscription{
		svc:    mockClient,
		ID:     ptr.String("sub-12345"),
		Arn:    ptr.String("arn:aws:eks:us-east-1:123456789012:subscription/sub-12345"),
		Status: ptr.String("EXPIRED"),
	}
	mockClient.On("DeleteEksAnywhereSubscription", mock.Anything, &eks.DeleteEksAnywhereSubscriptionInput{
		Id: r.ID,
	}).Return(&eks.DeleteEksAnywhereSubscriptionOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EKSAnywhereSubscription_Filter_Active(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{
		ID:     ptr.String("sub-12345"),
		Status: ptr.String("ACTIVE"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete subscription")
}

func Test_Mock_EKSAnywhereSubscription_Filter_Creating(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{
		ID:     ptr.String("sub-12345"),
		Status: ptr.String("CREATING"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete subscription")
}

func Test_Mock_EKSAnywhereSubscription_Filter_Deleting(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{
		ID:     ptr.String("sub-12345"),
		Status: ptr.String("DELETING"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete subscription")
}

func Test_Mock_EKSAnywhereSubscription_Filter_Inactive(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{
		ID:     ptr.String("sub-12345"),
		Status: ptr.String("INACTIVE"),
	}
	a.NoError(r.Filter())
}

func Test_Mock_EKSAnywhereSubscription_Filter_Expired(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{
		ID:     ptr.String("sub-12345"),
		Status: ptr.String("EXPIRED"),
	}
	a.NoError(r.Filter())
}

func Test_Mock_EKSAnywhereSubscription_Properties(t *testing.T) {
	a := assert.New(t)
	r := EKSAnywhereSubscription{ID: ptr.String("sub-12345"), Arn: ptr.String("arn:aws:eks:us-east-1:123456789012:subscription/sub-12345")}
	a.Equal("sub-12345", r.Properties().Get("Id"))
	a.Equal("arn:aws:eks:us-east-1:123456789012:subscription/sub-12345", r.Properties().Get("Arn"))
}

func Test_Mock_EKSAnywhereSubscription_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("sub-12345", (&EKSAnywhereSubscription{ID: ptr.String("sub-12345")}).String())
}

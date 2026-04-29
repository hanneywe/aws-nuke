package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func Test_Mock_IoTPolicyVersion_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListPolicies", mock.Anything, mock.Anything).
		Return(&iot.ListPoliciesOutput{
			Policies: []iottypes.Policy{
				{PolicyName: ptr.String("test-policy")},
			},
		}, nil)

	mockClient.On("ListPolicyVersions", mock.Anything, mock.Anything).
		Return(&iot.ListPolicyVersionsOutput{
			PolicyVersions: []iottypes.PolicyVersion{
				{VersionId: ptr.String("1"), IsDefaultVersion: false},
				{VersionId: ptr.String("2"), IsDefaultVersion: true},
			},
		}, nil)

	lister := &IoTPolicyVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 2)

	r := resources[0].(*IoTPolicyVersion)
	a.Equal("test-policy", *r.PolicyName)
	a.Equal("1", *r.PolicyVersionID)
	a.False(r.IsDefaultVersion)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTPolicyVersion_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListPolicies", mock.Anything, mock.Anything).
		Return(&iot.ListPoliciesOutput{
			Policies: []iottypes.Policy{},
		}, nil)

	lister := &IoTPolicyVersionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTPolicyVersion_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	r := &IoTPolicyVersion{
		svc:             mockClient,
		PolicyName:      ptr.String("test-policy"),
		PolicyVersionID: ptr.String("1"),
	}

	mockClient.On("DeletePolicyVersion", mock.Anything,
		&iot.DeletePolicyVersionInput{
			PolicyName:      r.PolicyName,
			PolicyVersionId: r.PolicyVersionID,
		}).Return(&iot.DeletePolicyVersionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTPolicyVersion_Properties(t *testing.T) {
	a := assert.New(t)
	r := &IoTPolicyVersion{
		PolicyName:      ptr.String("test-policy"),
		PolicyVersionID: ptr.String("1"),
	}
	props := r.Properties()
	a.Equal("test-policy", props.Get("PolicyName"))
	a.Equal("1", props.Get("PolicyVersionId"))
}

func Test_Mock_IoTPolicyVersion_String(t *testing.T) {
	a := assert.New(t)
	r := &IoTPolicyVersion{
		PolicyVersionID: ptr.String("1"),
	}
	a.Equal("1", r.String())
}

func Test_Mock_IoTPolicyVersion_Filter(t *testing.T) {
	a := assert.New(t)

	r := &IoTPolicyVersion{
		PolicyName:       ptr.String("test-policy"),
		PolicyVersionID:  ptr.String("1"),
		IsDefaultVersion: false,
	}
	a.Nil(r.Filter())

	r.IsDefaultVersion = true
	a.EqualError(r.Filter(), "cannot delete default policy version")
}

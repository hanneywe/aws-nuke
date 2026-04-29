package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func Test_Mock_ECRPullTimeUpdateExclusion_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("ListPullTimeUpdateExclusions", mock.Anything, mock.Anything).
		Return(&ecr.ListPullTimeUpdateExclusionsOutput{
			PullTimeUpdateExclusions: []string{
				"arn:aws:iam::123456789012:role/my-role",
			},
		}, nil)
	lister := &ECRPullTimeUpdateExclusionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("arn:aws:iam::123456789012:role/my-role", resources[0].(*ECRPullTimeUpdateExclusion).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullTimeUpdateExclusion_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("ListPullTimeUpdateExclusions", mock.Anything, mock.Anything).
		Return(&ecr.ListPullTimeUpdateExclusionsOutput{
			PullTimeUpdateExclusions: []string{},
		}, nil)
	lister := &ECRPullTimeUpdateExclusionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullTimeUpdateExclusion_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRPullTimeUpdateExclusion{
		svc:          mockClient,
		PrincipalArn: ptr.String("arn:aws:iam::123456789012:role/my-role"),
	}
	mockClient.On("DeregisterPullTimeUpdateExclusion", mock.Anything, &ecr.DeregisterPullTimeUpdateExclusionInput{
		PrincipalArn: r.PrincipalArn,
	}).Return(&ecr.DeregisterPullTimeUpdateExclusionOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullTimeUpdateExclusion_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRPullTimeUpdateExclusion{
		PrincipalArn: ptr.String("arn:aws:iam::123456789012:role/my-role"),
	}
	props := r.Properties()
	a.Equal("arn:aws:iam::123456789012:role/my-role", props.Get("PrincipalArn"))
}

func Test_Mock_ECRPullTimeUpdateExclusion_String(t *testing.T) {
	a := assert.New(t)
	r := &ECRPullTimeUpdateExclusion{PrincipalArn: ptr.String("arn:aws:iam::123456789012:role/my-role")}
	a.Equal("arn:aws:iam::123456789012:role/my-role", r.String())
}

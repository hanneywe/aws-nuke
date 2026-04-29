package resources

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/iam" //nolint:staticcheck

	"github.com/ekristen/aws-nuke/v3/mocks/mock_iamiface"
)

func Test_Mock_IAMUserPermissionBoundary_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_iamiface.NewMockIAMAPI(ctrl)

	r := &IAMUserPermissionBoundary{
		svc:                    mockSvc,
		UserName:               ptr.String("test-user"),
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
	}

	mockSvc.EXPECT().DeleteUserPermissionsBoundary(&iam.DeleteUserPermissionsBoundaryInput{
		UserName: r.UserName,
	}).Return(&iam.DeleteUserPermissionsBoundaryOutput{}, nil)

	err := r.Remove(context.TODO())
	a.Nil(err)
}

func Test_Mock_IAMUserPermissionBoundary_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	r := IAMUserPermissionBoundary{
		UserName:               ptr.String("test-user"),
		UserPath:               ptr.String("/"),
		UserCreateDate:         &now,
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
		UserTags: []*iam.Tag{
			{Key: ptr.String("env"), Value: ptr.String("test")},
		},
	}

	props := r.Properties()
	a.Equal("test-user", props.Get("UserName"))
	a.Equal("/", props.Get("UserPath"))
	a.Equal("arn:aws:iam::123456789012:policy/boundary", props.Get("PermissionsBoundaryARN"))
}

func Test_Mock_IAMUserPermissionBoundary_String(t *testing.T) {
	a := assert.New(t)

	r := IAMUserPermissionBoundary{
		UserName:               ptr.String("test-user"),
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
	}

	a.Equal("test-user -> arn:aws:iam::123456789012:policy/boundary", r.String())
}

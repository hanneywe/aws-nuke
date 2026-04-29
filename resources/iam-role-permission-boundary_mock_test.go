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

func Test_Mock_IAMRolePermissionBoundary_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_iamiface.NewMockIAMAPI(ctrl)

	r := &IAMRolePermissionBoundary{
		svc:                    mockSvc,
		RoleName:               ptr.String("test-role"),
		RolePath:               ptr.String("/"),
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
	}

	mockSvc.EXPECT().DeleteRolePermissionsBoundary(&iam.DeleteRolePermissionsBoundaryInput{
		RoleName: r.RoleName,
	}).Return(&iam.DeleteRolePermissionsBoundaryOutput{}, nil)

	err := r.Remove(context.TODO())
	a.Nil(err)
}

func Test_Mock_IAMRolePermissionBoundary_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	r := IAMRolePermissionBoundary{
		RoleName:               ptr.String("test-role"),
		RolePath:               ptr.String("/"),
		RoleCreateDate:         &now,
		RoleLastUsed:           &now,
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
		RoleTags: []*iam.Tag{
			{Key: ptr.String("env"), Value: ptr.String("test")},
		},
	}

	props := r.Properties()
	a.Equal("test-role", props.Get("RoleName"))
	a.Equal("/", props.Get("RolePath"))
	a.Equal("arn:aws:iam::123456789012:policy/boundary", props.Get("PermissionsBoundaryARN"))
}

func Test_Mock_IAMRolePermissionBoundary_String(t *testing.T) {
	a := assert.New(t)

	r := IAMRolePermissionBoundary{
		RoleName:               ptr.String("test-role"),
		PermissionsBoundaryARN: ptr.String("arn:aws:iam::123456789012:policy/boundary"),
	}

	a.Equal("test-role -> arn:aws:iam::123456789012:policy/boundary", r.String())
}

func Test_Mock_IAMRolePermissionBoundary_Filter_ServiceRole(t *testing.T) {
	a := assert.New(t)

	r := IAMRolePermissionBoundary{
		RolePath: ptr.String("/aws-service-role/elasticloadbalancing.amazonaws.com/"),
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "service-linked")
}

func Test_Mock_IAMRolePermissionBoundary_Filter_SSORole(t *testing.T) {
	a := assert.New(t)

	r := IAMRolePermissionBoundary{
		RolePath: ptr.String("/aws-reserved/sso.amazonaws.com/us-east-1/"),
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "SSO")
}

func Test_Mock_IAMRolePermissionBoundary_Filter_Normal(t *testing.T) {
	a := assert.New(t)

	r := IAMRolePermissionBoundary{
		RolePath: ptr.String("/"),
	}

	err := r.Filter()
	a.NoError(err)
}

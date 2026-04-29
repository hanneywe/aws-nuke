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

func Test_Mock_IAMPolicyVersion_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_iamiface.NewMockIAMAPI(ctrl)

	r := &IAMPolicyVersion{
		svc:       mockSvc,
		PolicyARN: ptr.String("arn:aws:iam::123456789012:policy/test-policy"),
		VersionID: ptr.String("v2"),
	}

	mockSvc.EXPECT().DeletePolicyVersion(&iam.DeletePolicyVersionInput{
		PolicyArn: r.PolicyARN,
		VersionId: r.VersionID,
	}).Return(&iam.DeletePolicyVersionOutput{}, nil)

	err := r.Remove(context.TODO())
	a.Nil(err)
}

func Test_Mock_IAMPolicyVersion_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	r := IAMPolicyVersion{
		PolicyARN:  ptr.String("arn:aws:iam::123456789012:policy/test-policy"),
		PolicyName: ptr.String("test-policy"),
		VersionID:  ptr.String("v2"),
		CreateDate: &now,
	}

	props := r.Properties()
	a.Equal("arn:aws:iam::123456789012:policy/test-policy", props.Get("PolicyARN"))
	a.Equal("test-policy", props.Get("PolicyName"))
	a.Equal("v2", props.Get("VersionID"))
}

func Test_Mock_IAMPolicyVersion_String(t *testing.T) {
	a := assert.New(t)

	r := IAMPolicyVersion{
		PolicyName: ptr.String("test-policy"),
		VersionID:  ptr.String("v2"),
	}

	a.Equal("test-policy -> v2", r.String())
}

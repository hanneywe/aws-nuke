package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/xray"
	xraytypes "github.com/aws/aws-sdk-go-v2/service/xray/types"
)

func Test_Mock_XRayEncryptionConfig_List_KMS(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockXRayClient)

	mockClient.On("GetEncryptionConfig", mock.Anything, mock.Anything).
		Return(&xray.GetEncryptionConfigOutput{
			EncryptionConfig: &xraytypes.EncryptionConfig{
				Type:  xraytypes.EncryptionTypeKms,
				KeyId: ptr.String("arn:aws:kms:us-east-1:123456789012:key/my-key"),
			},
		}, nil)

	lister := &XRayEncryptionConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*XRayEncryptionConfig)
	a.Equal(xraytypes.EncryptionTypeKms, r.Type)
	a.Equal("arn:aws:kms:us-east-1:123456789012:key/my-key", *r.KeyID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_XRayEncryptionConfig_List_Default(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockXRayClient)

	mockClient.On("GetEncryptionConfig", mock.Anything, mock.Anything).
		Return(&xray.GetEncryptionConfigOutput{
			EncryptionConfig: &xraytypes.EncryptionConfig{
				Type: xraytypes.EncryptionTypeNone,
			},
		}, nil)

	lister := &XRayEncryptionConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_XRayEncryptionConfig_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockXRayClient)

	r := &XRayEncryptionConfig{
		svc:   mockClient,
		Type:  xraytypes.EncryptionTypeKms,
		KeyID: ptr.String("arn:aws:kms:us-east-1:123456789012:key/my-key"),
	}

	mockClient.On("PutEncryptionConfig", mock.Anything,
		&xray.PutEncryptionConfigInput{
			Type: xraytypes.EncryptionTypeNone,
		}).Return(&xray.PutEncryptionConfigOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_XRayEncryptionConfig_Properties(t *testing.T) {
	a := assert.New(t)
	r := &XRayEncryptionConfig{
		Type:  xraytypes.EncryptionTypeKms,
		KeyID: ptr.String("arn:aws:kms:us-east-1:123456789012:key/my-key"),
	}
	props := r.Properties()
	a.Equal("KMS", props.Get("Type"))
	a.Equal("arn:aws:kms:us-east-1:123456789012:key/my-key", props.Get("KeyID"))
}

func Test_Mock_XRayEncryptionConfig_String(t *testing.T) {
	a := assert.New(t)
	r := &XRayEncryptionConfig{
		Type: xraytypes.EncryptionTypeKms,
	}
	a.Equal("KMS", r.String())
}

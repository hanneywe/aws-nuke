package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3control"
	s3controltypes "github.com/aws/aws-sdk-go-v2/service/s3control/types"
)

func Test_Mock_S3PublicAccessBlock_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3AccountClient)

	mockClient.On("GetPublicAccessBlock", mock.Anything, mock.Anything).
		Return(&s3control.GetPublicAccessBlockOutput{
			PublicAccessBlockConfiguration: &s3controltypes.PublicAccessBlockConfiguration{
				BlockPublicAcls:       ptr.Bool(true),
				BlockPublicPolicy:     ptr.Bool(true),
				IgnorePublicAcls:      ptr.Bool(false),
				RestrictPublicBuckets: ptr.Bool(true),
			},
		}, nil)

	lister := &S3PublicAccessBlockLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3AccountListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*S3PublicAccessBlock)
	a.Equal(true, *r.BlockPublicAcls)
	a.Equal(true, *r.BlockPublicPolicy)
	a.Equal(false, *r.IgnorePublicAcls)
	a.Equal(true, *r.RestrictPublicBuckets)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3PublicAccessBlock_List_NoConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3AccountClient)

	mockClient.On("GetPublicAccessBlock", mock.Anything, mock.Anything).
		Return(
			(*s3control.GetPublicAccessBlockOutput)(nil),
			&s3controltypes.NoSuchPublicAccessBlockConfiguration{
				Message: ptr.String("no config"),
			},
		)

	lister := &S3PublicAccessBlockLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3AccountListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3PublicAccessBlock_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3AccountClient)

	r := &S3PublicAccessBlock{
		svc:       mockClient,
		accountID: ptr.String("123456789012"),
	}

	mockClient.On("DeletePublicAccessBlock", mock.Anything,
		&s3control.DeletePublicAccessBlockInput{
			AccountId: r.accountID,
		}).Return(&s3control.DeletePublicAccessBlockOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3PublicAccessBlock_Properties(t *testing.T) {
	a := assert.New(t)
	r := &S3PublicAccessBlock{
		BlockPublicAcls:       ptr.Bool(true),
		BlockPublicPolicy:     ptr.Bool(false),
		IgnorePublicAcls:      ptr.Bool(true),
		RestrictPublicBuckets: ptr.Bool(false),
	}
	props := r.Properties()
	a.Equal("true", props.Get("BlockPublicAcls"))
	a.Equal("false", props.Get("BlockPublicPolicy"))
	a.Equal("true", props.Get("IgnorePublicAcls"))
	a.Equal("false", props.Get("RestrictPublicBuckets"))
}

func Test_Mock_S3PublicAccessBlock_String(t *testing.T) {
	a := assert.New(t)
	r := &S3PublicAccessBlock{}
	a.Equal("S3PublicAccessBlock", r.String())
}

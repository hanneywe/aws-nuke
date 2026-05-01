package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

func Test_Mock_SageMakerImage_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListImages", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListImagesOutput{
				Images: []sagemakertypes.Image{
					{
						ImageName:   ptr.String("test-image"),
						ImageArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:image/test-image"),
						ImageStatus: sagemakertypes.ImageStatusCreated,
					},
				},
			}, nil,
		)

	lister := &SageMakerImageLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	image := resources[0].(*SageMakerImage)
	assertions.Equal("test-image", *image.ImageName)
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:image/test-image", *image.ImageArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerImage_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListImages", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListImagesOutput{
				Images: []sagemakertypes.Image{},
			}, nil,
		)

	lister := &SageMakerImageLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerImage_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	image := &SageMakerImage{
		svc:         mockClient,
		ImageName:   ptr.String("test-image"),
		ImageArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:image/test-image"),
		ImageStatus: sagemakertypes.ImageStatusCreated,
	}

	mockClient.
		On("DeleteImage", mock.Anything,
			&sagemaker.DeleteImageInput{
				ImageName: image.ImageName,
			},
		).
		Return(&sagemaker.DeleteImageOutput{}, nil)

	err := image.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerImage_Properties(t *testing.T) {
	assertions := assert.New(t)

	image := SageMakerImage{
		ImageName:   ptr.String("test-image"),
		ImageArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:image/test-image"),
		ImageStatus: sagemakertypes.ImageStatusCreated,
	}

	properties := image.Properties()

	assertions.Equal("test-image", properties.Get("ImageName"))
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:image/test-image", properties.Get("ImageArn"))
}

func Test_Mock_SageMakerImage_String(t *testing.T) {
	assertions := assert.New(t)

	image := SageMakerImage{
		ImageName: ptr.String("test-image"),
	}

	assertions.Equal("test-image", image.String())
}

func Test_Mock_SageMakerImage_Filter(t *testing.T) {
	assertions := assert.New(t)

	deletingImage := SageMakerImage{
		ImageName:   ptr.String("deleting-image"),
		ImageStatus: sagemakertypes.ImageStatusDeleting,
	}
	err := deletingImage.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleting")

	createdImage := SageMakerImage{
		ImageName:   ptr.String("created-image"),
		ImageStatus: sagemakertypes.ImageStatusCreated,
	}
	err = createdImage.Filter()
	assertions.NoError(err)
}

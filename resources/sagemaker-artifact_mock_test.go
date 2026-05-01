package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testSageMakerV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_SageMakerArtifact_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListArtifacts", mock.Anything, mock.Anything).
		Return(&sagemaker.ListArtifactsOutput{
			ArtifactSummaries: []sagemakertypes.ArtifactSummary{
				{
					ArtifactArn: ptr.String("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact"),
				},
			},
		}, nil)

	lister := &SageMakerArtifactLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	artifact := resources[0].(*SageMakerArtifact)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact", *artifact.ArtifactArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerArtifact_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListArtifacts", mock.Anything, mock.Anything).
		Return(&sagemaker.ListArtifactsOutput{
			ArtifactSummaries: []sagemakertypes.ArtifactSummary{},
		}, nil)

	lister := &SageMakerArtifactLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerArtifact_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	artifact := &SageMakerArtifact{
		svc:         mockClient,
		ArtifactArn: ptr.String("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact"),
	}

	mockClient.On("DeleteArtifact", mock.Anything, &sagemaker.DeleteArtifactInput{
		ArtifactArn: artifact.ArtifactArn,
	}).Return(&sagemaker.DeleteArtifactOutput{}, nil)

	a.NoError(artifact.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerArtifact_Properties(t *testing.T) {
	a := assert.New(t)

	artifact := SageMakerArtifact{
		ArtifactArn: ptr.String("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact"),
	}

	props := artifact.Properties()
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact", props.Get("ArtifactArn"))
}

func Test_Mock_SageMakerArtifact_String(t *testing.T) {
	a := assert.New(t)
	artifact := SageMakerArtifact{ArtifactArn: ptr.String("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact")}
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:artifact/my-artifact", artifact.String())
}

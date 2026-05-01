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

func Test_Mock_SageMakerMlflowTrackingServer_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListMlflowTrackingServers", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListMlflowTrackingServersOutput{
				TrackingServerSummaries: []sagemakertypes.TrackingServerSummary{
					{
						TrackingServerName:   ptr.String("test-tracking-server"),
						TrackingServerArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-tracking-server/test-tracking-server"),
						TrackingServerStatus: sagemakertypes.TrackingServerStatusCreated,
					},
				},
			}, nil,
		)

	lister := &SageMakerMlflowTrackingServerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	trackingServer := resources[0].(*SageMakerMlflowTrackingServer)
	assertions.Equal("test-tracking-server", *trackingServer.TrackingServerName)
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:mlflow-tracking-server/test-tracking-server", *trackingServer.TrackingServerArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowTrackingServer_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListMlflowTrackingServers", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListMlflowTrackingServersOutput{
				TrackingServerSummaries: []sagemakertypes.TrackingServerSummary{},
			}, nil,
		)

	lister := &SageMakerMlflowTrackingServerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowTrackingServer_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	trackingServer := &SageMakerMlflowTrackingServer{
		svc:                  mockClient,
		TrackingServerName:   ptr.String("test-tracking-server"),
		TrackingServerArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-tracking-server/test-tracking-server"),
		TrackingServerStatus: sagemakertypes.TrackingServerStatusCreated,
	}

	mockClient.
		On("DeleteMlflowTrackingServer", mock.Anything,
			&sagemaker.DeleteMlflowTrackingServerInput{
				TrackingServerName: trackingServer.TrackingServerName,
			},
		).
		Return(&sagemaker.DeleteMlflowTrackingServerOutput{}, nil)

	err := trackingServer.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowTrackingServer_Properties(t *testing.T) {
	assertions := assert.New(t)

	trackingServer := SageMakerMlflowTrackingServer{
		TrackingServerName:   ptr.String("test-tracking-server"),
		TrackingServerArn:    ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-tracking-server/test-tracking-server"),
		TrackingServerStatus: sagemakertypes.TrackingServerStatusCreated,
	}

	properties := trackingServer.Properties()

	assertions.Equal("test-tracking-server", properties.Get("TrackingServerName"))
	assertions.Equal(
		"arn:aws:sagemaker:us-east-1:123456789012:mlflow-tracking-server/test-tracking-server",
		properties.Get("TrackingServerArn"),
	)
}

func Test_Mock_SageMakerMlflowTrackingServer_String(t *testing.T) {
	assertions := assert.New(t)

	trackingServer := SageMakerMlflowTrackingServer{
		TrackingServerName: ptr.String("test-tracking-server"),
	}

	assertions.Equal("test-tracking-server", trackingServer.String())
}

func Test_Mock_SageMakerMlflowTrackingServer_Filter(t *testing.T) {
	assertions := assert.New(t)

	deletingServer := SageMakerMlflowTrackingServer{
		TrackingServerName:   ptr.String("deleting-server"),
		TrackingServerStatus: sagemakertypes.TrackingServerStatusDeleting,
	}
	err := deletingServer.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleting")

	createdServer := SageMakerMlflowTrackingServer{
		TrackingServerName:   ptr.String("created-server"),
		TrackingServerStatus: sagemakertypes.TrackingServerStatusCreated,
	}
	err = createdServer.Filter()
	assertions.NoError(err)
}

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

func Test_Mock_SageMakerHub_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListHubs", mock.Anything, mock.Anything).
		Return(&sagemaker.ListHubsOutput{
			HubSummaries: []sagemakertypes.HubInfo{
				{
					HubName: ptr.String("my-hub"),
					HubArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:hub/my-hub"),
				},
			},
		}, nil)

	lister := &SageMakerHubLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	hub := resources[0].(*SageMakerHub)
	a.Equal("my-hub", *hub.HubName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:hub/my-hub", *hub.HubArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHub_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListHubs", mock.Anything, mock.Anything).
		Return(&sagemaker.ListHubsOutput{
			HubSummaries: []sagemakertypes.HubInfo{},
		}, nil)

	lister := &SageMakerHubLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHub_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	hub := &SageMakerHub{
		svc:     mockClient,
		HubName: ptr.String("my-hub"),
	}

	mockClient.On("DeleteHub", mock.Anything, &sagemaker.DeleteHubInput{
		HubName: hub.HubName,
	}).Return(&sagemaker.DeleteHubOutput{}, nil)

	a.NoError(hub.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHub_Properties(t *testing.T) {
	a := assert.New(t)

	hub := SageMakerHub{
		HubName: ptr.String("my-hub"),
		HubArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:hub/my-hub"),
	}

	props := hub.Properties()
	a.Equal("my-hub", props.Get("HubName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:hub/my-hub", props.Get("HubArn"))
}

func Test_Mock_SageMakerHub_String(t *testing.T) {
	a := assert.New(t)
	hub := SageMakerHub{HubName: ptr.String("my-hub")}
	a.Equal("my-hub", hub.String())
}

func Test_Mock_SageMakerHub_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := SageMakerHub{
		HubName: ptr.String("SageMakerPublicHub"),
		HubArn:  ptr.String("arn:aws:sagemaker:us-east-1:aws:hub/SageMakerPublicHub"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed hub")
}

func Test_Mock_SageMakerHub_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := SageMakerHub{
		HubName: ptr.String("my-hub"),
		HubArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:hub/my-hub"),
	}
	a.NoError(r.Filter())
}

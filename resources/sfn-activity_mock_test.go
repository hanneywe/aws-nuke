package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func Test_Mock_SFNActivity_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSFNv2Client)

	mockClient.On("ListActivities", mock.Anything, mock.Anything).
		Return(&sfn.ListActivitiesOutput{
			Activities: []sfntypes.ActivityListItem{
				{
					Name:        ptr.String("my-activity"),
					ActivityArn: ptr.String("arn:aws:states:us-east-1:123456789012:activity:my-activity"),
				},
			},
		}, nil)

	lister := &SFNActivityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSFNv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SFNActivity)
	a.Equal("my-activity", *r.Name)
	a.Equal("arn:aws:states:us-east-1:123456789012:activity:my-activity", *r.ActivityArn)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SFNActivity_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSFNv2Client)

	mockClient.On("ListActivities", mock.Anything, mock.Anything).
		Return(&sfn.ListActivitiesOutput{
			Activities: []sfntypes.ActivityListItem{},
		}, nil)

	lister := &SFNActivityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSFNv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SFNActivity_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSFNv2Client)

	r := &SFNActivity{
		svc:         mockClient,
		Name:        ptr.String("my-activity"),
		ActivityArn: ptr.String("arn:aws:states:us-east-1:123456789012:activity:my-activity"),
	}

	mockClient.On("DeleteActivity", mock.Anything,
		&sfn.DeleteActivityInput{
			ActivityArn: r.ActivityArn,
		}).Return(&sfn.DeleteActivityOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SFNActivity_Properties(t *testing.T) {
	a := assert.New(t)
	r := &SFNActivity{
		Name:        ptr.String("my-activity"),
		ActivityArn: ptr.String("arn:aws:states:us-east-1:123456789012:activity:my-activity"),
	}
	props := r.Properties()
	a.Equal("my-activity", props.Get("Name"))
	a.Equal("arn:aws:states:us-east-1:123456789012:activity:my-activity",
		props.Get("ActivityArn"))
}

func Test_Mock_SFNActivity_String(t *testing.T) {
	a := assert.New(t)
	r := &SFNActivity{
		Name: ptr.String("my-activity"),
	}
	a.Equal("my-activity", r.String())
}

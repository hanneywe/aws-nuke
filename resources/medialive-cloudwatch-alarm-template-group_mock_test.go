package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	mltypes "github.com/aws/aws-sdk-go-v2/service/medialive/types"
)

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListCloudWatchAlarmTemplateGroups", mock.Anything, mock.Anything).
		Return(&medialive.ListCloudWatchAlarmTemplateGroupsOutput{
			CloudWatchAlarmTemplateGroups: []mltypes.CloudWatchAlarmTemplateGroupSummary{
				{
					Id:   ptr.String("cwatg-123"),
					Name: ptr.String("my-alarm-group"),
				},
			},
		}, nil)

	lister := &MediaLiveCloudWatchAlarmTemplateGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveCloudWatchAlarmTemplateGroup)
	a.Equal("cwatg-123", *r.ID)
	a.Equal("my-alarm-group", *r.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListCloudWatchAlarmTemplateGroups", mock.Anything, mock.Anything).
		Return(&medialive.ListCloudWatchAlarmTemplateGroupsOutput{
			CloudWatchAlarmTemplateGroups: []mltypes.CloudWatchAlarmTemplateGroupSummary{},
		}, nil)

	lister := &MediaLiveCloudWatchAlarmTemplateGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveCloudWatchAlarmTemplateGroup{
		svc: mockClient,
		ID:  ptr.String("cwatg-123"),
	}

	mockClient.On("DeleteCloudWatchAlarmTemplateGroup", mock.Anything, &medialive.DeleteCloudWatchAlarmTemplateGroupInput{
		Identifier: r.ID,
	}).Return(&medialive.DeleteCloudWatchAlarmTemplateGroupOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_Properties(t *testing.T) {
	a := assert.New(t)

	r := MediaLiveCloudWatchAlarmTemplateGroup{
		ID:   ptr.String("cwatg-123"),
		Name: ptr.String("my-alarm-group"),
	}

	props := r.Properties()
	a.Equal("cwatg-123", props.Get("ID"))
	a.Equal("my-alarm-group", props.Get("Name"))
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_String(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveCloudWatchAlarmTemplateGroup{ID: ptr.String("cwatg-123")}
	a.Equal("cwatg-123", r.String())
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveCloudWatchAlarmTemplateGroup{Name: ptr.String("AWS-MediaConnectInputFailoverWorkflows")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed template group")
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplateGroup_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveCloudWatchAlarmTemplateGroup{Name: ptr.String("my-alarm-group")}
	a.NoError(r.Filter())
}

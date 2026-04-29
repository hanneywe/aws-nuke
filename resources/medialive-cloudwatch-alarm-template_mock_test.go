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

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListCloudWatchAlarmTemplates", mock.Anything, mock.Anything).
		Return(&medialive.ListCloudWatchAlarmTemplatesOutput{
			CloudWatchAlarmTemplates: []mltypes.CloudWatchAlarmTemplateSummary{
				{
					Id:      ptr.String("12345"),
					Name:    ptr.String("my-alarm"),
					GroupId: ptr.String("67890"),
				},
			},
		}, nil)

	lister := &MediaLiveCloudWatchAlarmTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveCloudWatchAlarmTemplate)
	a.Equal("12345", *r.ID)
	a.Equal("my-alarm", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListCloudWatchAlarmTemplates", mock.Anything, mock.Anything).
		Return(&medialive.ListCloudWatchAlarmTemplatesOutput{
			CloudWatchAlarmTemplates: []mltypes.CloudWatchAlarmTemplateSummary{},
		}, nil)

	lister := &MediaLiveCloudWatchAlarmTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveCloudWatchAlarmTemplate{
		svc: mockClient,
		ID:  ptr.String("12345"),
	}

	mockClient.On("DeleteCloudWatchAlarmTemplate", mock.Anything,
		&medialive.DeleteCloudWatchAlarmTemplateInput{
			Identifier: r.ID,
		}).Return(&medialive.DeleteCloudWatchAlarmTemplateOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	r := &MediaLiveCloudWatchAlarmTemplate{
		ID:      ptr.String("12345"),
		Name:    ptr.String("my-alarm"),
		GroupID: ptr.String("67890"),
	}
	props := r.Properties()
	a.Equal("12345", props.Get("ID"))
	a.Equal("my-alarm", props.Get("Name"))
	a.Equal("67890", props.Get("GroupID"))
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_String(t *testing.T) {
	a := assert.New(t)
	r := &MediaLiveCloudWatchAlarmTemplate{ID: ptr.String("12345")}
	a.Equal("12345", r.String())
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := &MediaLiveCloudWatchAlarmTemplate{
		ID:   ptr.String("99999"),
		Name: ptr.String("AWS-Default-Alarm"),
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "AWS-managed")
}

func Test_Mock_MediaLiveCloudWatchAlarmTemplate_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := &MediaLiveCloudWatchAlarmTemplate{
		ID:   ptr.String("12345"),
		Name: ptr.String("my-alarm"),
	}
	a.NoError(r.Filter())
}

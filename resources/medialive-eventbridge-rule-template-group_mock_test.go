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

func Test_Mock_MediaLiveEventBridgeRuleTemplateGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListEventBridgeRuleTemplateGroups", mock.Anything, mock.Anything).
		Return(&medialive.ListEventBridgeRuleTemplateGroupsOutput{
			EventBridgeRuleTemplateGroups: []mltypes.EventBridgeRuleTemplateGroupSummary{
				{
					Id:   ptr.String("ebrtg-123"),
					Name: ptr.String("my-rule-group"),
				},
			},
		}, nil)

	lister := &MediaLiveEventBridgeRuleTemplateGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveEventBridgeRuleTemplateGroup)
	a.Equal("ebrtg-123", *r.ID)
	a.Equal("my-rule-group", *r.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplateGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListEventBridgeRuleTemplateGroups", mock.Anything, mock.Anything).
		Return(&medialive.ListEventBridgeRuleTemplateGroupsOutput{
			EventBridgeRuleTemplateGroups: []mltypes.EventBridgeRuleTemplateGroupSummary{},
		}, nil)

	lister := &MediaLiveEventBridgeRuleTemplateGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplateGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveEventBridgeRuleTemplateGroup{
		svc: mockClient,
		ID:  ptr.String("ebrtg-123"),
	}

	mockClient.On("DeleteEventBridgeRuleTemplateGroup", mock.Anything, &medialive.DeleteEventBridgeRuleTemplateGroupInput{
		Identifier: r.ID,
	}).Return(&medialive.DeleteEventBridgeRuleTemplateGroupOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplateGroup_Properties(t *testing.T) {
	a := assert.New(t)

	r := MediaLiveEventBridgeRuleTemplateGroup{
		ID:   ptr.String("ebrtg-123"),
		Name: ptr.String("my-rule-group"),
	}

	props := r.Properties()
	a.Equal("ebrtg-123", props.Get("ID"))
	a.Equal("my-rule-group", props.Get("Name"))
}

func Test_Mock_MediaLiveEventBridgeRuleTemplateGroup_String(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveEventBridgeRuleTemplateGroup{Name: ptr.String("my-rule-group")}
	a.Equal("my-rule-group", r.String())
}

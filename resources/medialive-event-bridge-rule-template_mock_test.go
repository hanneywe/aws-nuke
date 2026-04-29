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

func Test_Mock_MediaLiveEventBridgeRuleTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListEventBridgeRuleTemplates", mock.Anything, mock.Anything).
		Return(&medialive.ListEventBridgeRuleTemplatesOutput{
			EventBridgeRuleTemplates: []mltypes.EventBridgeRuleTemplateSummary{
				{
					Id:   ptr.String("ebrt-123"),
					Name: ptr.String("my-rule-template"),
				},
			},
		}, nil)

	lister := &MediaLiveEventBridgeRuleTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*MediaLiveEventBridgeRuleTemplate)
	a.Equal("ebrt-123", *r.ID)
	a.Equal("my-rule-template", *r.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	mockClient.On("ListEventBridgeRuleTemplates", mock.Anything, mock.Anything).
		Return(&medialive.ListEventBridgeRuleTemplatesOutput{
			EventBridgeRuleTemplates: []mltypes.EventBridgeRuleTemplateSummary{},
		}, nil)

	lister := &MediaLiveEventBridgeRuleTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaLiveListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaLiveClient)

	r := &MediaLiveEventBridgeRuleTemplate{
		svc: mockClient,
		ID:  ptr.String("ebrt-123"),
	}

	mockClient.On("DeleteEventBridgeRuleTemplate", mock.Anything, &medialive.DeleteEventBridgeRuleTemplateInput{
		Identifier: r.ID,
	}).Return(&medialive.DeleteEventBridgeRuleTemplateOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaLiveEventBridgeRuleTemplate_Properties(t *testing.T) {
	a := assert.New(t)

	r := MediaLiveEventBridgeRuleTemplate{
		ID:   ptr.String("ebrt-123"),
		Name: ptr.String("my-rule-template"),
	}

	props := r.Properties()
	a.Equal("ebrt-123", props.Get("ID"))
	a.Equal("my-rule-template", props.Get("Name"))
}

func Test_Mock_MediaLiveEventBridgeRuleTemplate_String(t *testing.T) {
	a := assert.New(t)
	r := MediaLiveEventBridgeRuleTemplate{ID: ptr.String("ebrt-123")}
	a.Equal("ebrt-123", r.String())
}

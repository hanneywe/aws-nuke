package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
	mailmanagertypes "github.com/aws/aws-sdk-go-v2/service/mailmanager/types"
)

func Test_Mock_MailManagerAddonInstance_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddonInstances", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddonInstancesOutput{
			AddonInstances: []mailmanagertypes.AddonInstance{
				{AddonInstanceId: ptr.String("ai-12345")},
			},
		}, nil)
	lister := &MailManagerAddonInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("ai-12345", *resources[0].(*MailManagerAddonInstance).AddonInstanceID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonInstance_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddonInstances", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddonInstancesOutput{AddonInstances: []mailmanagertypes.AddonInstance{}}, nil)
	lister := &MailManagerAddonInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonInstance_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	ai := &MailManagerAddonInstance{svc: mockClient, AddonInstanceID: ptr.String("ai-12345")}
	mockClient.On("DeleteAddonInstance", mock.Anything, &mailmanager.DeleteAddonInstanceInput{AddonInstanceId: ai.AddonInstanceID}).
		Return(&mailmanager.DeleteAddonInstanceOutput{}, nil)
	a.NoError(ai.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonInstance_Properties(t *testing.T) {
	a := assert.New(t)
	ai := MailManagerAddonInstance{AddonInstanceID: ptr.String("ai-12345")}
	a.Equal("ai-12345", ai.Properties().Get("AddonInstanceId"))
}

func Test_Mock_MailManagerAddonInstance_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("ai-12345", (&MailManagerAddonInstance{AddonInstanceID: ptr.String("ai-12345")}).String())
}

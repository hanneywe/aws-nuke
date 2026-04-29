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

func Test_Mock_MailManagerAddonSubscription_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddonSubscriptions", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddonSubscriptionsOutput{
			AddonSubscriptions: []mailmanagertypes.AddonSubscription{
				{AddonSubscriptionId: ptr.String("as-12345")},
			},
		}, nil)
	lister := &MailManagerAddonSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("as-12345", *resources[0].(*MailManagerAddonSubscription).AddonSubscriptionID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonSubscription_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddonSubscriptions", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddonSubscriptionsOutput{AddonSubscriptions: []mailmanagertypes.AddonSubscription{}}, nil)
	lister := &MailManagerAddonSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonSubscription_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	as := &MailManagerAddonSubscription{svc: mockClient, AddonSubscriptionID: ptr.String("as-12345")}
	mockClient.On("DeleteAddonSubscription", mock.Anything,
		&mailmanager.DeleteAddonSubscriptionInput{AddonSubscriptionId: as.AddonSubscriptionID}).
		Return(&mailmanager.DeleteAddonSubscriptionOutput{}, nil)
	a.NoError(as.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddonSubscription_Properties(t *testing.T) {
	a := assert.New(t)
	as := MailManagerAddonSubscription{AddonSubscriptionID: ptr.String("as-12345")}
	a.Equal("as-12345", as.Properties().Get("AddonSubscriptionId"))
}

func Test_Mock_MailManagerAddonSubscription_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("as-12345", (&MailManagerAddonSubscription{AddonSubscriptionID: ptr.String("as-12345")}).String())
}

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

func Test_Mock_MailManagerAddressList_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddressLists", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddressListsOutput{
			AddressLists: []mailmanagertypes.AddressList{
				{AddressListId: ptr.String("al-12345"), AddressListName: ptr.String("my-list")},
			},
		}, nil)
	lister := &MailManagerAddressListLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	al := resources[0].(*MailManagerAddressList)
	a.Equal("al-12345", *al.AddressListID)
	a.Equal("my-list", *al.AddressListName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddressList_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	mockClient.On("ListAddressLists", mock.Anything, mock.Anything).
		Return(&mailmanager.ListAddressListsOutput{AddressLists: []mailmanagertypes.AddressList{}}, nil)
	lister := &MailManagerAddressListLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddressList_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMailManagerClient)
	al := &MailManagerAddressList{svc: mockClient, AddressListID: ptr.String("al-12345")}
	mockClient.On("DeleteAddressList", mock.Anything, &mailmanager.DeleteAddressListInput{AddressListId: al.AddressListID}).
		Return(&mailmanager.DeleteAddressListOutput{}, nil)
	a.NoError(al.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerAddressList_Properties(t *testing.T) {
	a := assert.New(t)
	al := MailManagerAddressList{AddressListID: ptr.String("al-12345"), AddressListName: ptr.String("my-list")}
	a.Equal("my-list", al.Properties().Get("AddressListName"))
}

func Test_Mock_MailManagerAddressList_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-list", (&MailManagerAddressList{AddressListName: ptr.String("my-list")}).String())
}

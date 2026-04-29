package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/account"
	accounttypes "github.com/aws/aws-sdk-go-v2/service/account/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testAccountListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_AccountAlternateContact_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAccountClient)

	mockClient.
		On("GetAlternateContact", mock.Anything, &account.GetAlternateContactInput{
			AlternateContactType: accounttypes.AlternateContactTypeBilling,
		}).
		Return(&account.GetAlternateContactOutput{
			AlternateContact: &accounttypes.AlternateContact{
				AlternateContactType: accounttypes.AlternateContactTypeBilling,
				Name:                 ptr.String("John Doe"),
				EmailAddress:         ptr.String("john@example.com"),
			},
		}, nil)

	// OPERATIONS and SECURITY return ResourceNotFoundException
	mockClient.
		On("GetAlternateContact", mock.Anything, &account.GetAlternateContactInput{
			AlternateContactType: accounttypes.AlternateContactTypeOperations,
		}).
		Return((*account.GetAlternateContactOutput)(nil),
			&accounttypes.ResourceNotFoundException{Message: ptr.String("not found")})

	mockClient.
		On("GetAlternateContact", mock.Anything, &account.GetAlternateContactInput{
			AlternateContactType: accounttypes.AlternateContactTypeSecurity,
		}).
		Return((*account.GetAlternateContactOutput)(nil),
			&accounttypes.ResourceNotFoundException{Message: ptr.String("not found")})

	lister := &AccountAlternateContactLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testAccountListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	contact := resources[0].(*AccountAlternateContact)
	assertions.Equal("BILLING", *contact.AlternateContactType)
	assertions.Equal("John Doe", *contact.Name)
	assertions.Equal("john@example.com", *contact.EmailAddress)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AccountAlternateContact_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAccountClient)

	// All three types return ResourceNotFoundException
	mockClient.
		On("GetAlternateContact", mock.Anything, mock.Anything).
		Return((*account.GetAlternateContactOutput)(nil),
			&accounttypes.ResourceNotFoundException{Message: ptr.String("not found")})

	lister := &AccountAlternateContactLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testAccountListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AccountAlternateContact_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAccountClient)

	contact := &AccountAlternateContact{
		svc:                  mockClient,
		AlternateContactType: ptr.String("BILLING"),
	}

	mockClient.
		On("DeleteAlternateContact", mock.Anything, &account.DeleteAlternateContactInput{
			AlternateContactType: accounttypes.AlternateContactTypeBilling,
		}).
		Return(&account.DeleteAlternateContactOutput{}, nil)

	err := contact.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AccountAlternateContact_Properties(t *testing.T) {
	assertions := assert.New(t)

	contact := AccountAlternateContact{
		AlternateContactType: ptr.String("BILLING"),
		Name:                 ptr.String("John Doe"),
		EmailAddress:         ptr.String("john@example.com"),
	}

	properties := contact.Properties()

	assertions.Equal("BILLING", properties.Get("AlternateContactType"))
	assertions.Equal("John Doe", properties.Get("Name"))
	assertions.Equal("john@example.com", properties.Get("EmailAddress"))
}

func Test_Mock_AccountAlternateContact_String(t *testing.T) {
	assertions := assert.New(t)

	contact := AccountAlternateContact{
		AlternateContactType: ptr.String("SECURITY"),
	}

	assertions.Equal("SECURITY", contact.String())
}

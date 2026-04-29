package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/notificationscontacts"
	nctypes "github.com/aws/aws-sdk-go-v2/service/notificationscontacts/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockNotificationsContactsClient struct {
	mock.Mock
}

func (m *mockNotificationsContactsClient) ListEmailContacts(
	ctx context.Context, params *notificationscontacts.ListEmailContactsInput,
	_ ...func(*notificationscontacts.Options),
) (*notificationscontacts.ListEmailContactsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*notificationscontacts.ListEmailContactsOutput), args.Error(1)
}

func (m *mockNotificationsContactsClient) DeleteEmailContact(
	ctx context.Context, params *notificationscontacts.DeleteEmailContactInput,
	_ ...func(*notificationscontacts.Options),
) (*notificationscontacts.DeleteEmailContactOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*notificationscontacts.DeleteEmailContactOutput), args.Error(1)
}

var testNotificationsContactsListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_NotificationsContactsEmailContact_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNotificationsContactsClient)

	mockClient.
		On("ListEmailContacts", mock.Anything, mock.Anything).
		Return(&notificationscontacts.ListEmailContactsOutput{
			EmailContacts: []nctypes.EmailContact{
				{
					Arn:     ptr.String("arn:aws:notificationscontacts:us-east-1:123456789012:emailcontact/ec-123"),
					Name:    ptr.String("test-contact"),
					Address: ptr.String("test@example.com"),
					Status:  nctypes.EmailContactStatusActive,
				},
			},
		}, nil)

	lister := &NotificationsContactsEmailContactLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNotificationsContactsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ec := resources[0].(*NotificationsContactsEmailContact)
	a.Equal("test-contact", *ec.Name)
	a.Equal("test@example.com", *ec.Address)
	a.Equal("active", ec.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NotificationsContactsEmailContact_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNotificationsContactsClient)

	mockClient.
		On("ListEmailContacts", mock.Anything, mock.Anything).
		Return(&notificationscontacts.ListEmailContactsOutput{
			EmailContacts: []nctypes.EmailContact{},
		}, nil)

	lister := &NotificationsContactsEmailContactLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNotificationsContactsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NotificationsContactsEmailContact_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNotificationsContactsClient)

	ec := &NotificationsContactsEmailContact{
		svc:  mockClient,
		ARN:  ptr.String("arn:aws:notificationscontacts:us-east-1:123456789012:emailcontact/ec-123"),
		Name: ptr.String("test-contact"),
	}

	mockClient.
		On("DeleteEmailContact", mock.Anything, &notificationscontacts.DeleteEmailContactInput{
			Arn: ec.ARN,
		}).
		Return(&notificationscontacts.DeleteEmailContactOutput{}, nil)

	err := ec.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NotificationsContactsEmailContact_Properties(t *testing.T) {
	a := assert.New(t)

	ec := &NotificationsContactsEmailContact{
		ARN:     ptr.String("arn:aws:notificationscontacts:us-east-1:123456789012:emailcontact/ec-123"),
		Name:    ptr.String("test-contact"),
		Address: ptr.String("test@example.com"),
		Status:  "ACTIVE",
	}

	props := ec.Properties()
	a.Equal("test-contact", props.Get("Name"))
	a.Equal("test@example.com", props.Get("Address"))
	a.Equal("ACTIVE", props.Get("Status"))
}

func Test_Mock_NotificationsContactsEmailContact_String(t *testing.T) {
	a := assert.New(t)

	ec := &NotificationsContactsEmailContact{
		Name: ptr.String("test-contact"),
	}

	a.Equal("test-contact", ec.String())
}

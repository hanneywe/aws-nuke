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

func Test_Mock_MailManagerRelay_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	mockClient.
		On("ListRelays", mock.Anything, mock.Anything).
		Return(
			&mailmanager.ListRelaysOutput{
				Relays: []mailmanagertypes.Relay{
					{
						RelayId:   ptr.String("relay-12345"),
						RelayName: ptr.String("test-relay"),
					},
				},
			}, nil,
		)

	lister := &MailManagerRelayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	relay := resources[0].(*MailManagerRelay)
	assertions.Equal("relay-12345", *relay.RelayID)
	assertions.Equal("test-relay", *relay.RelayName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRelay_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	mockClient.
		On("ListRelays", mock.Anything, mock.Anything).
		Return(
			&mailmanager.ListRelaysOutput{
				Relays: []mailmanagertypes.Relay{},
			}, nil,
		)

	lister := &MailManagerRelayLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRelay_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	relay := &MailManagerRelay{
		svc:       mockClient,
		RelayID:   ptr.String("relay-12345"),
		RelayName: ptr.String("test-relay"),
	}

	mockClient.
		On(
			"DeleteRelay",
			mock.Anything,
			&mailmanager.DeleteRelayInput{
				RelayId: relay.RelayID,
			},
		).
		Return(&mailmanager.DeleteRelayOutput{}, nil)

	err := relay.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRelay_Properties(t *testing.T) {
	assertions := assert.New(t)

	relay := MailManagerRelay{
		RelayID:   ptr.String("relay-12345"),
		RelayName: ptr.String("test-relay"),
	}

	properties := relay.Properties()

	assertions.Equal("relay-12345", properties.Get("RelayId"))
	assertions.Equal("test-relay", properties.Get("RelayName"))
}

func Test_Mock_MailManagerRelay_String(t *testing.T) {
	assertions := assert.New(t)

	relay := MailManagerRelay{
		RelayID:   ptr.String("relay-12345"),
		RelayName: ptr.String("test-relay"),
	}

	assertions.Equal("test-relay", relay.String())
}

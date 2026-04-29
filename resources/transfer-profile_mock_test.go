package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
	transfertypes "github.com/aws/aws-sdk-go-v2/service/transfer/types"
)

func Test_Mock_TransferProfile_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.
		On("ListProfiles", mock.Anything, mock.Anything).
		Return(&transfer.ListProfilesOutput{
			Profiles: []transfertypes.ListedProfile{
				{
					ProfileId: ptr.String("p-1234567890abcdef0"),
					As2Id:     ptr.String("my-as2-id"),
				},
			},
		}, nil)

	lister := &TransferProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	transferProfile := resources[0].(*TransferProfile)
	assertions.Equal("p-1234567890abcdef0", *transferProfile.ProfileID)
	assertions.Equal("my-as2-id", *transferProfile.As2ID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferProfile_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.
		On("ListProfiles", mock.Anything, mock.Anything).
		Return(&transfer.ListProfilesOutput{
			Profiles: []transfertypes.ListedProfile{},
		}, nil)

	lister := &TransferProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferProfile_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockTransferClient)

	transferProfile := &TransferProfile{
		svc:       mockClient,
		ProfileID: ptr.String("p-1234567890abcdef0"),
		As2ID:     ptr.String("my-as2-id"),
	}

	mockClient.
		On("DeleteProfile", mock.Anything, &transfer.DeleteProfileInput{
			ProfileId: transferProfile.ProfileID,
		}).
		Return(&transfer.DeleteProfileOutput{}, nil)

	err := transferProfile.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferProfile_Properties(t *testing.T) {
	assertions := assert.New(t)

	transferProfile := TransferProfile{
		ProfileID: ptr.String("p-1234567890abcdef0"),
		As2ID:     ptr.String("my-as2-id"),
	}

	properties := transferProfile.Properties()

	assertions.Equal("p-1234567890abcdef0", properties.Get("ProfileID"))
	assertions.Equal("my-as2-id", properties.Get("As2ID"))
}

func Test_Mock_TransferProfile_String(t *testing.T) {
	assertions := assert.New(t)

	transferProfile := TransferProfile{
		ProfileID: ptr.String("p-1234567890abcdef0"),
	}

	assertions.Equal("p-1234567890abcdef0", transferProfile.String())
}

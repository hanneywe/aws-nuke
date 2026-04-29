package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/b2bi"
	b2bitypes "github.com/aws/aws-sdk-go-v2/service/b2bi/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testB2BIListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_B2BIProfile_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockB2BIClient)

	mockClient.
		On("ListProfiles", mock.Anything, mock.Anything).
		Return(
			&b2bi.ListProfilesOutput{
				Profiles: []b2bitypes.ProfileSummary{
					{
						ProfileId: ptr.String("p-1234567890"),
						Name:      ptr.String("test-profile"),
					},
				},
			}, nil,
		)

	lister := &B2BIProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testB2BIListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	profile := resources[0].(*B2BIProfile)
	assertions.Equal("p-1234567890", *profile.ProfileID)
	assertions.Equal("test-profile", *profile.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_B2BIProfile_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockB2BIClient)

	mockClient.
		On("ListProfiles", mock.Anything, mock.Anything).
		Return(
			&b2bi.ListProfilesOutput{
				Profiles: []b2bitypes.ProfileSummary{},
			}, nil,
		)

	lister := &B2BIProfileLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testB2BIListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_B2BIProfile_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockB2BIClient)

	profile := &B2BIProfile{
		svc:       mockClient,
		ProfileID: ptr.String("p-1234567890"),
	}

	mockClient.
		On("DeleteProfile", mock.Anything, &b2bi.DeleteProfileInput{
			ProfileId: profile.ProfileID,
		}).
		Return(&b2bi.DeleteProfileOutput{}, nil)

	err := profile.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_B2BIProfile_Properties(t *testing.T) {
	assertions := assert.New(t)

	profile := B2BIProfile{
		ProfileID: ptr.String("p-1234567890"),
		Name:      ptr.String("test-profile"),
	}

	properties := profile.Properties()

	assertions.Equal("p-1234567890", properties.Get("ProfileId"))
	assertions.Equal("test-profile", properties.Get("Name"))
}

func Test_Mock_B2BIProfile_String(t *testing.T) {
	assertions := assert.New(t)

	profile := B2BIProfile{
		ProfileID: ptr.String("p-1234567890"),
	}

	assertions.Equal("p-1234567890", profile.String())
}

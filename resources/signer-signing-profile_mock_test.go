package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/signer"
	signertypes "github.com/aws/aws-sdk-go-v2/service/signer/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testSignerListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_SignerSigningProfile_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSignerClient)

	mockClient.On("ListSigningProfiles", mock.Anything, mock.Anything).
		Return(&signer.ListSigningProfilesOutput{
			Profiles: []signertypes.SigningProfile{
				{
					ProfileName:    ptr.String("my-profile"),
					ProfileVersion: ptr.String("v1"),
					Status:         signertypes.SigningProfileStatusActive,
				},
			},
		}, nil)

	lister := &SignerSigningProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSignerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	profile := resources[0].(*SignerSigningProfile)
	assertions.Equal("my-profile", *profile.ProfileName)
	assertions.Equal("v1", *profile.ProfileVersion)
	assertions.Equal("Active", *profile.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SignerSigningProfile_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSignerClient)

	mockClient.On("ListSigningProfiles", mock.Anything, mock.Anything).
		Return(&signer.ListSigningProfilesOutput{
			Profiles: []signertypes.SigningProfile{},
		}, nil)

	lister := &SignerSigningProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSignerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SignerSigningProfile_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSignerClient)

	profile := &SignerSigningProfile{
		svc:         mockClient,
		ProfileName: ptr.String("my-profile"),
	}

	mockClient.On("CancelSigningProfile", mock.Anything, &signer.CancelSigningProfileInput{
		ProfileName: profile.ProfileName,
	}).Return(&signer.CancelSigningProfileOutput{}, nil)

	err := profile.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SignerSigningProfile_Properties(t *testing.T) {
	assertions := assert.New(t)

	profile := SignerSigningProfile{
		ProfileName:    ptr.String("my-profile"),
		ProfileVersion: ptr.String("v1"),
		Status:         ptr.String("Active"),
	}

	properties := profile.Properties()
	assertions.Equal("my-profile", properties.Get("ProfileName"))
	assertions.Equal("v1", properties.Get("ProfileVersion"))
	assertions.Equal("Active", properties.Get("Status"))
}

func Test_Mock_SignerSigningProfile_String(t *testing.T) {
	assertions := assert.New(t)
	profile := SignerSigningProfile{ProfileName: ptr.String("my-profile")}
	assertions.Equal("my-profile", profile.String())
}

func Test_Mock_SignerSigningProfile_Filter(t *testing.T) {
	assertions := assert.New(t)

	canceledProfile := SignerSigningProfile{Status: ptr.String("Canceled")}
	assertions.Error(canceledProfile.Filter())

	revokedProfile := SignerSigningProfile{Status: ptr.String("Revoked")}
	assertions.Error(revokedProfile.Filter())

	activeProfile := SignerSigningProfile{Status: ptr.String("Active")}
	assertions.NoError(activeProfile.Filter())
}

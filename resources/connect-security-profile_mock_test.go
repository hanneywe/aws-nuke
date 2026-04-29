package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

func Test_Mock_ConnectSecurityProfile_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	// First, list instances
	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{
					Id:            ptr.String("i-12345"),
					InstanceAlias: ptr.String("my-instance"),
				},
			},
		}, nil)

	// Then, list security profiles for that instance
	mockClient.
		On("ListSecurityProfiles", mock.Anything, mock.Anything).
		Return(&connect.ListSecurityProfilesOutput{
			SecurityProfileSummaryList: []connecttypes.SecurityProfileSummary{
				{
					Id:   ptr.String("sp-12345"),
					Name: ptr.String("my-security-profile"),
				},
			},
		}, nil)

	lister := &ConnectSecurityProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	profile := resources[0].(*ConnectSecurityProfile)
	a.Equal("sp-12345", *profile.SecurityProfileID)
	a.Equal("my-security-profile", *profile.Name)
	a.Equal("i-12345", *profile.InstanceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectSecurityProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	// No instances means no security profiles
	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectSecurityProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectSecurityProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	profile := &ConnectSecurityProfile{
		svc:               mockClient,
		InstanceID:        ptr.String("i-12345"),
		SecurityProfileID: ptr.String("sp-12345"),
		Name:              ptr.String("my-security-profile"),
	}

	mockClient.
		On("DeleteSecurityProfile", mock.Anything, &connect.DeleteSecurityProfileInput{
			InstanceId:        profile.InstanceID,
			SecurityProfileId: profile.SecurityProfileID,
		}).
		Return(&connect.DeleteSecurityProfileOutput{}, nil)

	err := profile.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectSecurityProfile_Properties(t *testing.T) {
	a := assert.New(t)

	profile := ConnectSecurityProfile{
		InstanceID:        ptr.String("i-12345"),
		SecurityProfileID: ptr.String("sp-12345"),
		Name:              ptr.String("my-security-profile"),
	}

	props := profile.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("sp-12345", props.Get("SecurityProfileId"))
	a.Equal("my-security-profile", props.Get("Name"))
}

func Test_Mock_ConnectSecurityProfile_String(t *testing.T) {
	a := assert.New(t)

	profile := ConnectSecurityProfile{
		SecurityProfileID: ptr.String("sp-12345"),
	}

	a.Equal("sp-12345", profile.String())
}

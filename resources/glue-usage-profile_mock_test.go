package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueUsageProfile_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListUsageProfiles", mock.Anything, mock.Anything).
		Return(&glue.ListUsageProfilesOutput{
			Profiles: []gluetypes.UsageProfileDefinition{
				{Name: ptr.String("my-profile")},
			},
		}, nil)
	lister := &GlueUsageProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-profile", resources[0].(*GlueUsageProfile).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUsageProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListUsageProfiles", mock.Anything, mock.Anything).
		Return(&glue.ListUsageProfilesOutput{Profiles: []gluetypes.UsageProfileDefinition{}}, nil)
	lister := &GlueUsageProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUsageProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueUsageProfile{svc: mockClient, Name: ptr.String("my-profile")}
	mockClient.On("DeleteUsageProfile", mock.Anything, &glue.DeleteUsageProfileInput{
		Name: r.Name,
	}).Return(&glue.DeleteUsageProfileOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUsageProfile_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueUsageProfile{Name: ptr.String("my-profile")}
	a.Equal("my-profile", r.Properties().Get("Name"))
}

func Test_Mock_GlueUsageProfile_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-profile", (&GlueUsageProfile{Name: ptr.String("my-profile")}).String())
}

package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
	appstreamtypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
)

func Test_Mock_AppStreamUser_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)
	mockClient.On("DescribeUsers", mock.Anything, mock.Anything).
		Return(&appstream.DescribeUsersOutput{
			Users: []appstreamtypes.User{
				{
					UserName:           ptr.String("test@example.com"),
					AuthenticationType: appstreamtypes.AuthenticationTypeUserpool,
					Enabled:            ptr.Bool(true),
				},
			},
		}, nil)
	lister := &AppStreamUserLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppStreamListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	u := resources[0].(*AppStreamUser)
	a.Equal("test@example.com", *u.UserName)
	a.Equal(appstreamtypes.AuthenticationTypeUserpool, u.AuthenticationType)
	a.True(*u.Enabled)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamUser_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)
	mockClient.On("DescribeUsers", mock.Anything, mock.Anything).
		Return(&appstream.DescribeUsersOutput{Users: []appstreamtypes.User{}}, nil)
	lister := &AppStreamUserLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppStreamListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamUser_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)
	u := &AppStreamUser{
		svc:                mockClient,
		UserName:           ptr.String("test@example.com"),
		AuthenticationType: appstreamtypes.AuthenticationTypeUserpool,
	}
	mockClient.On("DeleteUser", mock.Anything, &appstream.DeleteUserInput{
		UserName:           u.UserName,
		AuthenticationType: appstreamtypes.AuthenticationTypeUserpool,
	}).Return(&appstream.DeleteUserOutput{}, nil)
	a.NoError(u.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamUser_Properties(t *testing.T) {
	a := assert.New(t)
	u := AppStreamUser{
		UserName:           ptr.String("test@example.com"),
		AuthenticationType: appstreamtypes.AuthenticationTypeUserpool,
		Enabled:            ptr.Bool(true),
	}
	props := u.Properties()
	a.Equal("test@example.com", props.Get("UserName"))
	a.Equal("true", props.Get("Enabled"))
}

func Test_Mock_AppStreamUser_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("test@example.com", (&AppStreamUser{UserName: ptr.String("test@example.com")}).String())
}

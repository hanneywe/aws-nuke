package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/location"
	locationtypes "github.com/aws/aws-sdk-go-v2/service/location/types"
)

func Test_Mock_LocationServiceKey_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	now := time.Now()
	mockClient.On("ListKeys", mock.Anything, mock.Anything).
		Return(&location.ListKeysOutput{
			Entries: []locationtypes.ListKeysResponseEntry{
				{
					KeyName:    ptr.String("my-key"),
					CreateTime: &now,
				},
			},
		}, nil)

	lister := &LocationServiceKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	k := resources[0].(*LocationServiceKey)
	a.Equal("my-key", *k.KeyName)
	a.Equal(now, *k.CreateTime)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListKeys", mock.Anything, mock.Anything).
		Return(&location.ListKeysOutput{
			Entries: []locationtypes.ListKeysResponseEntry{},
		}, nil)

	lister := &LocationServiceKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	k := &LocationServiceKey{
		svc:     mockClient,
		KeyName: ptr.String("my-key"),
	}

	mockClient.On("DeleteKey", mock.Anything, &location.DeleteKeyInput{
		KeyName: k.KeyName,
	}).Return(&location.DeleteKeyOutput{}, nil)

	a.NoError(k.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceKey_Properties(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	k := LocationServiceKey{
		KeyName:    ptr.String("my-key"),
		CreateTime: &now,
	}
	props := k.Properties()
	a.Equal("my-key", props.Get("KeyName"))
}

func Test_Mock_LocationServiceKey_String(t *testing.T) {
	a := assert.New(t)
	k := LocationServiceKey{KeyName: ptr.String("my-key")}
	a.Equal("my-key", k.String())
}

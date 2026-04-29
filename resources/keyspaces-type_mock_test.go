package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
	keyspacestypes "github.com/aws/aws-sdk-go-v2/service/keyspaces/types"
)

func Test_Mock_KeyspacesType_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	mockClient.On("ListKeyspaces", mock.Anything, mock.Anything).
		Return(&keyspaces.ListKeyspacesOutput{
			Keyspaces: []keyspacestypes.KeyspaceSummary{
				{KeyspaceName: ptr.String("test-keyspace")},
			},
		}, nil)

	mockClient.On("ListTypes", mock.Anything, mock.Anything).
		Return(&keyspaces.ListTypesOutput{
			Types: []string{"test-type"},
		}, nil)

	lister := &KeyspacesTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKeyspacesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*KeyspacesType)
	a.Equal("test-keyspace", *r.KeyspaceName)
	a.Equal("test-type", *r.TypeName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesType_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	mockClient.On("ListKeyspaces", mock.Anything, mock.Anything).
		Return(&keyspaces.ListKeyspacesOutput{
			Keyspaces: []keyspacestypes.KeyspaceSummary{},
		}, nil)

	lister := &KeyspacesTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKeyspacesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesType_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	r := &KeyspacesType{
		svc:          mockClient,
		KeyspaceName: ptr.String("test-keyspace"),
		TypeName:     ptr.String("test-type"),
	}

	mockClient.On("DeleteType", mock.Anything,
		&keyspaces.DeleteTypeInput{
			KeyspaceName: r.KeyspaceName,
			TypeName:     r.TypeName,
		}).Return(&keyspaces.DeleteTypeOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesType_Properties(t *testing.T) {
	a := assert.New(t)
	r := &KeyspacesType{
		KeyspaceName: ptr.String("test-keyspace"),
		TypeName:     ptr.String("test-type"),
	}
	props := r.Properties()
	a.Equal("test-keyspace", props.Get("KeyspaceName"))
	a.Equal("test-type", props.Get("TypeName"))
}

func Test_Mock_KeyspacesType_String(t *testing.T) {
	a := assert.New(t)
	r := &KeyspacesType{
		KeyspaceName: ptr.String("test-keyspace"),
		TypeName:     ptr.String("test-type"),
	}
	a.Equal(fmt.Sprintf("%s/%s", "test-keyspace", "test-type"), r.String())
}

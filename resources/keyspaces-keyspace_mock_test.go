package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
	keyspacestypes "github.com/aws/aws-sdk-go-v2/service/keyspaces/types"
)

func Test_Mock_KeyspacesKeyspace_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	mockClient.On("ListKeyspaces", mock.Anything, mock.Anything).
		Return(&keyspaces.ListKeyspacesOutput{
			Keyspaces: []keyspacestypes.KeyspaceSummary{
				{
					KeyspaceName: ptr.String("my_keyspace"),
					ResourceArn:  ptr.String("arn:aws:cassandra:us-east-1:123456789012:/keyspace/my_keyspace"),
				},
			},
		}, nil)

	lister := &KeyspacesKeyspaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKeyspacesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ks := resources[0].(*KeyspacesKeyspace)
	a.Equal("my_keyspace", *ks.KeyspaceName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesKeyspace_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	mockClient.On("ListKeyspaces", mock.Anything, mock.Anything).
		Return(&keyspaces.ListKeyspacesOutput{
			Keyspaces: []keyspacestypes.KeyspaceSummary{},
		}, nil)

	lister := &KeyspacesKeyspaceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKeyspacesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesKeyspace_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKeyspacesClient)

	ks := &KeyspacesKeyspace{
		svc:          mockClient,
		KeyspaceName: ptr.String("my_keyspace"),
	}

	mockClient.On("DeleteKeyspace", mock.Anything, &keyspaces.DeleteKeyspaceInput{
		KeyspaceName: ks.KeyspaceName,
	}).Return(&keyspaces.DeleteKeyspaceOutput{}, nil)

	a.NoError(ks.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_KeyspacesKeyspace_Filter_System(t *testing.T) {
	a := assert.New(t)

	for _, name := range []string{"system", "system_schema", "system_schema_mcs", "system_multiregion_info"} {
		ks := KeyspacesKeyspace{KeyspaceName: ptr.String(name)}
		a.Error(ks.Filter())
	}
}

func Test_Mock_KeyspacesKeyspace_Filter_User(t *testing.T) {
	a := assert.New(t)
	ks := KeyspacesKeyspace{KeyspaceName: ptr.String("my_keyspace")}
	a.NoError(ks.Filter())
}

func Test_Mock_KeyspacesKeyspace_Properties(t *testing.T) {
	a := assert.New(t)

	ks := KeyspacesKeyspace{
		KeyspaceName: ptr.String("my_keyspace"),
		ResourceArn:  ptr.String("arn:aws:cassandra:us-east-1:123456789012:/keyspace/my_keyspace"),
	}

	props := ks.Properties()
	a.Equal("my_keyspace", props.Get("KeyspaceName"))
}

func Test_Mock_KeyspacesKeyspace_String(t *testing.T) {
	a := assert.New(t)
	ks := KeyspacesKeyspace{KeyspaceName: ptr.String("my_keyspace")}
	a.Equal("my_keyspace", ks.String())
}

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

func Test_Mock_TransferSSHPublicKey_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListServers", mock.Anything, mock.Anything).
		Return(&transfer.ListServersOutput{
			Servers: []transfertypes.ListedServer{
				{ServerId: ptr.String("test-serverid")},
			},
		}, nil)

	mockClient.On("ListUsers", mock.Anything, mock.Anything).
		Return(&transfer.ListUsersOutput{
			Users: []transfertypes.ListedUser{
				{UserName: ptr.String("test-username")},
			},
		}, nil)

	mockClient.On("DescribeUser", mock.Anything, mock.Anything).
		Return(&transfer.DescribeUserOutput{
			User: &transfertypes.DescribedUser{
				Arn: ptr.String("arn:aws:transfer:us-east-1:123456789012:user/s-test/test-username"),
				SshPublicKeys: []transfertypes.SshPublicKey{
					{
						SshPublicKeyId:   ptr.String("test-sshpublickeyid"),
						SshPublicKeyBody: ptr.String("ssh-rsa AAAA..."),
					},
				},
			},
		}, nil)

	lister := &TransferSSHPublicKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*TransferSSHPublicKey)
	a.Equal("test-serverid", *r.ServerID)
	a.Equal("test-username", *r.UserName)
	a.Equal("test-sshpublickeyid", *r.SSHPublicKeyID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferSSHPublicKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListServers", mock.Anything, mock.Anything).
		Return(&transfer.ListServersOutput{
			Servers: []transfertypes.ListedServer{},
		}, nil)

	lister := &TransferSSHPublicKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferSSHPublicKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	r := &TransferSSHPublicKey{
		svc:            mockClient,
		ServerID:       ptr.String("test-serverid"),
		UserName:       ptr.String("test-username"),
		SSHPublicKeyID: ptr.String("test-sshpublickeyid"),
	}

	mockClient.On("DeleteSshPublicKey", mock.Anything,
		&transfer.DeleteSshPublicKeyInput{
			ServerId:       r.ServerID,
			UserName:       r.UserName,
			SshPublicKeyId: r.SSHPublicKeyID,
		}).Return(&transfer.DeleteSshPublicKeyOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferSSHPublicKey_Properties(t *testing.T) {
	a := assert.New(t)
	r := &TransferSSHPublicKey{
		ServerID:       ptr.String("test-serverid"),
		UserName:       ptr.String("test-username"),
		SSHPublicKeyID: ptr.String("test-sshpublickeyid"),
	}
	props := r.Properties()
	a.Equal("test-serverid", props.Get("ServerID"))
	a.Equal("test-username", props.Get("UserName"))
	a.Equal("test-sshpublickeyid", props.Get("SSHPublicKeyID"))
}

func Test_Mock_TransferSSHPublicKey_String(t *testing.T) {
	a := assert.New(t)
	r := &TransferSSHPublicKey{
		SSHPublicKeyID: ptr.String("test-sshpublickeyid"),
	}
	a.Equal("test-sshpublickeyid", r.String())
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const TransferSSHPublicKeyResource = "TransferSSHPublicKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     TransferSSHPublicKeyResource,
		Scope:    nuke.Account,
		Resource: &TransferSSHPublicKey{},
		Lister:   &TransferSSHPublicKeyLister{},
	})
}

type TransferSSHPublicKeyLister struct {
	svc TransferClient
}

func (l *TransferSSHPublicKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = transfer.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	serverPaginator := transfer.NewListServersPaginator(svc, &transfer.ListServersInput{})
	for serverPaginator.HasMorePages() {
		serverResp, err := serverPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, server := range serverResp.Servers {
			userPaginator := transfer.NewListUsersPaginator(svc, &transfer.ListUsersInput{
				ServerId: server.ServerId,
			})
			for userPaginator.HasMorePages() {
				userResp, err := userPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, user := range userResp.Users {
					keyResp, err := svc.DescribeUser(ctx, &transfer.DescribeUserInput{
						ServerId: server.ServerId,
						UserName: user.UserName,
					})
					if err != nil {
						return nil, err
					}
					for _, key := range keyResp.User.SshPublicKeys {
						resources = append(resources, &TransferSSHPublicKey{
							svc:            svc,
							ServerID:       server.ServerId,
							UserName:       user.UserName,
							SSHPublicKeyID: key.SshPublicKeyId,
						})
					}
				}
			}
		}
	}

	return resources, nil
}

type TransferSSHPublicKey struct {
	svc            TransferClient
	ServerID       *string
	UserName       *string
	SSHPublicKeyID *string
}

func (r *TransferSSHPublicKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSshPublicKey(ctx, &transfer.DeleteSshPublicKeyInput{
		ServerId:       r.ServerID,
		UserName:       r.UserName,
		SshPublicKeyId: r.SSHPublicKeyID,
	})
	return err
}

func (r *TransferSSHPublicKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *TransferSSHPublicKey) String() string {
	return *r.SSHPublicKeyID
}

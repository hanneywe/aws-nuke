package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/keyspaces"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KeyspacesTypeResource = "KeyspacesType"

func init() {
	registry.Register(&registry.Registration{
		Name:     KeyspacesTypeResource,
		Scope:    nuke.Account,
		Resource: &KeyspacesType{},
		Lister:   &KeyspacesTypeLister{},
	})
}

type KeyspacesTypeLister struct {
	svc KeyspacesClient
}

func (l *KeyspacesTypeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = keyspaces.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	ksPaginator := keyspaces.NewListKeyspacesPaginator(svc, &keyspaces.ListKeyspacesInput{})
	for ksPaginator.HasMorePages() {
		ksResp, err := ksPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, ks := range ksResp.Keyspaces {
			typeNamePaginator := keyspaces.NewListTypesPaginator(svc, &keyspaces.ListTypesInput{
				KeyspaceName: ks.KeyspaceName,
			})
			for typeNamePaginator.HasMorePages() {
				typeNameResp, err := typeNamePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, typeName := range typeNameResp.Types {
					resources = append(resources, &KeyspacesType{
						svc:          svc,
						KeyspaceName: ks.KeyspaceName,
						TypeName:     aws.String(typeName),
					})
				}
			}
		}
	}
	return resources, nil
}

type KeyspacesType struct {
	svc          KeyspacesClient
	KeyspaceName *string
	TypeName     *string
}

func (r *KeyspacesType) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteType(ctx, &keyspaces.DeleteTypeInput{
		KeyspaceName: r.KeyspaceName,
		TypeName:     r.TypeName,
	})
	return err
}

func (r *KeyspacesType) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KeyspacesType) String() string {
	return fmt.Sprintf("%s/%s", *r.KeyspaceName, *r.TypeName)
}

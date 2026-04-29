package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KeyspacesKeyspaceResource = "KeyspacesKeyspace"

var systemKeyspaces = map[string]struct{}{
	"system":                  {},
	"system_schema":           {},
	"system_schema_mcs":       {},
	"system_multiregion_info": {},
}

func init() {
	registry.Register(&registry.Registration{
		Name:     KeyspacesKeyspaceResource,
		Scope:    nuke.Account,
		Resource: &KeyspacesKeyspace{},
		Lister:   &KeyspacesKeyspaceLister{},
	})
}

type KeyspacesKeyspaceLister struct {
	svc KeyspacesClient
}

func (l *KeyspacesKeyspaceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = keyspaces.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := keyspaces.NewListKeyspacesPaginator(svc, &keyspaces.ListKeyspacesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, ks := range resp.Keyspaces {
			resources = append(resources, &KeyspacesKeyspace{
				svc:          svc,
				KeyspaceName: ks.KeyspaceName,
				ResourceArn:  ks.ResourceArn,
			})
		}
	}

	return resources, nil
}

type KeyspacesKeyspace struct {
	svc          KeyspacesClient
	KeyspaceName *string
	ResourceArn  *string
}

func (r *KeyspacesKeyspace) Filter() error {
	if _, ok := systemKeyspaces[*r.KeyspaceName]; ok {
		return fmt.Errorf("cannot delete system keyspace %s", *r.KeyspaceName)
	}
	return nil
}

func (r *KeyspacesKeyspace) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteKeyspace(ctx, &keyspaces.DeleteKeyspaceInput{
		KeyspaceName: r.KeyspaceName,
	})
	return err
}

func (r *KeyspacesKeyspace) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KeyspacesKeyspace) String() string {
	return *r.KeyspaceName
}

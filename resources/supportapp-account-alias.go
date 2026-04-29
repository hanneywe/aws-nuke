package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/supportapp"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SupportAppAccountAliasResource = "SupportAppAccountAlias"

func init() {
	registry.Register(&registry.Registration{
		Name:     SupportAppAccountAliasResource,
		Scope:    nuke.Account,
		Resource: &SupportAppAccountAlias{},
		Lister:   &SupportAppAccountAliasLister{},
	})
}

type SupportAppAccountAliasLister struct {
	svc SupportAppClient
}

func (l *SupportAppAccountAliasLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = supportapp.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetAccountAlias(ctx, &supportapp.GetAccountAliasInput{})
	if err != nil {
		return nil, err
	}

	var resources []resource.Resource
	if resp.AccountAlias != nil && *resp.AccountAlias != "" {
		resources = append(resources, &SupportAppAccountAlias{
			svc:          svc,
			AccountAlias: resp.AccountAlias,
		})
	}

	return resources, nil
}

type SupportAppAccountAlias struct {
	svc          SupportAppClient
	AccountAlias *string
}

func (r *SupportAppAccountAlias) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAccountAlias(ctx, &supportapp.DeleteAccountAliasInput{})
	return err
}

func (r *SupportAppAccountAlias) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SupportAppAccountAlias) String() string {
	return *r.AccountAlias
}

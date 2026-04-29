package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IAMAccountAliasResource = "IAMAccountAlias"

func init() {
	registry.Register(&registry.Registration{
		Name:     IAMAccountAliasResource,
		Scope:    nuke.Account,
		Resource: &IAMAccountAlias{},
		Lister:   &IAMAccountAliasLister{},
	})
}

type IAMAccountAliasLister struct {
	svc IAMClient
}

func (l *IAMAccountAliasLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = iam.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := iam.NewListAccountAliasesPaginator(svc, &iam.ListAccountAliasesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, alias := range resp.AccountAliases {
			resources = append(resources, &IAMAccountAlias{
				svc:          svc,
				AccountAlias: aws.String(alias),
			})
		}
	}

	return resources, nil
}

type IAMAccountAlias struct {
	svc          IAMClient
	AccountAlias *string
}

func (r *IAMAccountAlias) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAccountAlias(ctx, &iam.DeleteAccountAliasInput{
		AccountAlias: r.AccountAlias,
	})
	return err
}

func (r *IAMAccountAlias) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IAMAccountAlias) String() string {
	return *r.AccountAlias
}

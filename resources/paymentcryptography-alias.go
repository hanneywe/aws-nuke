package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PaymentCryptographyAliasResource = "PaymentCryptographyAlias"

func init() {
	registry.Register(&registry.Registration{
		Name:     PaymentCryptographyAliasResource,
		Scope:    nuke.Account,
		Resource: &PaymentCryptographyAlias{},
		Lister:   &PaymentCryptographyAliasLister{},
	})
}

type PaymentCryptographyAliasLister struct {
	svc PaymentCryptographyClient
}

func (l *PaymentCryptographyAliasLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = paymentcryptography.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &paymentcryptography.ListAliasesInput{}
	for {
		resp, err := svc.ListAliases(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, alias := range resp.Aliases {
			resources = append(resources, &PaymentCryptographyAlias{
				svc:       svc,
				AliasName: alias.AliasName,
				KeyArn:    alias.KeyArn,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type PaymentCryptographyAlias struct {
	svc       PaymentCryptographyClient
	AliasName *string
	KeyArn    *string
}

func (r *PaymentCryptographyAlias) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAlias(ctx, &paymentcryptography.DeleteAliasInput{
		AliasName: r.AliasName,
	})
	return err
}

func (r *PaymentCryptographyAlias) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PaymentCryptographyAlias) String() string {
	return *r.AliasName
}

package resources

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
	accounttypes "github.com/aws/aws-sdk-go-v2/service/account/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AccountAlternateContactResource = "AccountAlternateContact"

func init() {
	registry.Register(&registry.Registration{
		Name:     AccountAlternateContactResource,
		Scope:    nuke.Account,
		Resource: &AccountAlternateContact{},
		Lister:   &AccountAlternateContactLister{},
	})
}

type AccountAlternateContactLister struct {
	svc AccountClient
}

func (l *AccountAlternateContactLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = account.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	contactTypes := []accounttypes.AlternateContactType{
		accounttypes.AlternateContactTypeBilling,
		accounttypes.AlternateContactTypeOperations,
		accounttypes.AlternateContactTypeSecurity,
	}

	for _, ct := range contactTypes {
		resp, err := svc.GetAlternateContact(ctx, &account.GetAlternateContactInput{
			AlternateContactType: ct,
		})
		if err != nil {
			var rnfe *accounttypes.ResourceNotFoundException
			if errors.As(err, &rnfe) {
				continue
			}
			return nil, err
		}

		if resp.AlternateContact != nil {
			resources = append(resources, &AccountAlternateContact{
				svc:                  svc,
				AlternateContactType: aws.String(string(resp.AlternateContact.AlternateContactType)),
				Name:                 resp.AlternateContact.Name,
				EmailAddress:         resp.AlternateContact.EmailAddress,
			})
		}
	}

	return resources, nil
}

type AccountAlternateContact struct {
	svc                  AccountClient
	AlternateContactType *string
	Name                 *string
	EmailAddress         *string
}

func (r *AccountAlternateContact) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAlternateContact(ctx, &account.DeleteAlternateContactInput{
		AlternateContactType: accounttypes.AlternateContactType(*r.AlternateContactType),
	})
	return err
}

func (r *AccountAlternateContact) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AccountAlternateContact) String() string {
	return *r.AlternateContactType
}

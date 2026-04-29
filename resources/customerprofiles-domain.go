package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CustomerProfilesDomainResource = "CustomerProfilesDomain"

func init() {
	registry.Register(&registry.Registration{
		Name:     CustomerProfilesDomainResource,
		Scope:    nuke.Account,
		Resource: &CustomerProfilesDomain{},
		Lister:   &CustomerProfilesDomainLister{},
	})
}

type CustomerProfilesDomainLister struct {
	svc CustomerProfilesClient
}

func (l *CustomerProfilesDomainLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = customerprofiles.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &customerprofiles.ListDomainsInput{}
	for {
		resp, err := svc.ListDomains(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			resources = append(resources, &CustomerProfilesDomain{
				svc:        svc,
				DomainName: item.DomainName,
				Tags:       item.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type CustomerProfilesDomain struct {
	svc        CustomerProfilesClient
	DomainName *string
	Tags       map[string]string
}

func (r *CustomerProfilesDomain) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDomain(ctx, &customerprofiles.DeleteDomainInput{
		DomainName: r.DomainName,
	})
	return err
}

func (r *CustomerProfilesDomain) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CustomerProfilesDomain) String() string {
	return *r.DomainName
}

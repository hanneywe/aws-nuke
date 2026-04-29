package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectCasesDomainResource = "ConnectCasesDomain"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectCasesDomainResource,
		Scope:    nuke.Account,
		Resource: &ConnectCasesDomain{},
		Lister:   &ConnectCasesDomainLister{},
	})
}

type ConnectCasesDomainLister struct {
	svc ConnectCasesClient
}

func (l *ConnectCasesDomainLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connectcases.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := connectcases.NewListDomainsPaginator(svc, &connectcases.ListDomainsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, domain := range resp.Domains {
			resources = append(resources, &ConnectCasesDomain{
				svc:      svc,
				DomainID: domain.DomainId,
				Name:     domain.Name,
			})
		}
	}

	return resources, nil
}

type ConnectCasesDomain struct {
	svc      ConnectCasesClient
	DomainID *string `property:"name=DomainId"`
	Name     *string
}

func (r *ConnectCasesDomain) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDomain(ctx, &connectcases.DeleteDomainInput{
		DomainId: r.DomainID,
	})
	return err
}

func (r *ConnectCasesDomain) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectCasesDomain) String() string {
	return *r.DomainID
}

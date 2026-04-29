package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/swf"
	swftypes "github.com/aws/aws-sdk-go-v2/service/swf/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SWFDomainResource = "SWFDomain"

func init() {
	registry.Register(&registry.Registration{
		Name:     SWFDomainResource,
		Scope:    nuke.Account,
		Resource: &SWFDomain{},
		Lister:   &SWFDomainLister{},
	})
}

type SWFDomainLister struct {
	svc SWFClient
}

func (l *SWFDomainLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = swf.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := swf.NewListDomainsPaginator(svc, &swf.ListDomainsInput{
		RegistrationStatus: swftypes.RegistrationStatusRegistered,
	})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, d := range resp.DomainInfos {
			resources = append(resources, &SWFDomain{
				svc:    svc,
				Name:   d.Name,
				Status: aws.String(string(d.Status)),
			})
		}
	}
	return resources, nil
}

type SWFDomain struct {
	svc    SWFClient
	Name   *string
	Status *string
}

func (r *SWFDomain) Filter() error {
	if r.Status != nil && *r.Status == string(swftypes.RegistrationStatusDeprecated) {
		return fmt.Errorf("already deprecated")
	}
	return nil
}

func (r *SWFDomain) Remove(ctx context.Context) error {
	_, err := r.svc.DeprecateDomain(ctx, &swf.DeprecateDomainInput{
		Name: r.Name,
	})
	return err
}

func (r *SWFDomain) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SWFDomain) String() string {
	return *r.Name
}

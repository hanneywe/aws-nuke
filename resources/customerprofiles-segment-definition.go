package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CustomerProfilesSegmentDefinitionResource = "CustomerProfilesSegmentDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     CustomerProfilesSegmentDefinitionResource,
		Scope:    nuke.Account,
		Resource: &CustomerProfilesSegmentDefinition{},
		Lister:   &CustomerProfilesSegmentDefinitionLister{},
	})
}

type CustomerProfilesSegmentDefinitionLister struct {
	svc CustomerProfilesClient
}

func (l *CustomerProfilesSegmentDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = customerprofiles.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	domainParams := &customerprofiles.ListDomainsInput{}
	for {
		domainResp, err := svc.ListDomains(ctx, domainParams)
		if err != nil {
			return nil, err
		}

		for _, domain := range domainResp.Items {
			segParams := &customerprofiles.ListSegmentDefinitionsInput{
				DomainName: domain.DomainName,
			}
			for {
				segResp, err := svc.ListSegmentDefinitions(ctx, segParams)
				if err != nil {
					return nil, err
				}

				for _, item := range segResp.Items {
					resources = append(resources, &CustomerProfilesSegmentDefinition{
						svc:                   svc,
						DomainName:            domain.DomainName,
						SegmentDefinitionName: item.SegmentDefinitionName,
						DisplayName:           item.DisplayName,
					})
				}

				if segResp.NextToken == nil {
					break
				}
				segParams.NextToken = segResp.NextToken
			}
		}

		if domainResp.NextToken == nil {
			break
		}
		domainParams.NextToken = domainResp.NextToken
	}

	return resources, nil
}

type CustomerProfilesSegmentDefinition struct {
	svc                   CustomerProfilesClient
	DomainName            *string
	SegmentDefinitionName *string
	DisplayName           *string
}

func (r *CustomerProfilesSegmentDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSegmentDefinition(ctx, &customerprofiles.DeleteSegmentDefinitionInput{
		DomainName:            r.DomainName,
		SegmentDefinitionName: r.SegmentDefinitionName,
	})
	return err
}

func (r *CustomerProfilesSegmentDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CustomerProfilesSegmentDefinition) String() string {
	return *r.SegmentDefinitionName
}

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

const SWFActivityTypeResource = "SWFActivityType"

func init() {
	registry.Register(&registry.Registration{
		Name:     SWFActivityTypeResource,
		Scope:    nuke.Account,
		Resource: &SWFActivityType{},
		Lister:   &SWFActivityTypeLister{},
	})
}

type SWFActivityTypeLister struct {
	svc SWFClient
}

func (l *SWFActivityTypeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = swf.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	domainPaginator := swf.NewListDomainsPaginator(svc, &swf.ListDomainsInput{
		RegistrationStatus: swftypes.RegistrationStatusRegistered,
	})
	for domainPaginator.HasMorePages() {
		domainResp, err := domainPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, domain := range domainResp.DomainInfos {
			actPaginator := swf.NewListActivityTypesPaginator(svc, &swf.ListActivityTypesInput{
				Domain:             domain.Name,
				RegistrationStatus: swftypes.RegistrationStatusRegistered,
			})
			for actPaginator.HasMorePages() {
				actResp, err := actPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, info := range actResp.TypeInfos {
					resources = append(resources, &SWFActivityType{
						svc:     svc,
						Domain:  domain.Name,
						Name:    info.ActivityType.Name,
						Version: info.ActivityType.Version,
						Status:  aws.String(string(info.Status)),
					})
				}
			}
		}
	}

	return resources, nil
}

type SWFActivityType struct {
	svc     SWFClient
	Domain  *string
	Name    *string
	Version *string
	Status  *string
}

func (r *SWFActivityType) Filter() error {
	if r.Status != nil && *r.Status == string(swftypes.RegistrationStatusDeprecated) {
		return fmt.Errorf("already deprecated")
	}
	return nil
}

func (r *SWFActivityType) Remove(ctx context.Context) error {
	_, err := r.svc.DeprecateActivityType(ctx, &swf.DeprecateActivityTypeInput{
		Domain: r.Domain,
		ActivityType: &swftypes.ActivityType{
			Name:    r.Name,
			Version: r.Version,
		},
	})
	return err
}

func (r *SWFActivityType) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SWFActivityType) String() string {
	return *r.Name
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectCasesTemplateResource = "ConnectCasesTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectCasesTemplateResource,
		Scope:    nuke.Account,
		Resource: &ConnectCasesTemplate{},
		Lister:   &ConnectCasesTemplateLister{},
	})
}

type ConnectCasesTemplateLister struct {
	svc ConnectCasesClient
}

func (l *ConnectCasesTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connectcases.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	domainPaginator := connectcases.NewListDomainsPaginator(svc, &connectcases.ListDomainsInput{})

	for domainPaginator.HasMorePages() {
		domainResp, err := domainPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, domain := range domainResp.Domains {
			templatePaginator := connectcases.NewListTemplatesPaginator(svc, &connectcases.ListTemplatesInput{
				DomainId: domain.DomainId,
			})

			for templatePaginator.HasMorePages() {
				templateResp, err := templatePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, tmpl := range templateResp.Templates {
					resources = append(resources, &ConnectCasesTemplate{
						svc:        svc,
						DomainID:   domain.DomainId,
						TemplateID: tmpl.TemplateId,
						Name:       tmpl.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectCasesTemplate struct {
	svc        ConnectCasesClient
	DomainID   *string `property:"name=DomainId"`
	TemplateID *string `property:"name=TemplateId"`
	Name       *string
}

func (r *ConnectCasesTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTemplate(ctx, &connectcases.DeleteTemplateInput{
		DomainId:   r.DomainID,
		TemplateId: r.TemplateID,
	})
	return err
}

func (r *ConnectCasesTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectCasesTemplate) String() string {
	return *r.TemplateID
}

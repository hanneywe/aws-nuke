package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"
	connectcasestypes "github.com/aws/aws-sdk-go-v2/service/connectcases/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectCasesFieldResource = "ConnectCasesField"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectCasesFieldResource,
		Scope:    nuke.Account,
		Resource: &ConnectCasesField{},
		Lister:   &ConnectCasesFieldLister{},
	})
}

type ConnectCasesFieldLister struct {
	svc ConnectCasesClient
}

func (l *ConnectCasesFieldLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			fieldPaginator := connectcases.NewListFieldsPaginator(svc, &connectcases.ListFieldsInput{
				DomainId: domain.DomainId,
			})

			for fieldPaginator.HasMorePages() {
				fieldResp, err := fieldPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, field := range fieldResp.Fields {
					resources = append(resources, &ConnectCasesField{
						svc:       svc,
						DomainID:  domain.DomainId,
						FieldID:   field.FieldId,
						Name:      field.Name,
						Namespace: field.Namespace,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectCasesField struct {
	svc       ConnectCasesClient
	DomainID  *string `property:"name=DomainId"`
	FieldID   *string `property:"name=FieldId"`
	Name      *string
	Namespace connectcasestypes.FieldNamespace `property:"name=Namespace"`
}

func (r *ConnectCasesField) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteField(ctx, &connectcases.DeleteFieldInput{
		DomainId: r.DomainID,
		FieldId:  r.FieldID,
	})
	return err
}

func (r *ConnectCasesField) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectCasesField) String() string {
	return *r.FieldID
}

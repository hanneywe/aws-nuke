package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AmplifyDomainAssociationResource = "AmplifyDomainAssociation"

func init() {
	registry.Register(&registry.Registration{
		Name:     AmplifyDomainAssociationResource,
		Scope:    nuke.Account,
		Resource: &AmplifyDomainAssociation{},
		Lister:   &AmplifyDomainAssociationLister{},
	})
}

type AmplifyDomainAssociationLister struct {
	svc AmplifyClient
}

func (l *AmplifyDomainAssociationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = amplify.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	appPaginator := amplify.NewListAppsPaginator(svc, &amplify.ListAppsInput{})
	for appPaginator.HasMorePages() {
		appResp, err := appPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range appResp.Apps {
			domainPaginator := amplify.NewListDomainAssociationsPaginator(svc, &amplify.ListDomainAssociationsInput{
				AppId: appResp.Apps[i].AppId,
			})
			for domainPaginator.HasMorePages() {
				domainResp, err := domainPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for j := range domainResp.DomainAssociations {
					resources = append(resources, &AmplifyDomainAssociation{
						svc:                  svc,
						DomainName:           domainResp.DomainAssociations[j].DomainName,
						DomainAssociationArn: domainResp.DomainAssociations[j].DomainAssociationArn,
						AppID:                appResp.Apps[i].AppId,
					})
				}
			}
		}
	}
	return resources, nil
}

type AmplifyDomainAssociation struct {
	svc                  AmplifyClient
	DomainName           *string
	DomainAssociationArn *string
	AppID                *string `property:"name=AppId"`
}

func (r *AmplifyDomainAssociation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDomainAssociation(ctx, &amplify.DeleteDomainAssociationInput{
		AppId:      r.AppID,
		DomainName: r.DomainName,
	})
	return err
}

func (r *AmplifyDomainAssociation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AmplifyDomainAssociation) String() string {
	return *r.DomainName
}

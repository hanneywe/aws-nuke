package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueCatalogResource = "GlueCatalog"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueCatalogResource,
		Scope:    nuke.Account,
		Resource: &GlueCatalog{},
		Lister:   &GlueCatalogLister{},
	})
}

type GlueCatalogLister struct {
	svc GlueV2Client
}

func (l *GlueCatalogLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &glue.GetCatalogsInput{}
	for {
		resp, err := svc.GetCatalogs(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.CatalogList {
			resources = append(resources, &GlueCatalog{
				svc:       svc,
				CatalogID: resp.CatalogList[i].CatalogId,
				Name:      resp.CatalogList[i].Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GlueCatalog struct {
	svc       GlueV2Client
	CatalogID *string `property:"name=CatalogId"`
	Name      *string
}

func (r *GlueCatalog) Filter() error {
	if r.CatalogID != nil && *r.CatalogID == *r.Name {
		return fmt.Errorf("cannot delete default account catalog")
	}
	return nil
}

func (r *GlueCatalog) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCatalog(ctx, &glue.DeleteCatalogInput{
		CatalogId: r.CatalogID,
	})
	return err
}

func (r *GlueCatalog) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueCatalog) String() string {
	return *r.Name
}

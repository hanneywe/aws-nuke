package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const APIGatewayV2ProductPageResource = "APIGatewayV2ProductPage"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayV2ProductPageResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayV2ProductPage{},
		Lister:   &APIGatewayV2ProductPageLister{},
	})
}

type APIGatewayV2ProductPageLister struct {
	svc APIGatewayV2PortalClient
}

func (l *APIGatewayV2ProductPageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = apigatewayv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	productParams := &apigatewayv2.ListPortalProductsInput{}
	for {
		productResp, err := svc.ListPortalProducts(ctx, productParams)
		if err != nil {
			return nil, err
		}

		for i := range productResp.Items {
			product := &productResp.Items[i]

			pageParams := &apigatewayv2.ListProductPagesInput{
				PortalProductId: product.PortalProductId,
			}
			for {
				pageResp, err := svc.ListProductPages(ctx, pageParams)
				if err != nil {
					return nil, err
				}

				for j := range pageResp.Items {
					page := &pageResp.Items[j]
					resources = append(resources, &APIGatewayV2ProductPage{
						svc:             svc,
						PortalProductID: product.PortalProductId,
						ProductPageID:   page.ProductPageId,
						PageTitle:       page.PageTitle,
					})
				}

				if pageResp.NextToken == nil {
					break
				}
				pageParams.NextToken = pageResp.NextToken
			}
		}

		if productResp.NextToken == nil {
			break
		}
		productParams.NextToken = productResp.NextToken
	}

	return resources, nil
}

type APIGatewayV2ProductPage struct {
	svc             APIGatewayV2PortalClient
	PortalProductID *string
	ProductPageID   *string
	PageTitle       *string
}

func (r *APIGatewayV2ProductPage) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProductPage(ctx, &apigatewayv2.DeleteProductPageInput{
		PortalProductId: r.PortalProductID,
		ProductPageId:   r.ProductPageID,
	})
	return err
}

func (r *APIGatewayV2ProductPage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayV2ProductPage) String() string {
	if r.PageTitle != nil {
		return fmt.Sprintf("%s -> %s (%s)", *r.PortalProductID, *r.ProductPageID, *r.PageTitle)
	}
	return fmt.Sprintf("%s -> %s", *r.PortalProductID, *r.ProductPageID)
}

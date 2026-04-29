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

const APIGatewayV2ProductRestEndpointPageResource = "APIGatewayV2ProductRestEndpointPage"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayV2ProductRestEndpointPageResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayV2ProductRestEndpointPage{},
		Lister:   &APIGatewayV2ProductRestEndpointPageLister{},
	})
}

type APIGatewayV2ProductRestEndpointPageLister struct {
	svc APIGatewayV2PortalClient
}

func (l *APIGatewayV2ProductRestEndpointPageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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

			pageParams := &apigatewayv2.ListProductRestEndpointPagesInput{
				PortalProductId: product.PortalProductId,
			}
			for {
				pageResp, err := svc.ListProductRestEndpointPages(ctx, pageParams)
				if err != nil {
					return nil, err
				}

				for j := range pageResp.Items {
					page := &pageResp.Items[j]
					resources = append(resources, &APIGatewayV2ProductRestEndpointPage{
						svc:                       svc,
						PortalProductID:           product.PortalProductId,
						ProductRestEndpointPageID: page.ProductRestEndpointPageId,
						Endpoint:                  page.Endpoint,
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

type APIGatewayV2ProductRestEndpointPage struct {
	svc                       APIGatewayV2PortalClient
	PortalProductID           *string
	ProductRestEndpointPageID *string
	Endpoint                  *string
}

func (r *APIGatewayV2ProductRestEndpointPage) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteProductRestEndpointPage(ctx, &apigatewayv2.DeleteProductRestEndpointPageInput{
		PortalProductId:           r.PortalProductID,
		ProductRestEndpointPageId: r.ProductRestEndpointPageID,
	})
	return err
}

func (r *APIGatewayV2ProductRestEndpointPage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayV2ProductRestEndpointPage) String() string {
	if r.Endpoint != nil {
		return fmt.Sprintf("%s -> %s (%s)", *r.PortalProductID, *r.ProductRestEndpointPageID, *r.Endpoint)
	}
	return fmt.Sprintf("%s -> %s", *r.PortalProductID, *r.ProductRestEndpointPageID)
}

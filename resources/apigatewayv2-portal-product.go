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

const APIGatewayV2PortalProductResource = "APIGatewayV2PortalProduct"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayV2PortalProductResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayV2PortalProduct{},
		Lister:   &APIGatewayV2PortalProductLister{},
	})
}

type APIGatewayV2PortalProductLister struct {
	svc APIGatewayV2PortalClient
}

func (l *APIGatewayV2PortalProductLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = apigatewayv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &apigatewayv2.ListPortalProductsInput{}

	for {
		resp, err := svc.ListPortalProducts(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Items {
			product := &resp.Items[i]
			resources = append(resources, &APIGatewayV2PortalProduct{
				svc:             svc,
				PortalProductID: product.PortalProductId,
				DisplayName:     product.DisplayName,
				Tags:            product.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type APIGatewayV2PortalProduct struct {
	svc             APIGatewayV2PortalClient
	PortalProductID *string
	DisplayName     *string
	Tags            map[string]string
}

func (r *APIGatewayV2PortalProduct) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePortalProduct(ctx, &apigatewayv2.DeletePortalProductInput{
		PortalProductId: r.PortalProductID,
	})
	return err
}

func (r *APIGatewayV2PortalProduct) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayV2PortalProduct) String() string {
	if r.DisplayName != nil {
		return fmt.Sprintf("%s (%s)", *r.PortalProductID, *r.DisplayName)
	}
	return *r.PortalProductID
}

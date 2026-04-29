package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"

	liberrors "github.com/ekristen/libnuke/pkg/errors"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type APIGatewayV2PortalClient interface {
	ListPortals(ctx context.Context, params *apigatewayv2.ListPortalsInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.ListPortalsOutput, error)
	DeletePortal(ctx context.Context, params *apigatewayv2.DeletePortalInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeletePortalOutput, error)
	DisablePortal(ctx context.Context, params *apigatewayv2.DisablePortalInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DisablePortalOutput, error)
	ListPortalProducts(ctx context.Context, params *apigatewayv2.ListPortalProductsInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.ListPortalProductsOutput, error)
	DeletePortalProduct(ctx context.Context, params *apigatewayv2.DeletePortalProductInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeletePortalProductOutput, error)
	ListProductPages(ctx context.Context, params *apigatewayv2.ListProductPagesInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.ListProductPagesOutput, error)
	DeleteProductPage(ctx context.Context, params *apigatewayv2.DeleteProductPageInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteProductPageOutput, error)
	ListProductRestEndpointPages(ctx context.Context, params *apigatewayv2.ListProductRestEndpointPagesInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.ListProductRestEndpointPagesOutput, error)
	DeleteProductRestEndpointPage(ctx context.Context, params *apigatewayv2.DeleteProductRestEndpointPageInput,
		optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteProductRestEndpointPageOutput, error)
}

const APIGatewayV2PortalResource = "APIGatewayV2Portal"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayV2PortalResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayV2Portal{},
		Lister:   &APIGatewayV2PortalLister{},
		DependsOn: []string{
			APIGatewayV2PortalProductResource,
		},
	})
}

type APIGatewayV2PortalLister struct {
	svc APIGatewayV2PortalClient
}

func (l *APIGatewayV2PortalLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = apigatewayv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &apigatewayv2.ListPortalsInput{}

	for {
		resp, err := svc.ListPortals(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Items {
			portal := &resp.Items[i]

			var name *string
			if portal.PortalContent != nil {
				name = portal.PortalContent.DisplayName
			}

			resources = append(resources, &APIGatewayV2Portal{
				svc:           svc,
				PortalID:      portal.PortalId,
				Name:          name,
				PublishStatus: portal.PublishStatus,
				Tags:          portal.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type APIGatewayV2Portal struct {
	svc           APIGatewayV2PortalClient
	PortalID      *string
	Name          *string
	PublishStatus apigatewayv2types.PublishStatus `property:"name=PublishStatus"`
	Tags          map[string]string
}

func (r *APIGatewayV2Portal) Remove(ctx context.Context) error {
	if r.PublishStatus == apigatewayv2types.PublishStatusPublished {
		_, err := r.svc.DisablePortal(ctx, &apigatewayv2.DisablePortalInput{
			PortalId: r.PortalID,
		})
		if err != nil {
			return err
		}
		// DisablePortal is async — the portal needs time to transition.
		// Return a hold so aws-nuke retries after the portal is fully disabled.
		return liberrors.ErrHoldResource("waiting for portal to be disabled")
	}

	_, err := r.svc.DeletePortal(ctx, &apigatewayv2.DeletePortalInput{
		PortalId: r.PortalID,
	})
	return err
}

func (r *APIGatewayV2Portal) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayV2Portal) String() string {
	if r.Name != nil {
		return fmt.Sprintf("%s (%s)", *r.PortalID, *r.Name)
	}
	return *r.PortalID
}

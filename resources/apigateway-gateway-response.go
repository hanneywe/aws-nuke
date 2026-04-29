package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const APIGatewayGatewayResponseResource = "APIGatewayGatewayResponse"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayGatewayResponseResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayGatewayResponse{},
		Lister:   &APIGatewayGatewayResponseLister{},
	})
}

type APIGatewayGatewayResponseLister struct {
	svc APIGatewayV2Client
}

func (l *APIGatewayGatewayResponseLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = apigateway.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	restAPIParams := &apigateway.GetRestApisInput{
		Limit: ptr.Int32(500),
	}

	for {
		restAPIOutput, err := svc.GetRestApis(ctx, restAPIParams)
		if err != nil {
			return nil, err
		}

		for i := range restAPIOutput.Items {
			restAPI := &restAPIOutput.Items[i]

			gwOutput, err := svc.GetGatewayResponses(ctx, &apigateway.GetGatewayResponsesInput{
				RestApiId: restAPI.Id,
			})
			if err != nil {
				return nil, err
			}

			for j := range gwOutput.Items {
				gw := &gwOutput.Items[j]
				if gw.DefaultResponse {
					continue
				}
				resources = append(resources, &APIGatewayGatewayResponse{
					svc:          svc,
					RestAPIID:    restAPI.Id,
					ResponseType: gw.ResponseType,
					StatusCode:   gw.StatusCode,
				})
			}
		}

		if restAPIOutput.Position == nil {
			break
		}
		restAPIParams.Position = restAPIOutput.Position
	}

	return resources, nil
}

type APIGatewayGatewayResponse struct {
	svc          APIGatewayV2Client
	RestAPIID    *string
	ResponseType apigatewaytypes.GatewayResponseType
	StatusCode   *string
}

func (r *APIGatewayGatewayResponse) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGatewayResponse(ctx, &apigateway.DeleteGatewayResponseInput{
		RestApiId:    r.RestAPIID,
		ResponseType: r.ResponseType,
	})
	return err
}

func (r *APIGatewayGatewayResponse) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayGatewayResponse) String() string {
	return fmt.Sprintf("%s -> %s", *r.RestAPIID, string(r.ResponseType))
}

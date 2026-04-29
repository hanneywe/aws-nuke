package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const APIGatewayV2StageResource = "APIGatewayV2Stage"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayV2StageResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayV2Stage{},
		Lister:   &APIGatewayV2StageLister{},
	})
}

type APIGatewayV2StageLister struct {
	svc Apigatewayv2Client
}

func (l *APIGatewayV2StageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = apigatewayv2.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	apiResp, err := svc.GetApis(ctx, &apigatewayv2.GetApisInput{})
	if err != nil {
		return nil, err
	}
	for iAPI := range apiResp.Items {
		api := &apiResp.Items[iAPI]
		stageResp, err := svc.GetStages(ctx, &apigatewayv2.GetStagesInput{
			ApiId: api.ApiId,
		})
		if err != nil {
			return nil, err
		}
		for iStage := range stageResp.Items {
			stage := &stageResp.Items[iStage]
			resources = append(resources, &APIGatewayV2Stage{
				svc:       svc,
				APIID:     api.ApiId,
				StageName: stage.StageName,
				Tags:      stage.Tags,
			})
		}
	}

	return resources, nil
}

type APIGatewayV2Stage struct {
	svc       Apigatewayv2Client
	APIID     *string `property:"name=ApiId"`
	StageName *string
	Tags      map[string]string
}

func (r *APIGatewayV2Stage) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStage(ctx, &apigatewayv2.DeleteStageInput{
		ApiId:     r.APIID,
		StageName: r.StageName,
	})
	return err
}

func (r *APIGatewayV2Stage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayV2Stage) String() string {
	return *r.StageName
}

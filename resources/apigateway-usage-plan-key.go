package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const APIGatewayUsagePlanKeyResource = "APIGatewayUsagePlanKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     APIGatewayUsagePlanKeyResource,
		Scope:    nuke.Account,
		Resource: &APIGatewayUsagePlanKey{},
		Lister:   &APIGatewayUsagePlanKeyLister{},
	})
}

type APIGatewayUsagePlanKeyLister struct {
	svc ApigatewayClient
}

func (l *APIGatewayUsagePlanKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = apigateway.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	planResp, err := svc.GetUsagePlans(ctx, &apigateway.GetUsagePlansInput{})
	if err != nil {
		return nil, err
	}
	for _, plan := range planResp.Items {
		keyResp, err := svc.GetUsagePlanKeys(ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: plan.Id,
		})
		if err != nil {
			return nil, err
		}
		for _, key := range keyResp.Items {
			resources = append(resources, &APIGatewayUsagePlanKey{
				svc:         svc,
				UsagePlanID: plan.Id,
				KeyID:       key.Id,
				KeyName:     key.Name,
				KeyType:     key.Type,
			})
		}
	}

	return resources, nil
}

type APIGatewayUsagePlanKey struct {
	svc         ApigatewayClient
	UsagePlanID *string `property:"name=UsagePlanId"`
	KeyID       *string `property:"name=KeyId"`
	KeyName     *string
	KeyType     *string
}

func (r *APIGatewayUsagePlanKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUsagePlanKey(ctx, &apigateway.DeleteUsagePlanKeyInput{
		UsagePlanId: r.UsagePlanID,
		KeyId:       r.KeyID,
	})
	return err
}

func (r *APIGatewayUsagePlanKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *APIGatewayUsagePlanKey) String() string {
	return *r.KeyID
}

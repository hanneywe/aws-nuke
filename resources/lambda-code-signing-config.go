package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LambdaCodeSigningConfigResource = "LambdaCodeSigningConfig"

func init() {
	registry.Register(&registry.Registration{
		Name:     LambdaCodeSigningConfigResource,
		Scope:    nuke.Account,
		Resource: &LambdaCodeSigningConfig{},
		Lister:   &LambdaCodeSigningConfigLister{},
	})
}

type LambdaCodeSigningConfigLister struct {
	svc LambdaClient
}

func (l *LambdaCodeSigningConfigLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lambda.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &lambda.ListCodeSigningConfigsInput{}
	for {
		listOutput, err := svc.ListCodeSigningConfigs(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, codeSigningConfig := range listOutput.CodeSigningConfigs {
			resources = append(resources, &LambdaCodeSigningConfig{
				svc:                  svc,
				CodeSigningConfigID:  codeSigningConfig.CodeSigningConfigId,
				CodeSigningConfigArn: codeSigningConfig.CodeSigningConfigArn,
			})
		}

		if listOutput.NextMarker == nil {
			break
		}
		params.Marker = listOutput.NextMarker
	}

	return resources, nil
}

type LambdaCodeSigningConfig struct {
	svc                  LambdaClient
	CodeSigningConfigID  *string
	CodeSigningConfigArn *string
}

func (r *LambdaCodeSigningConfig) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCodeSigningConfig(ctx, &lambda.DeleteCodeSigningConfigInput{
		CodeSigningConfigArn: r.CodeSigningConfigArn,
	})
	return err
}

func (r *LambdaCodeSigningConfig) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LambdaCodeSigningConfig) String() string {
	return *r.CodeSigningConfigID
}

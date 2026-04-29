package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LambdaAliasResource = "LambdaAlias"

func init() {
	registry.Register(&registry.Registration{
		Name:     LambdaAliasResource,
		Scope:    nuke.Account,
		Resource: &LambdaAlias{},
		Lister:   &LambdaAliasLister{},
		DependsOn: []string{
			"LambdaVersion",
		},
	})
}

type LambdaAliasLister struct {
	svc LambdaClient
}

func (l *LambdaAliasLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lambda.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	functionParams := &lambda.ListFunctionsInput{}
	for {
		functionOutput, err := svc.ListFunctions(ctx, functionParams)
		if err != nil {
			return nil, err
		}

		for i := range functionOutput.Functions {
			function := &functionOutput.Functions[i]
			aliasParams := &lambda.ListAliasesInput{
				FunctionName: function.FunctionName,
			}
			for {
				aliasOutput, err := svc.ListAliases(ctx, aliasParams)
				if err != nil {
					return nil, err
				}

				for _, alias := range aliasOutput.Aliases {
					resources = append(resources, &LambdaAlias{
						svc:          svc,
						FunctionName: function.FunctionName,
						Name:         alias.Name,
						AliasArn:     alias.AliasArn,
					})
				}

				if aliasOutput.NextMarker == nil {
					break
				}
				aliasParams.Marker = aliasOutput.NextMarker
			}
		}

		if functionOutput.NextMarker == nil {
			break
		}
		functionParams.Marker = functionOutput.NextMarker
	}

	return resources, nil
}

type LambdaAlias struct {
	svc          LambdaClient
	FunctionName *string
	Name         *string
	AliasArn     *string
}

func (r *LambdaAlias) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAlias(ctx, &lambda.DeleteAliasInput{
		FunctionName: r.FunctionName,
		Name:         r.Name,
	})
	return err
}

func (r *LambdaAlias) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LambdaAlias) String() string {
	return *r.Name
}

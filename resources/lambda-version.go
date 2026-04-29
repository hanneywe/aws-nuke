package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LambdaVersionResource = "LambdaVersion"

const LambdaVersionLatest = "$LATEST"

func init() {
	registry.Register(&registry.Registration{
		Name:     LambdaVersionResource,
		Scope:    nuke.Account,
		Resource: &LambdaVersion{},
		Lister:   &LambdaVersionLister{},
	})
}

type LambdaVersionLister struct {
	svc LambdaClient
}

func (l *LambdaVersionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			versionParams := &lambda.ListVersionsByFunctionInput{
				FunctionName: function.FunctionName,
			}
			for {
				versionOutput, err := svc.ListVersionsByFunction(ctx, versionParams)
				if err != nil {
					return nil, err
				}

				for i := range versionOutput.Versions {
					version := &versionOutput.Versions[i]
					resources = append(resources, &LambdaVersion{
						svc:          svc,
						FunctionName: function.FunctionName,
						Version:      version.Version,
						FunctionArn:  version.FunctionArn,
					})
				}

				if versionOutput.NextMarker == nil {
					break
				}
				versionParams.Marker = versionOutput.NextMarker
			}
		}

		if functionOutput.NextMarker == nil {
			break
		}
		functionParams.Marker = functionOutput.NextMarker
	}

	return resources, nil
}

type LambdaVersion struct {
	svc          LambdaClient
	FunctionName *string
	Version      *string
	FunctionArn  *string
}

func (r *LambdaVersion) Filter() error {
	if r.Version != nil && *r.Version == LambdaVersionLatest {
		return fmt.Errorf("cannot delete %s version", LambdaVersionLatest)
	}
	return nil
}

func (r *LambdaVersion) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: r.FunctionName,
		Qualifier:    r.Version,
	})
	return err
}

func (r *LambdaVersion) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LambdaVersion) String() string {
	return *r.Version
}

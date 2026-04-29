package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueUserDefinedFunctionResource = "GlueUserDefinedFunction"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueUserDefinedFunctionResource,
		Scope:    nuke.Account,
		Resource: &GlueUserDefinedFunction{},
		Lister:   &GlueUserDefinedFunctionLister{},
	})
}

type GlueUserDefinedFunctionLister struct {
	svc GlueV2Client
}

func (l *GlueUserDefinedFunctionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = glue.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	dbParams := &glue.GetDatabasesInput{}
	for {
		dbResp, err := svc.GetDatabases(ctx, dbParams)
		if err != nil {
			return nil, err
		}
		for _, db := range dbResp.DatabaseList {
			fnParams := &glue.GetUserDefinedFunctionsInput{
				DatabaseName: db.Name,
				Pattern:      aws.String("*"),
			}
			for {
				fnResp, err := svc.GetUserDefinedFunctions(ctx, fnParams)
				if err != nil {
					return nil, err
				}
				for _, fn := range fnResp.UserDefinedFunctions {
					resources = append(resources, &GlueUserDefinedFunction{
						svc:          svc,
						DatabaseName: fn.DatabaseName,
						FunctionName: fn.FunctionName,
						CatalogID:    fn.CatalogId,
					})
				}
				if fnResp.NextToken == nil {
					break
				}
				fnParams.NextToken = fnResp.NextToken
			}
		}
		if dbResp.NextToken == nil {
			break
		}
		dbParams.NextToken = dbResp.NextToken
	}
	return resources, nil
}

type GlueUserDefinedFunction struct {
	svc          GlueV2Client
	DatabaseName *string
	FunctionName *string
	CatalogID    *string `property:"name=CatalogId"`
}

func (r *GlueUserDefinedFunction) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUserDefinedFunction(ctx, &glue.DeleteUserDefinedFunctionInput{
		DatabaseName: r.DatabaseName,
		FunctionName: r.FunctionName,
	})
	return err
}

func (r *GlueUserDefinedFunction) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueUserDefinedFunction) String() string {
	return fmt.Sprintf("%s/%s", *r.DatabaseName, *r.FunctionName)
}

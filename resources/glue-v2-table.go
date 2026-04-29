package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GlueTableResource = "GlueTable"

func init() {
	registry.Register(&registry.Registration{
		Name:     GlueTableResource,
		Scope:    nuke.Account,
		Resource: &GlueTable{},
		Lister:   &GlueTableLister{},
	})
}

type GlueTableLister struct {
	svc GlueV2Client
}

func (l *GlueTableLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			tableParams := &glue.GetTablesInput{
				DatabaseName: db.Name,
			}
			for {
				tableResp, err := svc.GetTables(ctx, tableParams)
				if err != nil {
					return nil, err
				}
				for i := range tableResp.TableList {
					resources = append(resources, &GlueTable{
						svc:          svc,
						DatabaseName: db.Name,
						Name:         tableResp.TableList[i].Name,
					})
				}
				if tableResp.NextToken == nil {
					break
				}
				tableParams.NextToken = tableResp.NextToken
			}
		}
		if dbResp.NextToken == nil {
			break
		}
		dbParams.NextToken = dbResp.NextToken
	}

	return resources, nil
}

type GlueTable struct {
	svc          GlueV2Client
	DatabaseName *string
	Name         *string
}

func (r *GlueTable) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTable(ctx, &glue.DeleteTableInput{
		DatabaseName: r.DatabaseName,
		Name:         r.Name,
	})
	return err
}

func (r *GlueTable) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GlueTable) String() string {
	return fmt.Sprintf("%s/%s", *r.DatabaseName, *r.Name)
}

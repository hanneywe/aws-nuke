package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DocDBGlobalClusterResource = "DocDBGlobalCluster"

func init() {
	registry.Register(&registry.Registration{
		Name:     DocDBGlobalClusterResource,
		Scope:    nuke.Account,
		Resource: &DocDBGlobalCluster{},
		Lister:   &DocDBGlobalClusterLister{},
		Settings: []string{
			"DisableDeletionProtection",
		},
	})
}

type DocDBGlobalClusterLister struct {
	svc DocDBV2Client
}

func (l *DocDBGlobalClusterLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = docdb.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := docdb.NewDescribeGlobalClustersPaginator(svc, &docdb.DescribeGlobalClustersInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, globalCluster := range output.GlobalClusters {
			resources = append(resources, &DocDBGlobalCluster{
				svc:                     svc,
				GlobalClusterIdentifier: globalCluster.GlobalClusterIdentifier,
				GlobalClusterArn:        globalCluster.GlobalClusterArn,
				Status:                  globalCluster.Status,
				DeletionProtection:      globalCluster.DeletionProtection,
			})
		}
	}

	return resources, nil
}

type DocDBGlobalCluster struct {
	svc      DocDBV2Client
	settings *libsettings.Setting

	GlobalClusterIdentifier *string
	GlobalClusterArn        *string
	Status                  *string
	DeletionProtection      *bool
}

func (r *DocDBGlobalCluster) Settings(setting *libsettings.Setting) {
	r.settings = setting
}

func (r *DocDBGlobalCluster) Remove(ctx context.Context) error {
	if ptr.ToBool(r.DeletionProtection) && r.settings.GetBool("DisableDeletionProtection") {
		_, err := r.svc.ModifyGlobalCluster(ctx, &docdb.ModifyGlobalClusterInput{
			GlobalClusterIdentifier: r.GlobalClusterIdentifier,
			DeletionProtection:      aws.Bool(false),
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteGlobalCluster(ctx, &docdb.DeleteGlobalClusterInput{
		GlobalClusterIdentifier: r.GlobalClusterIdentifier,
	})
	return err
}

func (r *DocDBGlobalCluster) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DocDBGlobalCluster) String() string {
	return *r.GlobalClusterIdentifier
}

func (r *DocDBGlobalCluster) Filter() error {
	if r.Status != nil && *r.Status == "deleting" { //nolint:goconst
		return fmt.Errorf("already deleting")
	}
	return nil
}

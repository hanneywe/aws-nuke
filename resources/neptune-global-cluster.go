package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NeptuneGlobalClusterResource = "NeptuneGlobalCluster"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneGlobalClusterResource,
		Scope:    nuke.Account,
		Resource: &NeptuneGlobalCluster{},
		Lister:   &NeptuneGlobalClusterLister{},
		Settings: []string{
			"DisableDeletionProtection",
		},
	})
}

type NeptuneGlobalClusterLister struct {
	svc NeptuneV2Client
}

func (l *NeptuneGlobalClusterLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptune.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := neptune.NewDescribeGlobalClustersPaginator(svc, &neptune.DescribeGlobalClustersInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range output.GlobalClusters {
			globalCluster := &output.GlobalClusters[i]
			resources = append(resources, &NeptuneGlobalCluster{
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

type NeptuneGlobalCluster struct {
	svc      NeptuneV2Client
	settings *libsettings.Setting

	GlobalClusterIdentifier *string
	GlobalClusterArn        *string
	Status                  *string
	DeletionProtection      *bool
}

func (r *NeptuneGlobalCluster) Settings(setting *libsettings.Setting) {
	r.settings = setting
}

func (r *NeptuneGlobalCluster) Remove(ctx context.Context) error {
	if ptr.ToBool(r.DeletionProtection) && r.settings.GetBool("DisableDeletionProtection") {
		_, err := r.svc.ModifyGlobalCluster(ctx, &neptune.ModifyGlobalClusterInput{
			GlobalClusterIdentifier: r.GlobalClusterIdentifier,
			DeletionProtection:      aws.Bool(false),
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteGlobalCluster(ctx, &neptune.DeleteGlobalClusterInput{
		GlobalClusterIdentifier: r.GlobalClusterIdentifier,
	})
	return err
}

func (r *NeptuneGlobalCluster) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneGlobalCluster) String() string {
	return *r.GlobalClusterIdentifier
}

func (r *NeptuneGlobalCluster) Filter() error {
	if r.Status != nil && *r.Status == "deleting" {
		return fmt.Errorf("already deleting")
	}
	return nil
}

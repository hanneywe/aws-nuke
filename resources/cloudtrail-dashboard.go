package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudTrailDashboardResource = "CloudTrailDashboard"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudTrailDashboardResource,
		Scope:    nuke.Account,
		Resource: &CloudTrailDashboard{},
		Lister:   &CloudTrailDashboardLister{},
		Settings: []string{
			"DisableTerminationProtection",
		},
	})
}

type CloudTrailDashboardLister struct {
	svc CloudTrailClient
}

func (l *CloudTrailDashboardLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cloudtrail.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &cloudtrail.ListDashboardsInput{}
	for {
		resp, err := svc.ListDashboards(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, d := range resp.Dashboards {
			dashboard := &CloudTrailDashboard{
				svc:          svc,
				DashboardArn: d.DashboardArn,
			}

			// Get dashboard details to check termination protection
			detail, err := svc.GetDashboard(ctx, &cloudtrail.GetDashboardInput{
				DashboardId: d.DashboardArn,
			})
			if err == nil {
				dashboard.TerminationProtectionEnabled = detail.TerminationProtectionEnabled
			}

			resources = append(resources, dashboard)
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type CloudTrailDashboard struct {
	svc                          CloudTrailClient
	settings                     *libsettings.Setting
	DashboardArn                 *string
	TerminationProtectionEnabled *bool
}

func (r *CloudTrailDashboard) Settings(setting *libsettings.Setting) {
	r.settings = setting
}

func (r *CloudTrailDashboard) Filter() error {
	if r.DashboardArn != nil && strings.Contains(*r.DashboardArn, "/AWSCloudTrail-") {
		return fmt.Errorf("cannot delete AWS-managed dashboard")
	}
	return nil
}

func (r *CloudTrailDashboard) Remove(ctx context.Context) error {
	if r.TerminationProtectionEnabled != nil && *r.TerminationProtectionEnabled &&
		r.settings.GetBool("DisableTerminationProtection") {
		_, err := r.svc.UpdateDashboard(ctx, &cloudtrail.UpdateDashboardInput{
			DashboardId:                  r.DashboardArn,
			TerminationProtectionEnabled: aws.Bool(false),
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteDashboard(ctx, &cloudtrail.DeleteDashboardInput{
		DashboardId: r.DashboardArn,
	})
	return err
}

func (r *CloudTrailDashboard) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudTrailDashboard) String() string {
	return *r.DashboardArn
}

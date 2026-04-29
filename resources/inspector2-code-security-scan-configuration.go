package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Inspector2CodeSecurityScanConfigurationResource = "Inspector2CodeSecurityScanConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     Inspector2CodeSecurityScanConfigurationResource,
		Scope:    nuke.Account,
		Resource: &Inspector2CodeSecurityScanConfiguration{},
		Lister:   &Inspector2CodeSecurityScanConfigurationLister{},
	})
}

type Inspector2CodeSecurityScanConfigurationLister struct {
	svc Inspector2V2Client
}

func (l *Inspector2CodeSecurityScanConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = inspector2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &inspector2.ListCodeSecurityScanConfigurationsInput{}
	for {
		resp, err := svc.ListCodeSecurityScanConfigurations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, cfg := range resp.Configurations {
			resources = append(resources, &Inspector2CodeSecurityScanConfiguration{
				svc:                  svc,
				ScanConfigurationArn: cfg.ScanConfigurationArn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type Inspector2CodeSecurityScanConfiguration struct {
	svc                  Inspector2V2Client
	ScanConfigurationArn *string
}

func (r *Inspector2CodeSecurityScanConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCodeSecurityScanConfiguration(ctx, &inspector2.DeleteCodeSecurityScanConfigurationInput{
		ScanConfigurationArn: r.ScanConfigurationArn,
	})
	return err
}

func (r *Inspector2CodeSecurityScanConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Inspector2CodeSecurityScanConfiguration) String() string {
	return *r.ScanConfigurationArn
}

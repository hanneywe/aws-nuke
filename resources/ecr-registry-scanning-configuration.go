package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRRegistryScanningConfigurationResource = "ECRRegistryScanningConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRRegistryScanningConfigurationResource,
		Scope:    nuke.Account,
		Resource: &ECRRegistryScanningConfiguration{},
		Lister:   &ECRRegistryScanningConfigurationLister{},
	})
}

type ECRRegistryScanningConfigurationLister struct {
	svc ECRv2Client
}

func (l *ECRRegistryScanningConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	resp, err := svc.GetRegistryScanningConfiguration(ctx, &ecr.GetRegistryScanningConfigurationInput{})
	if err != nil {
		return nil, err
	}

	var resources []resource.Resource
	if resp.ScanningConfiguration != nil {
		scanType := string(resp.ScanningConfiguration.ScanType)
		resources = append(resources, &ECRRegistryScanningConfiguration{
			svc:      svc,
			ScanType: &scanType,
			Rules:    resp.ScanningConfiguration.Rules,
		})
	}
	return resources, nil
}

type ECRRegistryScanningConfiguration struct {
	svc      ECRv2Client
	ScanType *string
	Rules    []ecrtypes.RegistryScanningRule `property:"-"`
}

func (r *ECRRegistryScanningConfiguration) Filter() error {
	if r.ScanType != nil && *r.ScanType == string(ecrtypes.ScanTypeBasic) && len(r.Rules) == 0 {
		return fmt.Errorf("already at default configuration")
	}
	return nil
}

func (r *ECRRegistryScanningConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.PutRegistryScanningConfiguration(ctx, &ecr.PutRegistryScanningConfigurationInput{
		ScanType: ecrtypes.ScanTypeBasic,
		Rules:    []ecrtypes.RegistryScanningRule{},
	})
	return err
}

func (r *ECRRegistryScanningConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRRegistryScanningConfiguration) String() string {
	return "Registry Scanning Configuration"
}

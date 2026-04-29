package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apprunner"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppRunnerObservabilityConfigurationResource = "AppRunnerObservabilityConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppRunnerObservabilityConfigurationResource,
		Scope:    nuke.Account,
		Resource: &AppRunnerObservabilityConfiguration{},
		Lister:   &AppRunnerObservabilityConfigurationLister{},
	})
}

type AppRunnerObservabilityConfigurationLister struct {
	svc AppRunnerClient
}

func (l *AppRunnerObservabilityConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = apprunner.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &apprunner.ListObservabilityConfigurationsInput{}
	for {
		resp, err := svc.ListObservabilityConfigurations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, oc := range resp.ObservabilityConfigurationSummaryList {
			resources = append(resources, &AppRunnerObservabilityConfiguration{
				svc:                            svc,
				ObservabilityConfigurationArn:  oc.ObservabilityConfigurationArn,
				ObservabilityConfigurationName: oc.ObservabilityConfigurationName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type AppRunnerObservabilityConfiguration struct {
	svc                            AppRunnerClient
	ObservabilityConfigurationArn  *string
	ObservabilityConfigurationName *string
}

func (r *AppRunnerObservabilityConfiguration) Filter() error {
	if r.ObservabilityConfigurationName != nil && *r.ObservabilityConfigurationName == "DefaultConfiguration" {
		return fmt.Errorf("cannot delete default configuration")
	}
	return nil
}

func (r *AppRunnerObservabilityConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteObservabilityConfiguration(ctx, &apprunner.DeleteObservabilityConfigurationInput{
		ObservabilityConfigurationArn: r.ObservabilityConfigurationArn,
	})
	return err
}

func (r *AppRunnerObservabilityConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppRunnerObservabilityConfiguration) String() string {
	return *r.ObservabilityConfigurationName
}

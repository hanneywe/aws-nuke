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

const AppRunnerAutoScalingConfigurationResource = "AppRunnerAutoScalingConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppRunnerAutoScalingConfigurationResource,
		Scope:    nuke.Account,
		Resource: &AppRunnerAutoScalingConfiguration{},
		Lister:   &AppRunnerAutoScalingConfigurationLister{},
	})
}

type AppRunnerAutoScalingConfigurationLister struct {
	svc AppRunnerClient
}

func (l *AppRunnerAutoScalingConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = apprunner.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &apprunner.ListAutoScalingConfigurationsInput{}
	for {
		resp, err := svc.ListAutoScalingConfigurations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, asc := range resp.AutoScalingConfigurationSummaryList {
			resources = append(resources, &AppRunnerAutoScalingConfiguration{
				svc:                          svc,
				AutoScalingConfigurationArn:  asc.AutoScalingConfigurationArn,
				AutoScalingConfigurationName: asc.AutoScalingConfigurationName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type AppRunnerAutoScalingConfiguration struct {
	svc                          AppRunnerClient
	AutoScalingConfigurationArn  *string
	AutoScalingConfigurationName *string
}

func (r *AppRunnerAutoScalingConfiguration) Filter() error {
	if r.AutoScalingConfigurationName != nil && *r.AutoScalingConfigurationName == "DefaultConfiguration" {
		return fmt.Errorf("cannot delete default configuration")
	}
	return nil
}

func (r *AppRunnerAutoScalingConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAutoScalingConfiguration(ctx, &apprunner.DeleteAutoScalingConfigurationInput{
		AutoScalingConfigurationArn: r.AutoScalingConfigurationArn,
	})
	return err
}

func (r *AppRunnerAutoScalingConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppRunnerAutoScalingConfiguration) String() string {
	return *r.AutoScalingConfigurationName
}

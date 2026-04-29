package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ApplicationSignalsGroupingConfigurationResource = "ApplicationSignalsGroupingConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     ApplicationSignalsGroupingConfigurationResource,
		Scope:    nuke.Account,
		Resource: &ApplicationSignalsGroupingConfiguration{},
		Lister:   &ApplicationSignalsGroupingConfigurationLister{},
	})
}

type ApplicationSignalsGroupingConfigurationLister struct {
	svc ApplicationSignalsClient
}

func (l *ApplicationSignalsGroupingConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = applicationsignals.NewFromConfig(*opts.Config)
	}

	resp, err := svc.ListGroupingAttributeDefinitions(ctx, &applicationsignals.ListGroupingAttributeDefinitionsInput{})
	if err != nil {
		return nil, err
	}

	if len(resp.GroupingAttributeDefinitions) == 0 {
		return nil, nil
	}

	var groupingNames []string
	for _, def := range resp.GroupingAttributeDefinitions {
		if def.GroupingName != nil {
			groupingNames = append(groupingNames, *def.GroupingName)
		}
	}

	return []resource.Resource{
		&ApplicationSignalsGroupingConfiguration{
			svc:          svc,
			GroupingType: fmt.Sprintf("custom (%d definitions)", len(resp.GroupingAttributeDefinitions)),
			Definitions:  strings.Join(groupingNames, ", "),
		},
	}, nil
}

type ApplicationSignalsGroupingConfiguration struct {
	svc          ApplicationSignalsClient
	GroupingType string `property:"name=GroupingType"`
	Definitions  string `property:"name=Definitions"`
}

func (r *ApplicationSignalsGroupingConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGroupingConfiguration(ctx, &applicationsignals.DeleteGroupingConfigurationInput{})
	return err
}

func (r *ApplicationSignalsGroupingConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ApplicationSignalsGroupingConfiguration) String() string {
	return "ApplicationSignalsGroupingConfiguration"
}

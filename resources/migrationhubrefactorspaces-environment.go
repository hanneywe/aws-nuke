package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MigrationHubRefactorSpacesEnvironmentResource = "MigrationHubRefactorSpacesEnvironment"

func init() {
	registry.Register(&registry.Registration{
		Name:     MigrationHubRefactorSpacesEnvironmentResource,
		Scope:    nuke.Account,
		Resource: &MigrationHubRefactorSpacesEnvironment{},
		Lister:   &MigrationHubRefactorSpacesEnvironmentLister{},
	})
}

type MigrationHubRefactorSpacesEnvironmentLister struct {
	svc MigrationHubRefactorSpacesClient
}

func (l *MigrationHubRefactorSpacesEnvironmentLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = migrationhubrefactorspaces.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := migrationhubrefactorspaces.NewListEnvironmentsPaginator(svc, &migrationhubrefactorspaces.ListEnvironmentsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, env := range resp.EnvironmentSummaryList {
			resources = append(resources, &MigrationHubRefactorSpacesEnvironment{
				svc:           svc,
				EnvironmentID: env.EnvironmentId,
				Name:          env.Name,
			})
		}
	}
	return resources, nil
}

type MigrationHubRefactorSpacesEnvironment struct {
	svc           MigrationHubRefactorSpacesClient
	EnvironmentID *string `property:"name=EnvironmentId"`
	Name          *string
}

func (r *MigrationHubRefactorSpacesEnvironment) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEnvironment(ctx, &migrationhubrefactorspaces.DeleteEnvironmentInput{
		EnvironmentIdentifier: r.EnvironmentID,
	})
	return err
}

func (r *MigrationHubRefactorSpacesEnvironment) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MigrationHubRefactorSpacesEnvironment) String() string {
	return *r.Name
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsRunGroupResource = "OmicsRunGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsRunGroupResource,
		Scope:    nuke.Account,
		Resource: &OmicsRunGroup{},
		Lister:   &OmicsRunGroupLister{},
	})
}

type OmicsRunGroupLister struct {
	svc OmicsClient
}

func (l *OmicsRunGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := omics.NewListRunGroupsPaginator(svc, &omics.ListRunGroupsInput{})
	for paginator.HasMorePages() {
		runGroupsOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, runGroup := range runGroupsOutput.Items {
			resources = append(resources, &OmicsRunGroup{
				svc:  svc,
				ID:   runGroup.Id,
				Name: runGroup.Name,
			})
		}
	}

	return resources, nil
}

type OmicsRunGroup struct {
	svc  OmicsClient
	ID   *string
	Name *string
}

func (r *OmicsRunGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRunGroup(ctx, &omics.DeleteRunGroupInput{
		Id: r.ID,
	})
	return err
}

func (r *OmicsRunGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsRunGroup) String() string {
	return *r.ID
}

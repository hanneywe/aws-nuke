package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/opensearch"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OpenSearchApplicationResource = "OpenSearchApplication"

func init() {
	registry.Register(&registry.Registration{
		Name:     OpenSearchApplicationResource,
		Scope:    nuke.Account,
		Resource: &OpenSearchApplication{},
		Lister:   &OpenSearchApplicationLister{},
	})
}

type OpenSearchApplicationLister struct {
	svc OpenSearchClient
}

func (l *OpenSearchApplicationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = opensearch.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := opensearch.NewListApplicationsPaginator(svc, &opensearch.ListApplicationsInput{
		MaxResults: 100,
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ApplicationSummaries {
			resources = append(resources, &OpenSearchApplication{
				svc:  svc,
				ID:   item.Id,
				Name: item.Name,
			})
		}
	}

	return resources, nil
}

type OpenSearchApplication struct {
	svc  OpenSearchClient
	ID   *string
	Name *string
}

func (r *OpenSearchApplication) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteApplication(ctx, &opensearch.DeleteApplicationInput{
		Id: r.ID,
	})
	return err
}

func (r *OpenSearchApplication) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OpenSearchApplication) String() string {
	return *r.ID
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsSequenceStoreResource = "OmicsSequenceStore"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsSequenceStoreResource,
		Scope:    nuke.Account,
		Resource: &OmicsSequenceStore{},
		Lister:   &OmicsSequenceStoreLister{},
		DependsOn: []string{
			OmicsReadSetResource,
		},
	})
}

type OmicsSequenceStoreLister struct {
	svc OmicsClient
}

func (l *OmicsSequenceStoreLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := omics.NewListSequenceStoresPaginator(svc, &omics.ListSequenceStoresInput{})
	for paginator.HasMorePages() {
		sequenceStoresOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sequenceStore := range sequenceStoresOutput.SequenceStores {
			resources = append(resources, &OmicsSequenceStore{
				svc:  svc,
				ID:   sequenceStore.Id,
				Name: sequenceStore.Name,
			})
		}
	}

	return resources, nil
}

type OmicsSequenceStore struct {
	svc  OmicsClient
	ID   *string
	Name *string
}

func (r *OmicsSequenceStore) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSequenceStore(ctx, &omics.DeleteSequenceStoreInput{
		Id: r.ID,
	})
	return err
}

func (r *OmicsSequenceStore) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsSequenceStore) String() string {
	return *r.ID
}

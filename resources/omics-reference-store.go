package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsReferenceStoreResource = "OmicsReferenceStore"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsReferenceStoreResource,
		Scope:    nuke.Account,
		Resource: &OmicsReferenceStore{},
		Lister:   &OmicsReferenceStoreLister{},
		DependsOn: []string{
			OmicsReferenceResource,
		},
	})
}

type OmicsReferenceStoreLister struct {
	svc OmicsClient
}

func (l *OmicsReferenceStoreLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := omics.NewListReferenceStoresPaginator(svc, &omics.ListReferenceStoresInput{})
	for paginator.HasMorePages() {
		referenceStoresOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, referenceStore := range referenceStoresOutput.ReferenceStores {
			resources = append(resources, &OmicsReferenceStore{
				svc:  svc,
				ID:   referenceStore.Id,
				Name: referenceStore.Name,
			})
		}
	}

	return resources, nil
}

type OmicsReferenceStore struct {
	svc  OmicsClient
	ID   *string
	Name *string
}

func (r *OmicsReferenceStore) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteReferenceStore(ctx, &omics.DeleteReferenceStoreInput{
		Id: r.ID,
	})
	return err
}

func (r *OmicsReferenceStore) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsReferenceStore) String() string {
	return *r.ID
}

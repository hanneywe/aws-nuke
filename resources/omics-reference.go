package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsReferenceResource = "OmicsReference"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsReferenceResource,
		Scope:    nuke.Account,
		Resource: &OmicsReference{},
		Lister:   &OmicsReferenceLister{},
		DependsOn: []string{
			OmicsReadSetResource,
		},
	})
}

type OmicsReferenceLister struct {
	svc OmicsClient
}

func (l *OmicsReferenceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, list all reference stores
	referenceStorePaginator := omics.NewListReferenceStoresPaginator(svc, &omics.ListReferenceStoresInput{})
	for referenceStorePaginator.HasMorePages() {
		referenceStoresOutput, err := referenceStorePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, referenceStore := range referenceStoresOutput.ReferenceStores {
			// For each reference store, list all references
			referencePaginator := omics.NewListReferencesPaginator(svc, &omics.ListReferencesInput{
				ReferenceStoreId: referenceStore.Id,
			})
			for referencePaginator.HasMorePages() {
				referencesOutput, err := referencePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, reference := range referencesOutput.References {
					resources = append(resources, &OmicsReference{
						svc:              svc,
						ID:               reference.Id,
						Name:             reference.Name,
						ReferenceStoreID: reference.ReferenceStoreId,
					})
				}
			}
		}
	}

	return resources, nil
}

type OmicsReference struct {
	svc              OmicsClient
	ID               *string
	Name             *string
	ReferenceStoreID *string
}

func (r *OmicsReference) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteReference(ctx, &omics.DeleteReferenceInput{
		Id:               r.ID,
		ReferenceStoreId: r.ReferenceStoreID,
	})
	return err
}

func (r *OmicsReference) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsReference) String() string {
	return *r.ID
}

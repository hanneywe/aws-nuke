package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsReadSetResource = "OmicsReadSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsReadSetResource,
		Scope:    nuke.Account,
		Resource: &OmicsReadSet{},
		Lister:   &OmicsReadSetLister{},
	})
}

type OmicsReadSetLister struct {
	svc OmicsClient
}

func (l *OmicsReadSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, list all sequence stores
	sequenceStorePaginator := omics.NewListSequenceStoresPaginator(svc, &omics.ListSequenceStoresInput{})
	for sequenceStorePaginator.HasMorePages() {
		sequenceStoresOutput, err := sequenceStorePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sequenceStore := range sequenceStoresOutput.SequenceStores {
			// For each sequence store, list all read sets
			readSetPaginator := omics.NewListReadSetsPaginator(svc, &omics.ListReadSetsInput{
				SequenceStoreId: sequenceStore.Id,
			})
			for readSetPaginator.HasMorePages() {
				readSetsOutput, err := readSetPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for index := range readSetsOutput.ReadSets {
					readSet := &readSetsOutput.ReadSets[index]
					resources = append(resources, &OmicsReadSet{
						svc:             svc,
						ID:              readSet.Id,
						Name:            readSet.Name,
						SequenceStoreID: readSet.SequenceStoreId,
					})
				}
			}
		}
	}

	return resources, nil
}

type OmicsReadSet struct {
	svc             OmicsClient
	ID              *string
	Name            *string
	SequenceStoreID *string
}

func (r *OmicsReadSet) Remove(ctx context.Context) error {
	_, err := r.svc.BatchDeleteReadSet(ctx, &omics.BatchDeleteReadSetInput{
		SequenceStoreId: r.SequenceStoreID,
		Ids:             []string{*r.ID},
	})
	return err
}

func (r *OmicsReadSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsReadSet) String() string {
	return *r.ID
}

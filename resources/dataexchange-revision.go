package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dataexchange"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DataExchangeRevisionResource = "DataExchangeRevision"

func init() {
	registry.Register(&registry.Registration{
		Name:     DataExchangeRevisionResource,
		Scope:    nuke.Account,
		Resource: &DataExchangeRevision{},
		Lister:   &DataExchangeRevisionLister{},
	})
}

type DataExchangeRevisionLister struct {
	svc DataexchangeClient
}

func (l *DataExchangeRevisionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = dataexchange.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	datasetPaginator := dataexchange.NewListDataSetsPaginator(svc, &dataexchange.ListDataSetsInput{})
	for datasetPaginator.HasMorePages() {
		datasetResp, err := datasetPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for iDataset := range datasetResp.DataSets {
			dataset := &datasetResp.DataSets[iDataset]
			revisionPaginator := dataexchange.NewListDataSetRevisionsPaginator(svc, &dataexchange.ListDataSetRevisionsInput{
				DataSetId: dataset.Id,
			})
			for revisionPaginator.HasMorePages() {
				revisionResp, err := revisionPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for iRevision := range revisionResp.Revisions {
					revision := &revisionResp.Revisions[iRevision]
					resources = append(resources, &DataExchangeRevision{
						svc:        svc,
						DataSetID:  dataset.Id,
						RevisionID: revision.Id,
						ARN:        revision.Arn,
						Comment:    revision.Comment,
					})
				}
			}
		}
	}

	return resources, nil
}

type DataExchangeRevision struct {
	svc        DataexchangeClient
	DataSetID  *string `property:"name=DataSetId"`
	RevisionID *string `property:"name=RevisionId"`
	ARN        *string `property:"name=Arn"`
	Comment    *string
}

func (r *DataExchangeRevision) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRevision(ctx, &dataexchange.DeleteRevisionInput{
		DataSetId:  r.DataSetID,
		RevisionId: r.RevisionID,
	})
	return err
}

func (r *DataExchangeRevision) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DataExchangeRevision) String() string {
	return *r.RevisionID
}

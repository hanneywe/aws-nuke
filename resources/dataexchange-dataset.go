package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dataexchange"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DataExchangeDataSetResource = "DataExchangeDataSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     DataExchangeDataSetResource,
		Scope:    nuke.Account,
		Resource: &DataExchangeDataSet{},
		Lister:   &DataExchangeDataSetLister{},
	})
}

type DataExchangeDataSetLister struct {
	svc DataexchangeClient
}

func (l *DataExchangeDataSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = dataexchange.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := dataexchange.NewListDataSetsPaginator(svc, &dataexchange.ListDataSetsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.DataSets {
			item := &resp.DataSets[i]
			resources = append(resources, &DataExchangeDataSet{
				svc:       svc,
				DataSetID: item.Id,
				Name:      item.Name,
				ARN:       item.Arn,
			})
		}
	}

	return resources, nil
}

type DataExchangeDataSet struct {
	svc       DataexchangeClient
	DataSetID *string `property:"name=DataSetId"`
	Name      *string
	ARN       *string `property:"name=Arn"`
}

func (r *DataExchangeDataSet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataSet(ctx, &dataexchange.DeleteDataSetInput{
		DataSetId: r.DataSetID,
	})
	return err
}

func (r *DataExchangeDataSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DataExchangeDataSet) String() string {
	return *r.DataSetID
}

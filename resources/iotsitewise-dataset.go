package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewiseTypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTSiteWiseDatasetResource = "IoTSiteWiseDataset"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTSiteWiseDatasetResource,
		Scope:    nuke.Account,
		Resource: &IoTSiteWiseDataset{},
		Lister:   &IoTSiteWiseDatasetLister{},
	})
}

type IoTSiteWiseDatasetLister struct {
	svc IotsitewiseClient
}

func (l *IoTSiteWiseDatasetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = iotsitewise.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := iotsitewise.NewListDatasetsPaginator(svc, &iotsitewise.ListDatasetsInput{
		SourceType: iotsitewiseTypes.DatasetSourceTypeKendra,
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.DatasetSummaries {
			resources = append(resources, &IoTSiteWiseDataset{
				svc:       svc,
				DatasetID: item.Id,
				Name:      item.Name,
			})
		}
	}

	return resources, nil
}

type IoTSiteWiseDataset struct {
	svc       IotsitewiseClient
	DatasetID *string `property:"name=DatasetId"`
	Name      *string
}

func (r *IoTSiteWiseDataset) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataset(ctx, &iotsitewise.DeleteDatasetInput{
		DatasetId: r.DatasetID,
	})
	return err
}

func (r *IoTSiteWiseDataset) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTSiteWiseDataset) String() string {
	return *r.DatasetID
}

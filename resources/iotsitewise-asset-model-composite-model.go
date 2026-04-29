package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTSiteWiseAssetModelCompositeModelResource = "IoTSiteWiseAssetModelCompositeModel"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTSiteWiseAssetModelCompositeModelResource,
		Scope:    nuke.Account,
		Resource: &IoTSiteWiseAssetModelCompositeModel{},
		Lister:   &IoTSiteWiseAssetModelCompositeModelLister{},
	})
}

type IoTSiteWiseAssetModelCompositeModelLister struct {
	svc IotsitewiseClient
}

func (l *IoTSiteWiseAssetModelCompositeModelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = iotsitewise.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	modelPaginator := iotsitewise.NewListAssetModelsPaginator(svc, &iotsitewise.ListAssetModelsInput{})
	for modelPaginator.HasMorePages() {
		modelResp, err := modelPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, model := range modelResp.AssetModelSummaries {
			compositePaginator := iotsitewise.NewListAssetModelCompositeModelsPaginator(svc, &iotsitewise.ListAssetModelCompositeModelsInput{
				AssetModelId: model.Id,
			})
			for compositePaginator.HasMorePages() {
				compositeResp, err := compositePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, composite := range compositeResp.AssetModelCompositeModelSummaries {
					resources = append(resources, &IoTSiteWiseAssetModelCompositeModel{
						svc:                        svc,
						AssetModelID:               model.Id,
						AssetModelCompositeModelID: composite.Id,
						Name:                       composite.Name,
						Type:                       composite.Type,
					})
				}
			}
		}
	}

	return resources, nil
}

type IoTSiteWiseAssetModelCompositeModel struct {
	svc                        IotsitewiseClient
	AssetModelID               *string `property:"name=AssetModelId"`
	AssetModelCompositeModelID *string `property:"name=AssetModelCompositeModelId"`
	Name                       *string
	Type                       *string
}

func (r *IoTSiteWiseAssetModelCompositeModel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAssetModelCompositeModel(ctx, &iotsitewise.DeleteAssetModelCompositeModelInput{
		AssetModelId:               r.AssetModelID,
		AssetModelCompositeModelId: r.AssetModelCompositeModelID,
	})
	return err
}

func (r *IoTSiteWiseAssetModelCompositeModel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTSiteWiseAssetModelCompositeModel) String() string {
	return *r.AssetModelCompositeModelID
}

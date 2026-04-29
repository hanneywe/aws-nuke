package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const TransferConnectorResource = "TransferConnector"

func init() {
	registry.Register(&registry.Registration{
		Name:     TransferConnectorResource,
		Scope:    nuke.Account,
		Resource: &TransferConnector{},
		Lister:   &TransferConnectorLister{},
	})
}

type TransferConnectorLister struct {
	svc TransferClient
}

func (l *TransferConnectorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = transfer.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := transfer.NewListConnectorsPaginator(svc, &transfer.ListConnectorsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Connectors {
			item := &resp.Connectors[i]
			resources = append(resources, &TransferConnector{
				svc:         svc,
				ConnectorID: item.ConnectorId,
				URL:         item.Url,
			})
		}
	}

	return resources, nil
}

type TransferConnector struct {
	svc         TransferClient
	ConnectorID *string
	URL         *string
}

func (r *TransferConnector) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConnector(ctx, &transfer.DeleteConnectorInput{
		ConnectorId: r.ConnectorID,
	})
	return err
}

func (r *TransferConnector) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *TransferConnector) String() string {
	return *r.ConnectorID
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	codeconnectionstypes "github.com/aws/aws-sdk-go-v2/service/codeconnections/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeConnectionsConnectionResource = "CodeConnectionsConnection"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeConnectionsConnectionResource,
		Scope:    nuke.Account,
		Resource: &CodeConnectionsConnection{},
		Lister:   &CodeConnectionsConnectionLister{},
	})
}

type CodeConnectionsConnectionLister struct {
	svc CodeConnectionsClient
}

func (l *CodeConnectionsConnectionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = codeconnections.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := codeconnections.NewListConnectionsPaginator(svc, &codeconnections.ListConnectionsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Connections {
			item := &resp.Connections[i]
			resources = append(resources, &CodeConnectionsConnection{
				svc:              svc,
				ConnectionArn:    item.ConnectionArn,
				ConnectionName:   item.ConnectionName,
				ConnectionStatus: item.ConnectionStatus,
				ProviderType:     item.ProviderType,
			})
		}
	}

	return resources, nil
}

type CodeConnectionsConnection struct {
	svc              CodeConnectionsClient
	ConnectionArn    *string
	ConnectionName   *string
	ConnectionStatus codeconnectionstypes.ConnectionStatus
	ProviderType     codeconnectionstypes.ProviderType
}

func (r *CodeConnectionsConnection) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConnection(ctx, &codeconnections.DeleteConnectionInput{
		ConnectionArn: r.ConnectionArn,
	})
	return err
}

func (r *CodeConnectionsConnection) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeConnectionsConnection) String() string {
	return *r.ConnectionName
}

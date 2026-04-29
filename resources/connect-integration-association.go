package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectIntegrationAssociationResource = "ConnectIntegrationAssociation"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectIntegrationAssociationResource,
		Scope:    nuke.Account,
		Resource: &ConnectIntegrationAssociation{},
		Lister:   &ConnectIntegrationAssociationLister{},
	})
}

type ConnectIntegrationAssociationLister struct {
	svc ConnectClient
}

func (l *ConnectIntegrationAssociationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	instancePaginator := connect.NewListInstancesPaginator(svc, &connect.ListInstancesInput{})
	for instancePaginator.HasMorePages() {
		instanceResp, err := instancePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for iInstance := range instanceResp.InstanceSummaryList {
			instance := &instanceResp.InstanceSummaryList[iInstance]
			assocPaginator := connect.NewListIntegrationAssociationsPaginator(svc, &connect.ListIntegrationAssociationsInput{
				InstanceId: instance.Id,
			})
			for assocPaginator.HasMorePages() {
				assocResp, err := assocPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for iAssoc := range assocResp.IntegrationAssociationSummaryList {
					assoc := &assocResp.IntegrationAssociationSummaryList[iAssoc]
					resources = append(resources, &ConnectIntegrationAssociation{
						svc:                      svc,
						InstanceID:               instance.Id,
						IntegrationAssociationID: assoc.IntegrationAssociationId,
						IntegrationARN:           assoc.IntegrationArn,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectIntegrationAssociation struct {
	svc                      ConnectClient
	InstanceID               *string `property:"name=InstanceId"`
	IntegrationAssociationID *string `property:"name=IntegrationAssociationId"`
	IntegrationARN           *string `property:"name=IntegrationArn"`
}

func (r *ConnectIntegrationAssociation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIntegrationAssociation(ctx, &connect.DeleteIntegrationAssociationInput{
		InstanceId:               r.InstanceID,
		IntegrationAssociationId: r.IntegrationAssociationID,
	})
	return err
}

func (r *ConnectIntegrationAssociation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectIntegrationAssociation) String() string {
	return *r.IntegrationAssociationID
}

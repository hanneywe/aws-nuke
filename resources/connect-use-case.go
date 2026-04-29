package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectUseCaseResource = "ConnectUseCase"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectUseCaseResource,
		Scope:    nuke.Account,
		Resource: &ConnectUseCase{},
		Lister:   &ConnectUseCaseLister{},
	})
}

type ConnectUseCaseLister struct {
	svc ConnectClient
}

func (l *ConnectUseCaseLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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

		for _, instance := range instanceResp.InstanceSummaryList {
			iaPaginator := connect.NewListIntegrationAssociationsPaginator(svc, &connect.ListIntegrationAssociationsInput{
				InstanceId: instance.Id,
			})
			for iaPaginator.HasMorePages() {
				iaResp, err := iaPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, ia := range iaResp.IntegrationAssociationSummaryList {
					ucPaginator := connect.NewListUseCasesPaginator(svc, &connect.ListUseCasesInput{
						InstanceId:               instance.Id,
						IntegrationAssociationId: ia.IntegrationAssociationId,
					})
					for ucPaginator.HasMorePages() {
						ucResp, err := ucPaginator.NextPage(ctx)
						if err != nil {
							return nil, err
						}

						for _, uc := range ucResp.UseCaseSummaryList {
							resources = append(resources, &ConnectUseCase{
								svc:                      svc,
								InstanceID:               instance.Id,
								IntegrationAssociationID: ia.IntegrationAssociationId,
								UseCaseID:                uc.UseCaseId,
							})
						}
					}
				}
			}
		}
	}

	return resources, nil
}

type ConnectUseCase struct {
	svc                      ConnectClient
	InstanceID               *string `property:"name=InstanceId"`
	IntegrationAssociationID *string `property:"name=IntegrationAssociationId"`
	UseCaseID                *string `property:"name=UseCaseId"`
}

func (r *ConnectUseCase) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUseCase(ctx, &connect.DeleteUseCaseInput{
		InstanceId:               r.InstanceID,
		IntegrationAssociationId: r.IntegrationAssociationID,
		UseCaseId:                r.UseCaseID,
	})
	return err
}

func (r *ConnectUseCase) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectUseCase) String() string {
	return *r.UseCaseID
}

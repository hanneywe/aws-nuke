package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/deadline"
	deadlinetypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
	"github.com/aws/smithy-go"

	liberrors "github.com/ekristen/libnuke/pkg/errors"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeadlineCloudQueueLimitAssociationResource = "DeadlineCloudQueueLimitAssociation"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeadlineCloudQueueLimitAssociationResource,
		Scope:    nuke.Account,
		Resource: &DeadlineCloudQueueLimitAssociation{},
		Lister:   &DeadlineCloudQueueLimitAssociationLister{},
	})
}

type DeadlineCloudQueueLimitAssociationLister struct {
	svc DeadlineCloudClient
}

func (l *DeadlineCloudQueueLimitAssociationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = deadline.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	farmPaginator := deadline.NewListFarmsPaginator(svc, &deadline.ListFarmsInput{})

	for farmPaginator.HasMorePages() {
		farmResp, err := farmPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, farm := range farmResp.Farms {
			assocPaginator := deadline.NewListQueueLimitAssociationsPaginator(svc, &deadline.ListQueueLimitAssociationsInput{
				FarmId: farm.FarmId,
			})

			for assocPaginator.HasMorePages() {
				assocResp, err := assocPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, assoc := range assocResp.QueueLimitAssociations {
					resources = append(resources, &DeadlineCloudQueueLimitAssociation{
						svc:     svc,
						FarmID:  farm.FarmId,
						QueueID: assoc.QueueId,
						LimitID: assoc.LimitId,
						Status:  string(assoc.Status),
					})
				}
			}
		}
	}

	return resources, nil
}

type DeadlineCloudQueueLimitAssociation struct {
	svc     DeadlineCloudClient
	FarmID  *string `property:"name=FarmId"`
	QueueID *string `property:"name=QueueId"`
	LimitID *string `property:"name=LimitId"`
	Status  string
}

func (r *DeadlineCloudQueueLimitAssociation) Remove(ctx context.Context) error {
	// Re-fetch current status since it may have changed since listing
	resp, err := r.svc.GetQueueLimitAssociation(ctx, &deadline.GetQueueLimitAssociationInput{
		FarmId:  r.FarmID,
		QueueId: r.QueueID,
		LimitId: r.LimitID,
	})
	if err != nil {
		return err
	}
	r.Status = string(resp.Status)

	switch deadlinetypes.QueueLimitAssociationStatus(r.Status) {
	case deadlinetypes.QueueLimitAssociationStatusActive:
		// Must stop the association before it can be deleted.
		// If the association was already moved to a stopping state between
		// the list and this call, the API returns a 409 ConflictException.
		// Treat that as a wait signal.
		_, err := r.svc.UpdateQueueLimitAssociation(ctx, &deadline.UpdateQueueLimitAssociationInput{
			FarmId:  r.FarmID,
			QueueId: r.QueueID,
			LimitId: r.LimitID,
			Status:  deadlinetypes.UpdateQueueLimitAssociationStatusStopLimitUsageAndCancelTasks,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ConflictException" {
				return liberrors.ErrHoldResource("association already transitioning")
			}
			return err
		}
		r.Status = string(deadlinetypes.QueueLimitAssociationStatusStopLimitUsageAndCancelTasks)
		return liberrors.ErrHoldResource("waiting for association to stop")

	case deadlinetypes.QueueLimitAssociationStatusStopLimitUsageAndCompleteTasks,
		deadlinetypes.QueueLimitAssociationStatusStopLimitUsageAndCancelTasks:
		// Still transitioning to STOPPED
		return liberrors.ErrHoldResource("waiting for association to stop")

	case deadlinetypes.QueueLimitAssociationStatusStopped:
		_, err := r.svc.DeleteQueueLimitAssociation(ctx, &deadline.DeleteQueueLimitAssociationInput{
			FarmId:  r.FarmID,
			QueueId: r.QueueID,
			LimitId: r.LimitID,
		})
		return err

	default:
		return fmt.Errorf("unexpected association status: %s", r.Status)
	}
}

func (r *DeadlineCloudQueueLimitAssociation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeadlineCloudQueueLimitAssociation) String() string {
	return *r.LimitID
}

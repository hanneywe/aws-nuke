package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/athena"
	athenatypes "github.com/aws/aws-sdk-go-v2/service/athena/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AthenaCapacityReservationResource = "AthenaCapacityReservation"

func init() {
	registry.Register(&registry.Registration{
		Name:     AthenaCapacityReservationResource,
		Scope:    nuke.Account,
		Resource: &AthenaCapacityReservation{},
		Lister:   &AthenaCapacityReservationLister{},
	})
}

type AthenaCapacityReservationLister struct {
	svc AthenaClient
}

func (l *AthenaCapacityReservationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = athena.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &athena.ListCapacityReservationsInput{}
	for {
		resp, err := svc.ListCapacityReservations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, cr := range resp.CapacityReservations {
			resources = append(resources, &AthenaCapacityReservation{
				svc:    svc,
				Name:   cr.Name,
				Status: cr.Status,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type AthenaCapacityReservation struct {
	svc    AthenaClient
	Name   *string
	Status athenatypes.CapacityReservationStatus
}

func (r *AthenaCapacityReservation) Remove(ctx context.Context) error {
	// A reservation must be canceled before it can be deleted
	if r.Status != athenatypes.CapacityReservationStatusCancelled {
		_, err := r.svc.CancelCapacityReservation(ctx, &athena.CancelCapacityReservationInput{
			Name: r.Name,
		})
		if err != nil {
			return err
		}
		r.Status = athenatypes.CapacityReservationStatusCancelled
	}

	_, err := r.svc.DeleteCapacityReservation(ctx, &athena.DeleteCapacityReservationInput{
		Name: r.Name,
	})
	return err
}

func (r *AthenaCapacityReservation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AthenaCapacityReservation) String() string {
	return *r.Name
}

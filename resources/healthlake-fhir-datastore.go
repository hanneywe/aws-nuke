package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/healthlake"
	healthlaketypes "github.com/aws/aws-sdk-go-v2/service/healthlake/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const HealthLakeFHIRDatastoreResource = "HealthLakeFHIRDatastore"

func init() {
	registry.Register(&registry.Registration{
		Name:     HealthLakeFHIRDatastoreResource,
		Scope:    nuke.Account,
		Resource: &HealthLakeFHIRDatastore{},
		Lister:   &HealthLakeFHIRDatastoreLister{},
	})
}

type HealthLakeFHIRDatastoreLister struct {
	svc HealthLakeClient
}

func (l *HealthLakeFHIRDatastoreLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = healthlake.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := healthlake.NewListFHIRDatastoresPaginator(svc, &healthlake.ListFHIRDatastoresInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, datastore := range output.DatastorePropertiesList {
			datastoreStatus := string(datastore.DatastoreStatus)
			resources = append(resources, &HealthLakeFHIRDatastore{
				svc:           svc,
				DatastoreID:   datastore.DatastoreId,
				DatastoreName: datastore.DatastoreName,
				Status:        &datastoreStatus,
			})
		}
	}

	return resources, nil
}

type HealthLakeFHIRDatastore struct {
	svc           HealthLakeClient
	DatastoreID   *string
	DatastoreName *string
	Status        *string
}

func (r *HealthLakeFHIRDatastore) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFHIRDatastore(ctx, &healthlake.DeleteFHIRDatastoreInput{
		DatastoreId: r.DatastoreID,
	})
	return err
}

func (r *HealthLakeFHIRDatastore) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *HealthLakeFHIRDatastore) String() string {
	return *r.DatastoreID
}

func (r *HealthLakeFHIRDatastore) Filter() error {
	if r.Status != nil && (*r.Status == string(healthlaketypes.DatastoreStatusDeleting) ||
		*r.Status == string(healthlaketypes.DatastoreStatusDeleted)) {
		return fmt.Errorf("already %s", *r.Status)
	}
	return nil
}

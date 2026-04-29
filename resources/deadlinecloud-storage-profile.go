package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/deadline"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DeadlineCloudStorageProfileResource = "DeadlineCloudStorageProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     DeadlineCloudStorageProfileResource,
		Scope:    nuke.Account,
		Resource: &DeadlineCloudStorageProfile{},
		Lister:   &DeadlineCloudStorageProfileLister{},
	})
}

type DeadlineCloudStorageProfileLister struct {
	svc DeadlineCloudClient
}

func (l *DeadlineCloudStorageProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			spPaginator := deadline.NewListStorageProfilesPaginator(svc, &deadline.ListStorageProfilesInput{
				FarmId: farm.FarmId,
			})

			for spPaginator.HasMorePages() {
				spResp, err := spPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, sp := range spResp.StorageProfiles {
					resources = append(resources, &DeadlineCloudStorageProfile{
						svc:              svc,
						FarmID:           farm.FarmId,
						StorageProfileID: sp.StorageProfileId,
						DisplayName:      sp.DisplayName,
					})
				}
			}
		}
	}

	return resources, nil
}

type DeadlineCloudStorageProfile struct {
	svc              DeadlineCloudClient
	FarmID           *string `property:"name=FarmId"`
	StorageProfileID *string `property:"name=StorageProfileId"`
	DisplayName      *string
}

func (r *DeadlineCloudStorageProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStorageProfile(ctx, &deadline.DeleteStorageProfileInput{
		FarmId:           r.FarmID,
		StorageProfileId: r.StorageProfileID,
	})
	return err
}

func (r *DeadlineCloudStorageProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DeadlineCloudStorageProfile) String() string {
	return *r.StorageProfileID
}

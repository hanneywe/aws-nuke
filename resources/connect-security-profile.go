package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectSecurityProfileResource = "ConnectSecurityProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectSecurityProfileResource,
		Scope:    nuke.Account,
		Resource: &ConnectSecurityProfile{},
		Lister:   &ConnectSecurityProfileLister{},
	})
}

type ConnectSecurityProfileLister struct {
	svc ConnectClient
}

func (l *ConnectSecurityProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First list all instances
	instancePaginator := connect.NewListInstancesPaginator(svc, &connect.ListInstancesInput{})

	for instancePaginator.HasMorePages() {
		instanceResp, err := instancePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, instance := range instanceResp.InstanceSummaryList {
			// Then list security profiles for each instance
			profilePaginator := connect.NewListSecurityProfilesPaginator(svc, &connect.ListSecurityProfilesInput{
				InstanceId: instance.Id,
			})

			for profilePaginator.HasMorePages() {
				profileResp, err := profilePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, profile := range profileResp.SecurityProfileSummaryList {
					resources = append(resources, &ConnectSecurityProfile{
						svc:               svc,
						InstanceID:        instance.Id,
						SecurityProfileID: profile.Id,
						Name:              profile.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectSecurityProfile struct {
	svc               ConnectClient
	InstanceID        *string `property:"name=InstanceId"`
	SecurityProfileID *string `property:"name=SecurityProfileId"`
	Name              *string
}

func (r *ConnectSecurityProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSecurityProfile(ctx, &connect.DeleteSecurityProfileInput{
		InstanceId:        r.InstanceID,
		SecurityProfileId: r.SecurityProfileID,
	})
	return err
}

func (r *ConnectSecurityProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectSecurityProfile) String() string {
	return *r.SecurityProfileID
}

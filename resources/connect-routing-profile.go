package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectRoutingProfileResource = "ConnectRoutingProfile"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectRoutingProfileResource,
		Scope:    nuke.Account,
		Resource: &ConnectRoutingProfile{},
		Lister:   &ConnectRoutingProfileLister{},
	})
}

type ConnectRoutingProfileLister struct {
	svc ConnectClient
}

func (l *ConnectRoutingProfileLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			profilePaginator := connect.NewListRoutingProfilesPaginator(svc, &connect.ListRoutingProfilesInput{
				InstanceId: instance.Id,
			})
			for profilePaginator.HasMorePages() {
				profileResp, err := profilePaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for iProfile := range profileResp.RoutingProfileSummaryList {
					profile := &profileResp.RoutingProfileSummaryList[iProfile]
					resources = append(resources, &ConnectRoutingProfile{
						svc:              svc,
						InstanceID:       instance.Id,
						RoutingProfileID: profile.Id,
						Name:             profile.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectRoutingProfile struct {
	svc              ConnectClient
	InstanceID       *string `property:"name=InstanceId"`
	RoutingProfileID *string `property:"name=RoutingProfileId"`
	Name             *string
}

func (r *ConnectRoutingProfile) Filter() error {
	if r.Name != nil && *r.Name == "Basic Routing Profile" {
		return fmt.Errorf("cannot delete default routing profile")
	}
	return nil
}

func (r *ConnectRoutingProfile) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRoutingProfile(ctx, &connect.DeleteRoutingProfileInput{
		InstanceId:       r.InstanceID,
		RoutingProfileId: r.RoutingProfileID,
	})
	return err
}

func (r *ConnectRoutingProfile) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectRoutingProfile) String() string {
	return *r.Name
}

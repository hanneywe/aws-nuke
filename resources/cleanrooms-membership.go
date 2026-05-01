package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CleanRoomsMembershipResource = "CleanRoomsMembership"

func init() {
	registry.Register(&registry.Registration{
		Name:     CleanRoomsMembershipResource,
		Scope:    nuke.Account,
		Resource: &CleanRoomsMembership{},
		Lister:   &CleanRoomsMembershipLister{},
	})
}

type CleanRoomsMembershipLister struct {
	svc CleanRoomsClient
}

func (l *CleanRoomsMembershipLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cleanrooms.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := cleanrooms.NewListMembershipsPaginator(svc, &cleanrooms.ListMembershipsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.MembershipSummaries {
			item := &resp.MembershipSummaries[i]
			resources = append(resources, &CleanRoomsMembership{
				svc:                  svc,
				MembershipIdentifier: item.Id,
				CollaborationName:    item.CollaborationName,
				CollaborationID:      item.CollaborationId,
			})
		}
	}

	return resources, nil
}

type CleanRoomsMembership struct {
	svc                  CleanRoomsClient
	MembershipIdentifier *string
	CollaborationName    *string
	CollaborationID      *string
}

func (r *CleanRoomsMembership) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMembership(ctx, &cleanrooms.DeleteMembershipInput{
		MembershipIdentifier: r.MembershipIdentifier,
	})
	return err
}

func (r *CleanRoomsMembership) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CleanRoomsMembership) String() string {
	return *r.MembershipIdentifier
}

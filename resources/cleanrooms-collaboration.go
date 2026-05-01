package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CleanRoomsCollaborationResource = "CleanRoomsCollaboration"

func init() {
	registry.Register(&registry.Registration{
		Name:     CleanRoomsCollaborationResource,
		Scope:    nuke.Account,
		Resource: &CleanRoomsCollaboration{},
		Lister:   &CleanRoomsCollaborationLister{},
	})
}

type CleanRoomsCollaborationLister struct {
	svc CleanRoomsClient
}

func (l *CleanRoomsCollaborationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cleanrooms.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := cleanrooms.NewListCollaborationsPaginator(svc, &cleanrooms.ListCollaborationsInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, collaboration := range output.CollaborationList {
			resources = append(resources, &CleanRoomsCollaboration{
				svc:                     svc,
				CollaborationIdentifier: collaboration.Id,
				Name:                    collaboration.Name,
				CreatorAccountID:        collaboration.CreatorAccountId,
				AccountID:               opts.AccountID,
			})
		}
	}

	return resources, nil
}

type CleanRoomsCollaboration struct {
	svc                     CleanRoomsClient
	CollaborationIdentifier *string
	Name                    *string
	CreatorAccountID        *string
	AccountID               *string
}

func (r *CleanRoomsCollaboration) Filter() error {
	if r.CreatorAccountID != nil && r.AccountID != nil && *r.CreatorAccountID != *r.AccountID {
		return fmt.Errorf("collaboration owned by account %s, not this account", *r.CreatorAccountID)
	}
	return nil
}

func (r *CleanRoomsCollaboration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCollaboration(ctx, &cleanrooms.DeleteCollaborationInput{
		CollaborationIdentifier: r.CollaborationIdentifier,
	})
	return err
}

func (r *CleanRoomsCollaboration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CleanRoomsCollaboration) String() string {
	return *r.CollaborationIdentifier
}

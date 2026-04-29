package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/notificationscontacts"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NotificationsContactsEmailContactResource = "NotificationsContactsEmailContact"

func init() {
	registry.Register(&registry.Registration{
		Name:     NotificationsContactsEmailContactResource,
		Scope:    nuke.Account,
		Resource: &NotificationsContactsEmailContact{},
		Lister:   &NotificationsContactsEmailContactLister{},
	})
}

type NotificationsContactsEmailContactClient interface {
	ListEmailContacts(ctx context.Context, params *notificationscontacts.ListEmailContactsInput,
		optFns ...func(*notificationscontacts.Options)) (*notificationscontacts.ListEmailContactsOutput, error)
	DeleteEmailContact(ctx context.Context, params *notificationscontacts.DeleteEmailContactInput,
		optFns ...func(*notificationscontacts.Options)) (*notificationscontacts.DeleteEmailContactOutput, error)
}

type NotificationsContactsEmailContactLister struct {
	svc NotificationsContactsEmailContactClient
}

func (l *NotificationsContactsEmailContactLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = notificationscontacts.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := notificationscontacts.NewListEmailContactsPaginator(svc, &notificationscontacts.ListEmailContactsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.EmailContacts {
			ec := &resp.EmailContacts[i]
			resources = append(resources, &NotificationsContactsEmailContact{
				svc:     svc,
				ARN:     ec.Arn,
				Name:    ec.Name,
				Address: ec.Address,
				Status:  string(ec.Status),
			})
		}
	}

	return resources, nil
}

type NotificationsContactsEmailContact struct {
	svc     NotificationsContactsEmailContactClient
	ARN     *string
	Name    *string
	Address *string
	Status  string
}

func (r *NotificationsContactsEmailContact) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEmailContact(ctx, &notificationscontacts.DeleteEmailContactInput{
		Arn: r.ARN,
	})
	return err
}

func (r *NotificationsContactsEmailContact) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NotificationsContactsEmailContact) String() string {
	return *r.Name
}

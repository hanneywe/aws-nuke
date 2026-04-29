package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerAddressListResource = "MailManagerAddressList"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerAddressListResource,
		Scope:    nuke.Account,
		Resource: &MailManagerAddressList{},
		Lister:   &MailManagerAddressListLister{},
	})
}

type MailManagerAddressListLister struct {
	svc MailManagerClient
}

func (l *MailManagerAddressListLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mailmanager.NewListAddressListsPaginator(svc, &mailmanager.ListAddressListsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, al := range resp.AddressLists {
			resources = append(resources, &MailManagerAddressList{
				svc:             svc,
				AddressListID:   al.AddressListId,
				AddressListName: al.AddressListName,
			})
		}
	}
	return resources, nil
}

type MailManagerAddressList struct {
	svc             MailManagerClient
	AddressListID   *string `property:"name=AddressListId"`
	AddressListName *string
}

func (r *MailManagerAddressList) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAddressList(ctx, &mailmanager.DeleteAddressListInput{
		AddressListId: r.AddressListID,
	})
	return err
}

func (r *MailManagerAddressList) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerAddressList) String() string {
	return *r.AddressListName
}

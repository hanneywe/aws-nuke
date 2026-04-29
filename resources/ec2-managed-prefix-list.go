package resources

import (
	"context"
	"fmt"

	"github.com/gotidy/ptr"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2ManagedPrefixListResource = "EC2ManagedPrefixList"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2ManagedPrefixListResource,
		Scope:    nuke.Account,
		Resource: &EC2ManagedPrefixList{},
		Lister:   &EC2ManagedPrefixListLister{},
	})
}

type EC2ManagedPrefixListLister struct {
	svc EC2Client
}

func (l *EC2ManagedPrefixListLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeManagedPrefixListsPaginator(svc,
		&ec2.DescribeManagedPrefixListsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.PrefixLists {
			resources = append(resources, &EC2ManagedPrefixList{
				svc:            svc,
				PrefixListID:   resp.PrefixLists[i].PrefixListId,
				PrefixListName: resp.PrefixLists[i].PrefixListName,
				OwnerID:        resp.PrefixLists[i].OwnerId,
				AddressFamily:  resp.PrefixLists[i].AddressFamily,
				State:          string(resp.PrefixLists[i].State),
				Tags:           resp.PrefixLists[i].Tags,
				accountID:      opts.AccountID,
			})
		}
	}

	return resources, nil
}

type EC2ManagedPrefixList struct {
	svc            EC2Client
	PrefixListID   *string `property:"name=PrefixListId"`
	PrefixListName *string
	OwnerID        *string `property:"name=OwnerId"`
	AddressFamily  *string
	State          string
	Tags           []ec2types.Tag
	accountID      *string
}

func (r *EC2ManagedPrefixList) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteManagedPrefixList(ctx, &ec2.DeleteManagedPrefixListInput{
		PrefixListId: r.PrefixListID,
	})
	return err
}

func (r *EC2ManagedPrefixList) Filter() error {
	if ptr.ToString(r.OwnerID) != ptr.ToString(r.accountID) {
		return fmt.Errorf("cannot delete AWS-managed prefix list")
	}
	if r.State == string(ec2types.PrefixListStateDeleteComplete) {
		return fmt.Errorf("already deleted")
	}
	if r.State == string(ec2types.PrefixListStateDeleteFailed) {
		return fmt.Errorf("delete failed, manual intervention required")
	}
	return nil
}

func (r *EC2ManagedPrefixList) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2ManagedPrefixList) String() string {
	return *r.PrefixListID
}

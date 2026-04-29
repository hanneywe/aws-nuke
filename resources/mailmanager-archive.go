package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
	mailmanagertypes "github.com/aws/aws-sdk-go-v2/service/mailmanager/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerArchiveResource = "MailManagerArchive"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerArchiveResource,
		Scope:    nuke.Account,
		Resource: &MailManagerArchive{},
		Lister:   &MailManagerArchiveLister{},
	})
}

type MailManagerArchiveLister struct {
	svc MailManagerClient
}

func (l *MailManagerArchiveLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mailmanager.NewListArchivesPaginator(svc, &mailmanager.ListArchivesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, ar := range resp.Archives {
			resources = append(resources, &MailManagerArchive{
				svc:          svc,
				ArchiveID:    ar.ArchiveId,
				ArchiveName:  ar.ArchiveName,
				ArchiveState: ar.ArchiveState,
			})
		}
	}
	return resources, nil
}

type MailManagerArchive struct {
	svc          MailManagerClient
	ArchiveID    *string `property:"name=ArchiveId"`
	ArchiveName  *string
	ArchiveState mailmanagertypes.ArchiveState
}

func (r *MailManagerArchive) Filter() error {
	if r.ArchiveState == mailmanagertypes.ArchiveStatePendingDeletion {
		return fmt.Errorf("already pending deletion")
	}
	return nil
}

func (r *MailManagerArchive) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteArchive(ctx, &mailmanager.DeleteArchiveInput{
		ArchiveId: r.ArchiveID,
	})
	return err
}

func (r *MailManagerArchive) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerArchive) String() string {
	return *r.ArchiveName
}

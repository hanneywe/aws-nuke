package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerAddonInstanceResource = "MailManagerAddonInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerAddonInstanceResource,
		Scope:    nuke.Account,
		Resource: &MailManagerAddonInstance{},
		Lister:   &MailManagerAddonInstanceLister{},
	})
}

type MailManagerAddonInstanceLister struct {
	svc MailManagerClient
}

func (l *MailManagerAddonInstanceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := mailmanager.NewListAddonInstancesPaginator(svc, &mailmanager.ListAddonInstancesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, ai := range resp.AddonInstances {
			resources = append(resources, &MailManagerAddonInstance{
				svc:             svc,
				AddonInstanceID: ai.AddonInstanceId,
			})
		}
	}
	return resources, nil
}

type MailManagerAddonInstance struct {
	svc             MailManagerClient
	AddonInstanceID *string `property:"name=AddonInstanceId"`
}

func (r *MailManagerAddonInstance) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAddonInstance(ctx, &mailmanager.DeleteAddonInstanceInput{
		AddonInstanceId: r.AddonInstanceID,
	})
	return err
}

func (r *MailManagerAddonInstance) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerAddonInstance) String() string {
	return *r.AddonInstanceID
}

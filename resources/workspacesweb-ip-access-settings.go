package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const WorkSpacesWebIPAccessSettingsResource = "WorkSpacesWebIPAccessSettings"

func init() {
	registry.Register(&registry.Registration{
		Name:     WorkSpacesWebIPAccessSettingsResource,
		Scope:    nuke.Account,
		Resource: &WorkSpacesWebIPAccessSettings{},
		Lister:   &WorkSpacesWebIPAccessSettingsLister{},
	})
}

type WorkSpacesWebIPAccessSettingsLister struct {
	svc WorkSpacesWebClient
}

func (l *WorkSpacesWebIPAccessSettingsLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = workspacesweb.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := workspacesweb.NewListIpAccessSettingsPaginator(svc, &workspacesweb.ListIpAccessSettingsInput{})
	for paginator.HasMorePages() {
		listIPAccessSettingsOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for idx := range listIPAccessSettingsOutput.IpAccessSettings {
			ipAccessSettings := &listIPAccessSettingsOutput.IpAccessSettings[idx]
			resources = append(resources, &WorkSpacesWebIPAccessSettings{
				svc:                 svc,
				IPAccessSettingsArn: ipAccessSettings.IpAccessSettingsArn,
				DisplayName:         ipAccessSettings.DisplayName,
			})
		}
	}

	return resources, nil
}

type WorkSpacesWebIPAccessSettings struct {
	svc                 WorkSpacesWebClient
	IPAccessSettingsArn *string
	DisplayName         *string
}

func (r *WorkSpacesWebIPAccessSettings) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIpAccessSettings(ctx, &workspacesweb.DeleteIpAccessSettingsInput{
		IpAccessSettingsArn: r.IPAccessSettingsArn,
	})
	return err
}

func (r *WorkSpacesWebIPAccessSettings) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *WorkSpacesWebIPAccessSettings) String() string {
	return *r.IPAccessSettingsArn
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/workspacesweb"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const WorkSpacesWebDataProtectionSettingsResource = "WorkSpacesWebDataProtectionSettings"

func init() {
	registry.Register(&registry.Registration{
		Name:     WorkSpacesWebDataProtectionSettingsResource,
		Scope:    nuke.Account,
		Resource: &WorkSpacesWebDataProtectionSettings{},
		Lister:   &WorkSpacesWebDataProtectionSettingsLister{},
	})
}

type WorkSpacesWebDataProtectionSettingsLister struct {
	svc WorkSpacesWebClient
}

func (l *WorkSpacesWebDataProtectionSettingsLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = workspacesweb.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &workspacesweb.ListDataProtectionSettingsInput{}
	for {
		output, err := svc.ListDataProtectionSettings(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, dataProtectionSettings := range output.DataProtectionSettings {
			resources = append(resources, &WorkSpacesWebDataProtectionSettings{
				svc:                       svc,
				DataProtectionSettingsArn: dataProtectionSettings.DataProtectionSettingsArn,
				DisplayName:               dataProtectionSettings.DisplayName,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type WorkSpacesWebDataProtectionSettings struct {
	svc                       WorkSpacesWebClient
	DataProtectionSettingsArn *string
	DisplayName               *string
}

func (r *WorkSpacesWebDataProtectionSettings) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataProtectionSettings(ctx, &workspacesweb.DeleteDataProtectionSettingsInput{
		DataProtectionSettingsArn: r.DataProtectionSettingsArn,
	})
	return err
}

func (r *WorkSpacesWebDataProtectionSettings) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *WorkSpacesWebDataProtectionSettings) String() string {
	return *r.DataProtectionSettingsArn
}

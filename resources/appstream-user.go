package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
	appstreamtypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppStreamUserResource = "AppStreamUser"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppStreamUserResource,
		Scope:    nuke.Account,
		Resource: &AppStreamUser{},
		Lister:   &AppStreamUserLister{},
	})
}

type AppStreamUserLister struct {
	svc AppStreamClient
}

func (l *AppStreamUserLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = appstream.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &appstream.DescribeUsersInput{
		AuthenticationType: appstreamtypes.AuthenticationTypeUserpool,
	}

	for {
		resp, err := svc.DescribeUsers(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Users {
			u := resp.Users[i]
			resources = append(resources, &AppStreamUser{
				svc:                svc,
				UserName:           u.UserName,
				AuthenticationType: u.AuthenticationType,
				Enabled:            u.Enabled,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type AppStreamUser struct {
	svc                AppStreamClient
	UserName           *string
	AuthenticationType appstreamtypes.AuthenticationType
	Enabled            *bool
}

func (r *AppStreamUser) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteUser(ctx, &appstream.DeleteUserInput{
		UserName:           r.UserName,
		AuthenticationType: r.AuthenticationType,
	})
	return err
}

func (r *AppStreamUser) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppStreamUser) String() string {
	return *r.UserName
}

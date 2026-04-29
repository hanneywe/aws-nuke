package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CognitoUserPoolGroupResource = "CognitoUserPoolGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     CognitoUserPoolGroupResource,
		Scope:    nuke.Account,
		Resource: &CognitoUserPoolGroup{},
		Lister:   &CognitoUserPoolGroupLister{},
	})
}

type CognitoUserPoolGroupLister struct {
	svc CognitoidentityproviderClient
}

func (l *CognitoUserPoolGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cognitoidentityprovider.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	poolPaginator := cognitoidentityprovider.NewListUserPoolsPaginator(svc, &cognitoidentityprovider.ListUserPoolsInput{})
	for poolPaginator.HasMorePages() {
		poolResp, err := poolPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, pool := range poolResp.UserPools {
			groupPaginator := cognitoidentityprovider.NewListGroupsPaginator(svc, &cognitoidentityprovider.ListGroupsInput{
				UserPoolId: pool.Id,
			})
			for groupPaginator.HasMorePages() {
				groupResp, err := groupPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, group := range groupResp.Groups {
					resources = append(resources, &CognitoUserPoolGroup{
						svc:         svc,
						UserPoolID:  pool.Id,
						GroupName:   group.GroupName,
						Description: group.Description,
					})
				}
			}
		}
	}

	return resources, nil
}

type CognitoUserPoolGroup struct {
	svc         CognitoidentityproviderClient
	UserPoolID  *string `property:"name=UserPoolId"`
	GroupName   *string
	Description *string
}

func (r *CognitoUserPoolGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
		GroupName:  r.GroupName,
		UserPoolId: r.UserPoolID,
	})
	return err
}

func (r *CognitoUserPoolGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CognitoUserPoolGroup) String() string {
	return *r.GroupName
}

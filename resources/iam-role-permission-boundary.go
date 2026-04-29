package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/iam" //nolint:staticcheck
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IAMRolePermissionBoundaryResource = "IAMRolePermissionBoundary"

func init() {
	registry.Register(&registry.Registration{
		Name:     IAMRolePermissionBoundaryResource,
		Scope:    nuke.Account,
		Resource: &IAMRolePermissionBoundary{},
		Lister:   &IAMRolePermissionBoundaryLister{},
	})
}

type IAMRolePermissionBoundaryLister struct{}

func (l *IAMRolePermissionBoundaryLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := iam.New(opts.Session)

	var resources []resource.Resource

	params := &iam.ListRolesInput{}
	for {
		resp, err := svc.ListRoles(params)
		if err != nil {
			return nil, err
		}

		for _, listedRole := range resp.Roles {
			role, err := GetIAMRole(svc, listedRole.RoleName)
			if err != nil {
				logrus.Errorf("Failed to get role %s: %v", *listedRole.RoleName, err)
				continue
			}

			if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != nil {
				resources = append(resources, &IAMRolePermissionBoundary{
					svc:                    svc,
					RoleName:               role.RoleName,
					RolePath:               role.Path,
					RoleCreateDate:         role.CreateDate,
					RoleLastUsed:           getLastUsedDate(role),
					PermissionsBoundaryARN: role.PermissionsBoundary.PermissionsBoundaryArn,
					RoleTags:               role.Tags,
				})
			}
		}

		if !*resp.IsTruncated {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

type IAMRolePermissionBoundary struct {
	svc                    iamiface.IAMAPI
	RoleName               *string
	RolePath               *string
	RoleCreateDate         *time.Time
	RoleLastUsed           *time.Time
	PermissionsBoundaryARN *string
	RoleTags               []*iam.Tag
}

func (r *IAMRolePermissionBoundary) Filter() error {
	if strings.HasPrefix(*r.RolePath, "/aws-service-role/") {
		return fmt.Errorf("cannot modify service-linked roles")
	}
	if strings.HasPrefix(*r.RolePath, "/aws-reserved/sso.amazonaws.com/") {
		return fmt.Errorf("cannot modify SSO roles")
	}
	return nil
}

func (r *IAMRolePermissionBoundary) Remove(_ context.Context) error {
	_, err := r.svc.DeleteRolePermissionsBoundary(&iam.DeleteRolePermissionsBoundaryInput{
		RoleName: r.RoleName,
	})
	return err
}

func (r *IAMRolePermissionBoundary) Properties() types.Properties {
	properties := types.NewProperties().
		Set("RoleName", r.RoleName).
		Set("RolePath", r.RolePath).
		Set("RoleCreateDate", r.RoleCreateDate.Format(time.RFC3339)).
		Set("RoleLastUsed", r.RoleLastUsed).
		Set("PermissionsBoundaryARN", r.PermissionsBoundaryARN)

	for _, tag := range r.RoleTags {
		properties.SetTagWithPrefix("role", tag.Key, tag.Value)
	}
	return properties
}

func (r *IAMRolePermissionBoundary) String() string {
	return fmt.Sprintf("%s -> %s", *r.RoleName, *r.PermissionsBoundaryARN)
}

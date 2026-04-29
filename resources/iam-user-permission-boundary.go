package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/iam" //nolint:staticcheck
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IAMUserPermissionBoundaryResource = "IAMUserPermissionBoundary"

func init() {
	registry.Register(&registry.Registration{
		Name:     IAMUserPermissionBoundaryResource,
		Scope:    nuke.Account,
		Resource: &IAMUserPermissionBoundary{},
		Lister:   &IAMUserPermissionBoundaryLister{},
	})
}

type IAMUserPermissionBoundaryLister struct{}

func (l *IAMUserPermissionBoundaryLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := iam.New(opts.Session)

	var resources []resource.Resource

	allUsers, err := ListIAMUsers(svc)
	if err != nil {
		return nil, err
	}

	for _, out := range allUsers {
		user, err := GetIAMUser(svc, out.UserName)
		if err != nil {
			logrus.Errorf("Failed to get user %s: %v", *out.UserName, err)
			continue
		}

		if user.PermissionsBoundary != nil && user.PermissionsBoundary.PermissionsBoundaryArn != nil {
			resources = append(resources, &IAMUserPermissionBoundary{
				svc:                    svc,
				UserName:               user.UserName,
				UserPath:               user.Path,
				UserCreateDate:         user.CreateDate,
				PermissionsBoundaryARN: user.PermissionsBoundary.PermissionsBoundaryArn,
				UserTags:               user.Tags,
			})
		}
	}

	return resources, nil
}

type IAMUserPermissionBoundary struct {
	svc                    iamiface.IAMAPI
	UserName               *string
	UserPath               *string
	UserCreateDate         *time.Time
	PermissionsBoundaryARN *string
	UserTags               []*iam.Tag
}

func (r *IAMUserPermissionBoundary) Remove(_ context.Context) error {
	_, err := r.svc.DeleteUserPermissionsBoundary(&iam.DeleteUserPermissionsBoundaryInput{
		UserName: r.UserName,
	})
	return err
}

func (r *IAMUserPermissionBoundary) Properties() types.Properties {
	properties := types.NewProperties().
		Set("UserName", r.UserName).
		Set("UserPath", r.UserPath).
		Set("UserCreateDate", r.UserCreateDate.Format(time.RFC3339)).
		Set("PermissionsBoundaryARN", r.PermissionsBoundaryARN)

	for _, tag := range r.UserTags {
		properties.SetTagWithPrefix("user", tag.Key, tag.Value)
	}
	return properties
}

func (r *IAMUserPermissionBoundary) String() string {
	return fmt.Sprintf("%s -> %s", *r.UserName, *r.PermissionsBoundaryARN)
}

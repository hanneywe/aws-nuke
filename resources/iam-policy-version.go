package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"         //nolint:staticcheck
	"github.com/aws/aws-sdk-go/service/iam" //nolint:staticcheck
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IAMPolicyVersionResource = "IAMPolicyVersion"

func init() {
	registry.Register(&registry.Registration{
		Name:     IAMPolicyVersionResource,
		Scope:    nuke.Account,
		Resource: &IAMPolicyVersion{},
		Lister:   &IAMPolicyVersionLister{},
	})
}

type IAMPolicyVersionLister struct{}

func (l *IAMPolicyVersionLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := iam.New(opts.Session)

	var resources []resource.Resource

	if err := svc.ListPoliciesPages(&iam.ListPoliciesInput{
		Scope: aws.String("Local"),
	}, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		for _, policy := range page.Policies {
			versResp, err := svc.ListPolicyVersions(&iam.ListPolicyVersionsInput{
				PolicyArn: policy.Arn,
			})
			if err != nil {
				logrus.Errorf("Failed to list versions for policy %s: %v", *policy.PolicyName, err)
				continue
			}

			for _, version := range versResp.Versions {
				if !*version.IsDefaultVersion {
					resources = append(resources, &IAMPolicyVersion{
						svc:        svc,
						PolicyARN:  policy.Arn,
						PolicyName: policy.PolicyName,
						VersionID:  version.VersionId,
						CreateDate: version.CreateDate,
					})
				}
			}
		}
		return true
	}); err != nil {
		return nil, err
	}

	return resources, nil
}

type IAMPolicyVersion struct {
	svc        iamiface.IAMAPI
	PolicyARN  *string
	PolicyName *string
	VersionID  *string
	CreateDate *time.Time
}

func (r *IAMPolicyVersion) Remove(_ context.Context) error {
	_, err := r.svc.DeletePolicyVersion(&iam.DeletePolicyVersionInput{
		PolicyArn: r.PolicyARN,
		VersionId: r.VersionID,
	})
	return err
}

func (r *IAMPolicyVersion) Properties() types.Properties {
	return types.NewProperties().
		Set("PolicyARN", r.PolicyARN).
		Set("PolicyName", r.PolicyName).
		Set("VersionID", r.VersionID).
		Set("CreateDate", r.CreateDate.Format(time.RFC3339))
}

func (r *IAMPolicyVersion) String() string {
	return fmt.Sprintf("%s -> %s", *r.PolicyName, *r.VersionID)
}

package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTPolicyVersionResource = "IoTPolicyVersion"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTPolicyVersionResource,
		Scope:    nuke.Account,
		Resource: &IoTPolicyVersion{},
		Lister:   &IoTPolicyVersionLister{},
	})
}

type IoTPolicyVersionLister struct {
	svc IoTClient
}

func (l *IoTPolicyVersionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	policyPaginator := iot.NewListPoliciesPaginator(svc, &iot.ListPoliciesInput{})
	for policyPaginator.HasMorePages() {
		policyResp, err := policyPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for iPolicy := range policyResp.Policies {
			policy := &policyResp.Policies[iPolicy]
			versionResp, err := svc.ListPolicyVersions(ctx, &iot.ListPolicyVersionsInput{
				PolicyName: policy.PolicyName,
			})
			if err != nil {
				return nil, err
			}
			for iVersion := range versionResp.PolicyVersions {
				version := &versionResp.PolicyVersions[iVersion]
				resources = append(resources, &IoTPolicyVersion{
					svc:              svc,
					PolicyName:       policy.PolicyName,
					PolicyVersionID:  version.VersionId,
					IsDefaultVersion: version.IsDefaultVersion,
				})
			}
		}
	}

	return resources, nil
}

type IoTPolicyVersion struct {
	svc              IoTClient
	PolicyName       *string
	PolicyVersionID  *string `property:"name=PolicyVersionId"`
	IsDefaultVersion bool
}

func (r *IoTPolicyVersion) Filter() error {
	if r.IsDefaultVersion {
		return fmt.Errorf("cannot delete default policy version")
	}
	return nil
}

func (r *IoTPolicyVersion) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePolicyVersion(ctx, &iot.DeletePolicyVersionInput{
		PolicyName:      r.PolicyName,
		PolicyVersionId: r.PolicyVersionID,
	})
	return err
}

func (r *IoTPolicyVersion) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTPolicyVersion) String() string {
	return *r.PolicyVersionID
}

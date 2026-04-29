package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeDeployOnPremisesInstanceResource = "CodeDeployOnPremisesInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeDeployOnPremisesInstanceResource,
		Scope:    nuke.Account,
		Resource: &CodeDeployOnPremisesInstance{},
		Lister:   &CodeDeployOnPremisesInstanceLister{},
	})
}

type CodeDeployOnPremisesInstanceLister struct {
	svc CodeDeployV2Client
}

func (l *CodeDeployOnPremisesInstanceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = codedeploy.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &codedeploy.ListOnPremisesInstancesInput{}
	for {
		resp, err := svc.ListOnPremisesInstances(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, name := range resp.InstanceNames {
			resources = append(resources, &CodeDeployOnPremisesInstance{
				svc:          svc,
				InstanceName: aws.String(name),
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type CodeDeployOnPremisesInstance struct {
	svc          CodeDeployV2Client
	InstanceName *string
}

func (r *CodeDeployOnPremisesInstance) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterOnPremisesInstance(ctx, &codedeploy.DeregisterOnPremisesInstanceInput{
		InstanceName: r.InstanceName,
	})
	return err
}

func (r *CodeDeployOnPremisesInstance) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeDeployOnPremisesInstance) String() string {
	return *r.InstanceName
}

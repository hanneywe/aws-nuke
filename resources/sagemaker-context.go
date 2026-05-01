package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerContextResource = "SageMakerContext"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerContextResource,
		Scope:    nuke.Account,
		Resource: &SageMakerContext{},
		Lister:   &SageMakerContextLister{},
	})
}

type SageMakerContextLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerContextLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListContextsInput{}

	for {
		output, err := svc.ListContexts(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, sagemakerContext := range output.ContextSummaries {
			resources = append(resources, &SageMakerContext{
				svc:         svc,
				ContextName: sagemakerContext.ContextName,
				ContextArn:  sagemakerContext.ContextArn,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SageMakerContext struct {
	svc SageMakerV2Client

	ContextName *string
	ContextArn  *string
}

func (r *SageMakerContext) Remove(ctx context.Context) error {
	// Delete all associations where this context is the source or destination
	for _, arn := range []*string{r.ContextArn} {
		if arn == nil {
			continue
		}
		for _, direction := range []struct {
			params *sagemaker.ListAssociationsInput
		}{
			{params: &sagemaker.ListAssociationsInput{SourceArn: arn}},
			{params: &sagemaker.ListAssociationsInput{DestinationArn: arn}},
		} {
			resp, err := r.svc.ListAssociations(ctx, direction.params)
			if err != nil {
				return err
			}
			for _, assoc := range resp.AssociationSummaries {
				_, err := r.svc.DeleteAssociation(ctx, &sagemaker.DeleteAssociationInput{
					SourceArn:      assoc.SourceArn,
					DestinationArn: assoc.DestinationArn,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	_, err := r.svc.DeleteContext(ctx, &sagemaker.DeleteContextInput{
		ContextName: r.ContextName,
	})
	return err
}

func (r *SageMakerContext) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerContext) String() string {
	return *r.ContextName
}

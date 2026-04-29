package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	codeconnectionstypes "github.com/aws/aws-sdk-go-v2/service/codeconnections/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeConnectionsRepositoryLinkResource = "CodeConnectionsRepositoryLink"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeConnectionsRepositoryLinkResource,
		Scope:    nuke.Account,
		Resource: &CodeConnectionsRepositoryLink{},
		Lister:   &CodeConnectionsRepositoryLinkLister{},
	})
}

type CodeConnectionsRepositoryLinkLister struct {
	svc CodeConnectionsClient
}

func (l *CodeConnectionsRepositoryLinkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = codeconnections.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListRepositoryLinks(ctx, &codeconnections.ListRepositoryLinksInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, link := range resp.RepositoryLinks {
			resources = append(resources, &CodeConnectionsRepositoryLink{
				svc:              svc,
				RepositoryLinkID: link.RepositoryLinkId,
				RepositoryName:   link.RepositoryName,
				ProviderType:     link.ProviderType,
				ConnectionArn:    link.ConnectionArn,
			})
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return resources, nil
}

type CodeConnectionsRepositoryLink struct {
	svc              CodeConnectionsClient
	RepositoryLinkID *string
	RepositoryName   *string
	ProviderType     codeconnectionstypes.ProviderType
	ConnectionArn    *string
}

func (r *CodeConnectionsRepositoryLink) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRepositoryLink(ctx, &codeconnections.DeleteRepositoryLinkInput{
		RepositoryLinkId: r.RepositoryLinkID,
	})
	return err
}

func (r *CodeConnectionsRepositoryLink) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeConnectionsRepositoryLink) String() string {
	return *r.RepositoryLinkID
}

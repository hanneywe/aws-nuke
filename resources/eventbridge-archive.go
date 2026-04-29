package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EventBridgeArchiveResource = "EventBridgeArchive"

func init() {
	registry.Register(&registry.Registration{
		Name:     EventBridgeArchiveResource,
		Scope:    nuke.Account,
		Resource: &EventBridgeArchive{},
		Lister:   &EventBridgeArchiveLister{},
	})
}

type EventBridgeArchiveLister struct {
	svc EventBridgeClient
}

func (l *EventBridgeArchiveLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = eventbridge.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &eventbridge.ListArchivesInput{}

	for {
		output, err := svc.ListArchives(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, archive := range output.Archives {
			resources = append(resources, &EventBridgeArchive{
				svc:         svc,
				ArchiveName: archive.ArchiveName,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type EventBridgeArchive struct {
	svc         EventBridgeClient
	ArchiveName *string
}

func (r *EventBridgeArchive) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteArchive(ctx, &eventbridge.DeleteArchiveInput{
		ArchiveName: r.ArchiveName,
	})
	return err
}

func (r *EventBridgeArchive) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EventBridgeArchive) String() string {
	return *r.ArchiveName
}

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

const IoTCommandResource = "IoTCommand"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTCommandResource,
		Scope:    nuke.Account,
		Resource: &IoTCommand{},
		Lister:   &IoTCommandLister{},
	})
}

type IoTCommandLister struct {
	svc IoTClient
}

func (l *IoTCommandLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iot.NewListCommandsPaginator(svc, &iot.ListCommandsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cmd := range resp.Commands {
			resources = append(resources, &IoTCommand{
				svc:             svc,
				CommandID:       cmd.CommandId,
				CommandArn:      cmd.CommandArn,
				PendingDeletion: cmd.PendingDeletion,
			})
		}
	}

	return resources, nil
}

type IoTCommand struct {
	svc             IoTClient
	CommandID       *string `property:"name=CommandId"`
	CommandArn      *string
	PendingDeletion *bool
}

func (r *IoTCommand) Filter() error {
	if r.PendingDeletion != nil && *r.PendingDeletion {
		return fmt.Errorf("already pending deletion")
	}
	return nil
}

func (r *IoTCommand) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCommand(ctx, &iot.DeleteCommandInput{
		CommandId: r.CommandID,
	})
	return err
}

func (r *IoTCommand) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTCommand) String() string {
	return *r.CommandID
}

package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSChatRoomResource = "IVSChatRoom"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSChatRoomResource,
		Scope:    nuke.Account,
		Resource: &IVSChatRoom{},
		Lister:   &IVSChatRoomLister{},
	})
}

type IVSChatRoomLister struct {
	svc IVSChatClient
}

func (l *IVSChatRoomLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivschat.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivschat.NewListRoomsPaginator(svc, &ivschat.ListRoomsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, room := range resp.Rooms {
			resources = append(resources, &IVSChatRoom{
				svc:  svc,
				ARN:  room.Arn,
				Name: room.Name,
				Tags: room.Tags,
			})
		}
	}

	return resources, nil
}

type IVSChatRoom struct {
	svc  IVSChatClient
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *IVSChatRoom) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRoom(ctx, &ivschat.DeleteRoomInput{
		Identifier: r.ARN,
	})
	return err
}

func (r *IVSChatRoom) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSChatRoom) String() string {
	return *r.Name
}

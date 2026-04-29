package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailContainerService_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetContainerServices", mock.Anything, mock.Anything).
		Return(&lightsail.GetContainerServicesOutput{
			ContainerServices: []lstypes.ContainerService{
				{
					ContainerServiceName: ptr.String("my-container-svc"),
				},
			},
		}, nil)

	lister := &LightsailContainerServiceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cs := resources[0].(*LightsailContainerService)
	a.Equal("my-container-svc", *cs.ContainerServiceName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContainerService_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetContainerServices", mock.Anything, mock.Anything).
		Return(&lightsail.GetContainerServicesOutput{
			ContainerServices: []lstypes.ContainerService{},
		}, nil)

	lister := &LightsailContainerServiceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContainerService_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	cs := &LightsailContainerService{
		svc:                  mockClient,
		ContainerServiceName: ptr.String("my-container-svc"),
	}

	mockClient.On("DeleteContainerService", mock.Anything, &lightsail.DeleteContainerServiceInput{
		ServiceName: cs.ContainerServiceName,
	}).Return(&lightsail.DeleteContainerServiceOutput{}, nil)

	a.NoError(cs.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContainerService_Properties(t *testing.T) {
	a := assert.New(t)

	cs := LightsailContainerService{
		ContainerServiceName: ptr.String("my-container-svc"),
	}

	props := cs.Properties()
	a.Equal("my-container-svc", props.Get("ContainerServiceName"))
}

func Test_Mock_LightsailContainerService_String(t *testing.T) {
	a := assert.New(t)
	cs := LightsailContainerService{ContainerServiceName: ptr.String("my-container-svc")}
	a.Equal("my-container-svc", cs.String())
}

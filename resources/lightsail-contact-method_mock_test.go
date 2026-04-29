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

func Test_Mock_LightsailContactMethod_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetContactMethods", mock.Anything, mock.Anything).
		Return(&lightsail.GetContactMethodsOutput{
			ContactMethods: []lstypes.ContactMethod{
				{
					Protocol:        lstypes.ContactProtocolEmail,
					ContactEndpoint: ptr.String("test@example.com"),
				},
			},
		}, nil)

	lister := &LightsailContactMethodLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cm := resources[0].(*LightsailContactMethod)
	a.Equal("Email", *cm.Protocol)
	a.Equal("test@example.com", *cm.ContactEndpoint)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContactMethod_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetContactMethods", mock.Anything, mock.Anything).
		Return(&lightsail.GetContactMethodsOutput{
			ContactMethods: []lstypes.ContactMethod{},
		}, nil)

	lister := &LightsailContactMethodLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContactMethod_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	cm := &LightsailContactMethod{
		svc:             mockClient,
		Protocol:        ptr.String("Email"),
		ContactEndpoint: ptr.String("test@example.com"),
	}

	mockClient.On("DeleteContactMethod", mock.Anything, &lightsail.DeleteContactMethodInput{
		Protocol: lstypes.ContactProtocolEmail,
	}).Return(&lightsail.DeleteContactMethodOutput{}, nil)

	a.NoError(cm.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailContactMethod_Properties(t *testing.T) {
	a := assert.New(t)

	cm := LightsailContactMethod{
		Protocol:        ptr.String("Email"),
		ContactEndpoint: ptr.String("test@example.com"),
	}

	props := cm.Properties()
	a.Equal("Email", props.Get("Protocol"))
	a.Equal("test@example.com", props.Get("ContactEndpoint"))
}

func Test_Mock_LightsailContactMethod_String(t *testing.T) {
	a := assert.New(t)
	cm := LightsailContactMethod{ContactEndpoint: ptr.String("test@example.com")}
	a.Equal("test@example.com", cm.String())
}

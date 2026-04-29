package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appintegrations"
	aitypes "github.com/aws/aws-sdk-go-v2/service/appintegrations/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testAppIntegrationsListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_AppIntegrationsEventIntegration_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAppIntegrationsClient)

	mockClient.
		On("ListEventIntegrations", mock.Anything, mock.Anything).
		Return(&appintegrations.ListEventIntegrationsOutput{
			EventIntegrations: []aitypes.EventIntegration{
				{
					Name:                ptr.String("my-event-integration"),
					EventIntegrationArn: ptr.String("arn:aws:app-integrations:us-east-1:123456789012:event-integration/my-event-integration"),
				},
			},
		}, nil)

	lister := &AppIntegrationsEventIntegrationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testAppIntegrationsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	ei := resources[0].(*AppIntegrationsEventIntegration)
	assertions.Equal("my-event-integration", *ei.Name)
	assertions.Equal("arn:aws:app-integrations:us-east-1:123456789012:event-integration/my-event-integration", *ei.EventIntegrationArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AppIntegrationsEventIntegration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAppIntegrationsClient)

	mockClient.
		On("ListEventIntegrations", mock.Anything, mock.Anything).
		Return(&appintegrations.ListEventIntegrationsOutput{
			EventIntegrations: []aitypes.EventIntegration{},
		}, nil)

	lister := &AppIntegrationsEventIntegrationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testAppIntegrationsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AppIntegrationsEventIntegration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockAppIntegrationsClient)

	ei := &AppIntegrationsEventIntegration{
		svc:  mockClient,
		Name: ptr.String("my-event-integration"),
	}

	mockClient.
		On("DeleteEventIntegration", mock.Anything, &appintegrations.DeleteEventIntegrationInput{
			Name: ei.Name,
		}).
		Return(&appintegrations.DeleteEventIntegrationOutput{}, nil)

	err := ei.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_AppIntegrationsEventIntegration_Properties(t *testing.T) {
	assertions := assert.New(t)

	ei := AppIntegrationsEventIntegration{
		Name:                ptr.String("my-event-integration"),
		EventIntegrationArn: ptr.String("arn:aws:app-integrations:us-east-1:123456789012:event-integration/my-event-integration"),
	}

	properties := ei.Properties()

	assertions.Equal("my-event-integration", properties.Get("Name"))
	expectedArn := "arn:aws:app-integrations:us-east-1:123456789012:event-integration/my-event-integration"
	assertions.Equal(expectedArn, properties.Get("EventIntegrationArn"))
}

func Test_Mock_AppIntegrationsEventIntegration_String(t *testing.T) {
	assertions := assert.New(t)

	ei := AppIntegrationsEventIntegration{
		Name: ptr.String("my-event-integration"),
	}

	assertions.Equal("my-event-integration", ei.String())
}

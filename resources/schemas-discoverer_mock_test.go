package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/schemas"
	schemastypes "github.com/aws/aws-sdk-go-v2/service/schemas/types"
)

func Test_Mock_SchemasDiscoverer_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	mockClient.On("ListDiscoverers", mock.Anything, mock.Anything).
		Return(&schemas.ListDiscoverersOutput{
			Discoverers: []schemastypes.DiscovererSummary{
				{
					DiscovererId: ptr.String("d-123456"),
					SourceArn:    ptr.String("arn:aws:events:us-east-1:123456789012:event-bus/default"),
					Tags:         map[string]string{"env": "test"},
				},
			},
		}, nil)
	lister := &SchemasDiscovererLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchemasListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	d := resources[0].(*SchemasDiscoverer)
	a.Equal("d-123456", *d.DiscovererID)
	a.Equal("arn:aws:events:us-east-1:123456789012:event-bus/default", *d.SourceArn)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasDiscoverer_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	mockClient.On("ListDiscoverers", mock.Anything, mock.Anything).
		Return(&schemas.ListDiscoverersOutput{Discoverers: []schemastypes.DiscovererSummary{}}, nil)
	lister := &SchemasDiscovererLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchemasListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasDiscoverer_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	d := &SchemasDiscoverer{svc: mockClient, DiscovererID: ptr.String("d-123456")}
	mockClient.On("DeleteDiscoverer", mock.Anything, &schemas.DeleteDiscovererInput{DiscovererId: d.DiscovererID}).
		Return(&schemas.DeleteDiscovererOutput{}, nil)
	a.NoError(d.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasDiscoverer_Properties(t *testing.T) {
	a := assert.New(t)
	d := SchemasDiscoverer{
		DiscovererID: ptr.String("d-123456"),
		SourceArn:    ptr.String("arn:aws:events:us-east-1:123456789012:event-bus/default"),
		Tags:         map[string]string{"env": "test"},
	}
	props := d.Properties()
	a.Equal("d-123456", props.Get("DiscovererID"))
	a.Equal("arn:aws:events:us-east-1:123456789012:event-bus/default", props.Get("SourceArn"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_SchemasDiscoverer_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("d-123456", (&SchemasDiscoverer{DiscovererID: ptr.String("d-123456")}).String())
}

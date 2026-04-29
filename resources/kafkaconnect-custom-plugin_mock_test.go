package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"
	kafkaconnecttypes "github.com/aws/aws-sdk-go-v2/service/kafkaconnect/types"
)

func Test_Mock_KafkaConnectCustomPlugin_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	mockClient.On("ListCustomPlugins", mock.Anything, mock.Anything).
		Return(&kafkaconnect.ListCustomPluginsOutput{
			CustomPlugins: []kafkaconnecttypes.CustomPluginSummary{
				{
					CustomPluginArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:custom-plugin/my-plugin/abc123"),
					Name:            ptr.String("my-plugin"),
				},
			},
		}, nil)

	lister := &KafkaConnectCustomPluginLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKafkaConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cp := resources[0].(*KafkaConnectCustomPlugin)
	a.Equal("my-plugin", *cp.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectCustomPlugin_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	mockClient.On("ListCustomPlugins", mock.Anything, mock.Anything).
		Return(&kafkaconnect.ListCustomPluginsOutput{
			CustomPlugins: []kafkaconnecttypes.CustomPluginSummary{},
		}, nil)

	lister := &KafkaConnectCustomPluginLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKafkaConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectCustomPlugin_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	cp := &KafkaConnectCustomPlugin{
		svc:             mockClient,
		CustomPluginArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:custom-plugin/my-plugin/abc123"),
	}

	mockClient.On("DeleteCustomPlugin", mock.Anything, &kafkaconnect.DeleteCustomPluginInput{
		CustomPluginArn: cp.CustomPluginArn,
	}).Return(&kafkaconnect.DeleteCustomPluginOutput{}, nil)

	a.NoError(cp.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectCustomPlugin_Properties(t *testing.T) {
	a := assert.New(t)

	cp := KafkaConnectCustomPlugin{
		CustomPluginArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:custom-plugin/my-plugin/abc123"),
		Name:            ptr.String("my-plugin"),
	}

	props := cp.Properties()
	a.Equal("my-plugin", props.Get("Name"))
}

func Test_Mock_KafkaConnectCustomPlugin_String(t *testing.T) {
	a := assert.New(t)
	cp := KafkaConnectCustomPlugin{Name: ptr.String("my-plugin")}
	a.Equal("my-plugin", cp.String())
}

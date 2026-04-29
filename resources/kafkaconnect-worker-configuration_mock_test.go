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

func Test_Mock_KafkaConnectWorkerConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	mockClient.On("ListWorkerConfigurations", mock.Anything, mock.Anything).
		Return(&kafkaconnect.ListWorkerConfigurationsOutput{
			WorkerConfigurations: []kafkaconnecttypes.WorkerConfigurationSummary{
				{
					WorkerConfigurationArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:worker-configuration/my-wc/abc123"),
					Name:                   ptr.String("my-wc"),
				},
			},
		}, nil)

	lister := &KafkaConnectWorkerConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKafkaConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	wc := resources[0].(*KafkaConnectWorkerConfiguration)
	a.Equal("my-wc", *wc.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectWorkerConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	mockClient.On("ListWorkerConfigurations", mock.Anything, mock.Anything).
		Return(&kafkaconnect.ListWorkerConfigurationsOutput{
			WorkerConfigurations: []kafkaconnecttypes.WorkerConfigurationSummary{},
		}, nil)

	lister := &KafkaConnectWorkerConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKafkaConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectWorkerConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKafkaConnectClient)

	wc := &KafkaConnectWorkerConfiguration{
		svc:                    mockClient,
		WorkerConfigurationArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:worker-configuration/my-wc/abc123"),
	}

	mockClient.On("DeleteWorkerConfiguration", mock.Anything, &kafkaconnect.DeleteWorkerConfigurationInput{
		WorkerConfigurationArn: wc.WorkerConfigurationArn,
	}).Return(&kafkaconnect.DeleteWorkerConfigurationOutput{}, nil)

	a.NoError(wc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_KafkaConnectWorkerConfiguration_Properties(t *testing.T) {
	a := assert.New(t)

	wc := KafkaConnectWorkerConfiguration{
		WorkerConfigurationArn: ptr.String("arn:aws:kafkaconnect:us-east-1:123456789012:worker-configuration/my-wc/abc123"),
		Name:                   ptr.String("my-wc"),
	}

	props := wc.Properties()
	a.Equal("my-wc", props.Get("Name"))
}

func Test_Mock_KafkaConnectWorkerConfiguration_String(t *testing.T) {
	a := assert.New(t)
	wc := KafkaConnectWorkerConfiguration{Name: ptr.String("my-wc")}
	a.Equal("my-wc", wc.String())
}

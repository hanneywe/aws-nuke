package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

func Test_Mock_IoTCustomMetric_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListCustomMetrics", mock.Anything, mock.Anything).
		Return(&iot.ListCustomMetricsOutput{
			MetricNames: []string{"my-custom-metric"},
		}, nil)

	lister := &IoTCustomMetricLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	metric := resources[0].(*IoTCustomMetric)
	a.Equal("my-custom-metric", *metric.MetricName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCustomMetric_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListCustomMetrics", mock.Anything, mock.Anything).
		Return(&iot.ListCustomMetricsOutput{
			MetricNames: []string{},
		}, nil)

	lister := &IoTCustomMetricLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCustomMetric_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	metric := &IoTCustomMetric{
		svc:        mockClient,
		MetricName: ptr.String("my-custom-metric"),
	}

	mockClient.On("DeleteCustomMetric", mock.Anything, &iot.DeleteCustomMetricInput{
		MetricName: metric.MetricName,
	}).Return(&iot.DeleteCustomMetricOutput{}, nil)

	a.NoError(metric.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCustomMetric_Properties(t *testing.T) {
	a := assert.New(t)

	metric := IoTCustomMetric{
		MetricName: ptr.String("my-custom-metric"),
	}

	props := metric.Properties()
	a.Equal("my-custom-metric", props.Get("MetricName"))
}

func Test_Mock_IoTCustomMetric_String(t *testing.T) {
	a := assert.New(t)
	metric := IoTCustomMetric{MetricName: ptr.String("my-custom-metric")}
	a.Equal("my-custom-metric", metric.String())
}
